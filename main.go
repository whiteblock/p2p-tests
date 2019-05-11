package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"github.com/spf13/pflag"
	"io/ioutil"
	//"net/rpc"
	"strings"
	"crypto/rand"
	mrand "math/rand"
	"encoding/hex"
	"encoding/json"
	"time"
	"github.com/libp2p/go-libp2p"
	logrus "github.com/sirupsen/logrus"
	//tmpps "github.com/libp2p/go-libp2p-peerstore/pstoremem"
	//man "github.com/multiformats/go-multiaddr-net"
	ma "github.com/multiformats/go-multiaddr"
	relay "github.com/libp2p/go-libp2p-circuit"
	crypto "github.com/libp2p/go-libp2p-crypto"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	c "github.com/libp2p/go-libp2p-daemon/p2pclient"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
)

var (
	portStartPoint int
	bindIP string
	maddrs []string
	sendInterval int64
)

func main() {
	var generateOnly bool

	log.SetFlags(log.LstdFlags | log.Llongfile)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	flag := pflag.NewFlagSet("p2p", pflag.ExitOnError)

	//var rawPeers []string
	var fileName string

	identify.ClientVersion = "p2pd/0.1"
	id := flag.String("id", "", "peer identity; private key file")
	seed := flag.Int64("seed", 0, "seed to generate peer identity deterministically")
	flag.BoolVar(&generateOnly,"generate-only",false,"generate ")
	connMgr := flag.Bool("connManager", false, "Enables the Connection Manager")
	connMgrLo := flag.Int("connLo", 256, "Connection Manager Low Water mark")
	connMgrHi := flag.Int("connHi", 512, "Connection Manager High Water mark")
	connMgrGrace := flag.Duration("connGrace", 120, "Connection Manager grace period (in seconds)")
	natPortMap := flag.Bool("natPortMap", false, "Enables NAT port mapping")
	pubsubRouter := flag.String("pubsubRouter", "gossipsub", "Specifies the pubsub router implementation")
	pubsubSign := flag.Bool("pubsubSign", false, "Enables pubsub message signing")
	pubsubSignStrict := flag.Bool("pubsubSignStrict", false, "Enables pubsub strict signature verification")
	gossipsubHeartbeatInterval := flag.Duration("gossipsubHeartbeatInterval", 0, "Specifies the gossipsub heartbeat interval")
	gossipsubHeartbeatInitialDelay := flag.Duration("gossipsubHeartbeatInitialDelay", 0, "Specifies the gossipsub initial heartbeat delay")
	relayEnabled := flag.Bool("relay", true, "Enables circuit relay")
	relayActive := flag.Bool("relayActive", false, "Enables active mode for relay")
	relayHop := flag.Bool("relayHop", false, "Enables hop for relay")
	hostAddrs := flag.String("hostAddrs", "", "comma separated list of multiaddrs the host should listen on")
	announceAddrs := flag.String("announceAddrs", "", "comma separated list of multiaddrs the host should announce to the network")
	noListen := flag.Bool("noListenAddrs", false, "sets the host to listen on no addresses")
	//peerStore := flag.String("peerstore", "", "peers to add to the daemon's peerstore")

	//flag.StringSliceVarP(&rawPeers,"peer","p",[]string{},"peers")
	flag.IntVar(&portStartPoint,"port-start",8999,"port start")
	flag.Int64Var(&sendInterval,"send-interval",-1,"interval to send messages, -1 means don't send")
	flag.StringVar(&bindIP,"ip","127.0.0.1","ip address to bind on")
	//flag.StringSliceVar(&maddrs, "maddr", []string{}, "addresses that daemon owns")
	flag.StringVar(&fileName,"file","static-peers.json","file of the peers")

	flag.Parse(os.Args)

	peers,err := CreatePeerInfosFromFile(fileName)
	if err != nil {
		panic(err)
	}

	var opts []libp2p.Option


	if *id == "" {
		var r io.Reader
		if *seed == int64(0) {
			r = rand.Reader
		} else {
			r = mrand.New(mrand.NewSource(*seed))
		}
		priv, _, err := crypto.GenerateEd25519Key(r)
		if err != nil {
			panic(err)
		}
		opts = append(opts, libp2p.Identity(priv))
	}

	if *hostAddrs != "" {
		addrs := strings.Split(*hostAddrs, ",")
		opts = append(opts, libp2p.ListenAddrStrings(addrs...))
	}

	if *announceAddrs != "" {
		addrs := strings.Split(*announceAddrs, ",")
		maddrs := make([]ma.Multiaddr, 0, len(addrs))
		for _, a := range addrs {
			maddr, err := ma.NewMultiaddr(a)
			if err != nil {
				log.Fatal(err)
			}
			maddrs = append(maddrs, maddr)
		}
		opts = append(opts, libp2p.AddrsFactory(func([]ma.Multiaddr) []ma.Multiaddr {
			return maddrs
		}))
	}

	if *connMgr {
		cm := connmgr.NewConnManager(*connMgrLo, *connMgrHi, *connMgrGrace)
		opts = append(opts, libp2p.ConnectionManager(cm))
	}

	if *natPortMap {
		opts = append(opts, libp2p.NATPortMap())
	}

	if *relayEnabled {
		var relayOpts []relay.RelayOpt
		if *relayActive {
			relayOpts = append(relayOpts, relay.OptActive)
		}
		if *relayHop {
			relayOpts = append(relayOpts, relay.OptHop)
		}
		opts = append(opts, libp2p.EnableRelay(relayOpts...))
	}

	if *noListen {
		opts = append(opts, libp2p.NoListenAddrs)
	}

	// gets the options to pass to the daemon
	d, cl, closer, ctx, err := createDaemonClientPair(opts)
	if err != nil {
		panic(err)
	}
	clientPID,clientAddrs,_ := cl.Identify()
	logrus.WithFields(logrus.Fields{
		"pid":clientPID,
		"addrs":clientAddrs,
	}).Info("Created a client")

	logrus.WithFields(logrus.Fields{
		"pid":d.ID(),
		"addrs":d.Addrs(),
	}).Info("Created a Daemon")
	if generateOnly {
		os.Exit(0)
	}

	err = pubsub(d, *pubsubRouter, *pubsubSign, *pubsubSignStrict, *gossipsubHeartbeatInterval, *gossipsubHeartbeatInitialDelay)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %s\n",d.ID().Pretty())


	fmt.Printf("DAEMON PEERLIST: %v\n", d.Addrs())

	for _,peer := range peers {
		logrus.WithFields(logrus.Fields{
					"peer":peer.ID.Pretty(),
					"addrs":peer.Addrs,
				}).Info("Dialing peer")
		for i := 0; i < 20; i++ {
			err = cl.Connect(peer.ID,peer.Addrs)
			if err == nil {

			}else{
				logrus.WithFields(logrus.Fields{"timeout":"2s"}).Warn("Failed to connect")
			}
			time.Sleep(2*time.Second)
		}
		if err != nil {
			panic(err)
		}
	}

	dch, err := cl.Subscribe(ctx, "jargon")
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for{
			msg := <- dch
			var data interface{}
			err := json.Unmarshal([]byte(msg.Data),&data)
			if err != nil {
				data = string(msg.Data)
			}
			logrus.WithFields(logrus.Fields{
				"from":hex.EncodeToString(msg.From),
				"data":data,
				"seqno":hex.EncodeToString(msg.Seqno),
				"topicIDs":msg.TopicIDs,
				"signature":hex.EncodeToString(msg.Signature),
				"key":hex.EncodeToString(msg.Key),
			}).Info("Received a message")
		}
	}()
 	var counter int64
 	go func(){
 		for sendInterval > 0 || counter == 0 { //infinitely loop if > 0
			out,err := json.Marshal(map[string]interface{}{
				"id":counter,
				"jsonrpc":"2.0",
				"params":[]string{"foo","bar","foofoo","foobar","barfoo","barbar"},
				"method":"far",
			})
			logrus.WithFields(logrus.Fields{
					"sending":string(out),
					"error":err,
				}).Info("Sending a message")
			cl.Publish("jargon", out)

			counter++
			time.Sleep(time.Duration(sendInterval)*time.Microsecond)
			if counter > 1000000000 {
				counter = 0
			}
		}
 	}()
	
	defer closer()

	testProtos := []string{"/test"}

	err = cl.NewStreamHandler(testProtos, func(info *c.StreamInfo, conn io.ReadWriteCloser) {
		defer conn.Close()
		logrus.WithFields(logrus.Fields{
			"peer":info.Peer.Pretty(),
			"addr":info.Addr,
			"proto":info.Proto,
		}).Info("Received a stream")

		in,err :=  ioutil.ReadAll(conn)
		if err != nil {
			logrus.WithFields(logrus.Fields{"err":err}).Info("Error getting data")
		}
		logrus.WithFields(logrus.Fields{"data":in}).Info("Got some data")		
	})

	if err != nil {
		panic(err)
	}
	chanwait()

}
