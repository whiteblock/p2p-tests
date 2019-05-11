package main

import (
	"os"
	"fmt"
	"log"
	"sync"
	"os/signal"
	"io/ioutil"
	"strings"
	"encoding/json"
	logrus "github.com/sirupsen/logrus"
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

type SerialPeerInfo struct {
	ID 		string
	MAddrs []string
}

func (spi SerialPeerInfo) Convert() (ps.PeerInfo,error) {
	out := ps.PeerInfo{}
	pid,err := peer.IDB58Decode(spi.ID)
	if err != nil {
		return ps.PeerInfo{},err
	}
	out.ID = pid
	for _,maddrStr :=  range spi.MAddrs{
		mAddr,err := ma.NewMultiaddr(maddrStr)
		if err != nil {
			return ps.PeerInfo{},err
		}
		out.Addrs = append(out.Addrs,mAddr)
	}
	logrus.WithFields(logrus.Fields{
			"peer":out,
	}).Info("Parsed peer...")
	return out,nil
}

func CreatePeerInfosFromFile(filename string) ([]ps.PeerInfo,error) {
	res,err := ioutil.ReadFile(filename)
	if err != nil {
		return nil,nil
	}
	var rawPeers []SerialPeerInfo
	err = json.Unmarshal(res,&rawPeers)
	if err != nil {
		log.Println(err)
		return nil,nil
	}
	out := []ps.PeerInfo{}
	for _,rawPeer := range rawPeers{
		logrus.WithFields(logrus.Fields{
				"ID":rawPeer.ID,
				"MAddrs":rawPeer.MAddrs,
		}).Info("Parsing peer...")

		peer, err := rawPeer.Convert()
		if err != nil {
			return nil,err
		}
		out = append(out,peer)
	}
	return out,nil
}