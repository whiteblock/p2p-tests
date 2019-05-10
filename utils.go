package main

import (
	"os"
	"fmt"
	"sync"
	"os/signal"
	"strings"
	peer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
	ps "github.com/libp2p/go-libp2p-peerstore"
)

func chanwait() {
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()
}

// func handleSignals() {
// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
// 	<-signals
// 	log.Println("signal received")
// }

func CreatePeerInfos(peers []string) ([]ps.PeerInfo,error) {
	out := []ps.PeerInfo{}
	for _,rawPeer := range peers {
		pidSock := strings.Split(rawPeer,"@")
		socket := pidSock[1]
		rawPid := pidSock[0]
		ipPort := strings.Split(socket,":")
		port := ipPort[1]
		ip := ipPort[0]

		mAddr,err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s",ip,port))
		if err != nil {
			return nil,err
		}
		test,err := ps.InfoFromP2pAddr(mAddr)
		fmt.Printf("INFO: %#vn\n",test)
		pid,err := peer.IDB58Decode(rawPid)
		if err != nil {
			return nil,err
		}

		out = append(out,ps.PeerInfo{
			ID:pid,
			Addrs: []ma.Multiaddr{mAddr},
		})
	}
	return out,nil
}