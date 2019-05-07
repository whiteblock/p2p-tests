package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/kovetskiy/godocs"
	libp2p "github.com/libp2p/go-libp2p"
	relay "github.com/libp2p/go-libp2p-circuit"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	crypto "github.com/libp2p/go-libp2p-crypto"
	c "github.com/libp2p/go-libp2p-daemon/p2pclient"
	identify "github.com/libp2p/go-libp2p/p2p/protocol/identify"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	// "context"
	"os/signal"
)

func main() {

	identify.ClientVersion = "p2pd/0.1"
	id := flag.String("id", "", "peer identity; private key file")
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

	flag.Parse()

	var opts []libp2p.Option
	// var staticPeers []ma.Multiaddr

	if *id != "" {
		var r io.Reader
		r = rand.Reader
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

	// Logrus provides JSON logs.
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		("******Daemon Configuration******")
		"id":                             *id,
		"pubsubRouter":                   *pubsubRouter,
		"gossipsubHeartbeatInterval":     *gossipsubHeartbeatInterval,
		"gossipsubHeartbeatInitialDelay": *gossipsubHeartbeatInitialDelay,
		"relayEnabled":                   *relayEnabled,
		"relayActive":                    *relayActive,
		"relayHop":                       *relayHop,
		"hostAddrs":                      *hostAddrs,
		"announceAddrs":                  *announceAddrs,
	}).Fatal("something bad happened")

	// gets the options to pass to the daemon

	d1, c1, closer, err := createDaemonClientPair(opts, *pubsubRouter, *pubsubSign, *pubsubSignStrict, *gossipsubHeartbeatInterval, *gossipsubHeartbeatInitialDelay)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Printf("%#v", *d1))
	fmt.Println(fmt.Printf("%#v", *c1))

	defer closer()

	testProtos := []string{"/test"}

	err = c1.NewStreamHandler(testProtos, func(info *c.StreamInfo, conn io.ReadWriteCloser) {
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

	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		endWaiter.Done()
	}()

	fmt.Printf("Daemon started")
	endWaiter.Wait()
}
