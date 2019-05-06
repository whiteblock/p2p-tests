package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"sync"
	"flag"
	"strings"
	// "context"
	"os/signal"
	"crypto/rand"
	mrand "math/rand"
	libp2p "github.com/libp2p/go-libp2p"
	ma "github.com/multiformats/go-multiaddr"
	crypto "github.com/libp2p/go-libp2p-crypto"
	relay "github.com/libp2p/go-libp2p-circuit"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	c "github.com/libp2p/go-libp2p-daemon/p2pclient"
	identify "github.com/libp2p/go-libp2p/p2p/protocol/identify"
)

func main() {

	identify.ClientVersion = "p2pd/0.1"

	seed := flag.Int64("seed", 0, "seed to generate privKey")
	id := flag.String("id", "", "peer identity; private key file")
	connMgr := flag.Bool("connManager", false, "Enables the Connection Manager")
	connMgrLo := flag.Int("connLo", 256, "Connection Manager Low Water mark")
	connMgrHi := flag.Int("connHi", 512, "Connection Manager High Water mark")
	connMgrGrace := flag.Duration("connGrace", 120, "Connection Manager grace period (in seconds)")
	natPortMap := flag.Bool("natPortMap", false, "Enables NAT port mapping")
	// pubsub := flag.Bool("pubsub", false, "Enables pubsub")
	// pubsubRouter := flag.String("pubsubRouter", "gossipsub", "Specifies the pubsub router implementation")
	// pubsubSign := flag.Bool("pubsubSign", true, "Enables pubsub message signing")
	// pubsubSignStrict := flag.Bool("pubsubSignStrict", false, "Enables pubsub strict signature verification")
	// gossipsubHeartbeatInterval := flag.Duration("gossipsubHeartbeatInterval", 0, "Specifies the gossipsub heartbeat interval")
	// gossipsubHeartbeatInitialDelay := flag.Duration("gossipsubHeartbeatInitialDelay", 0, "Specifies the gossipsub initial heartbeat delay")
	relayEnabled := flag.Bool("relay", true, "Enables circuit relay")
	relayActive := flag.Bool("relayActive", false, "Enables active mode for relay")
	relayHop := flag.Bool("relayHop", false, "Enables hop for relay")
	hostAddrs := flag.String("hostAddrs", "", "comma separated list of multiaddrs the host should listen on")
	announceAddrs := flag.String("announceAddrs", "", "comma separated list of multiaddrs the host should announce to the network")
	noListen := flag.Bool("noListenAddrs", false, "sets the host to listen on no addresses")

	flag.Parse()

	var opts []libp2p.Option
	// var staticPeers []ma.Multiaddr

	if *id != "" {
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
		// if *relayDiscovery {
		// 	relayOpts = append(relayOpts, relay.OptDiscovery)
		// }
		opts = append(opts, libp2p.EnableRelay(relayOpts...))
	}

	if *noListen {
		opts = append(opts, libp2p.NoListenAddrs)
	}

	fmt.Println("******Daemon Configuration******")
	fmt.Println("seed: ",*seed)
	fmt.Println("id: ", *id)
	fmt.Println("connection manager: ", *connMgr)
	fmt.Println("connmgrLo: ",*connMgrLo)
	fmt.Println("connmgrHi: ",*connMgrHi)
	fmt.Println("connMgrGrace: ", *connMgrGrace)
	fmt.Println("natPortMap: ", *natPortMap)
	// fmt.Println(*pubsub)
	// fmt.Println(*pubsubRouter)
	// fmt.Println(*pubsubSign)
	// fmt.Println(*pubsubSignStrict)
	// fmt.Println(*gossipsubHeartbeatInterval)
	// fmt.Println(*gossipsubHeartbeatInitialDelay)
	fmt.Println("relayEnabled: ", *relayEnabled)
	fmt.Println("relayActive: ", *relayActive)
	fmt.Println("relayHop: ", *relayHop)
	fmt.Println("hostAddrs: ", *hostAddrs)
	fmt.Println("announceAddrs: ", *announceAddrs)
	fmt.Println("noListen: ", *noListen)
	fmt.Println("********************************")

	// gets the options to pass to the deamon

	d1, c1, closer, err := createDaemonClientPair(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Printf("%#v",*d1))
	fmt.Println(fmt.Printf("%#v",*c1))

	defer closer()

	testprotos := []string{"/test"}

	err = c1.NewStreamHandler(testprotos, func(info *c.StreamInfo, conn io.ReadWriteCloser) {
		defer conn.Close()
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		// TODO print to JSON file.
		fmt.Printf("%s\n", buf[0:n])
	})

	if err != nil {
		panic(err)
	}

	// TODO open a RPC port to publish messages

	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()

	fmt.Printf("Daemon started")
	end_waiter.Wait()
}
