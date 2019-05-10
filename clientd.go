package main

import (
	"context"
	"github.com/libp2p/go-libp2p"
	p2pd "github.com/libp2p/go-libp2p-daemon"
	c "github.com/libp2p/go-libp2p-daemon/p2pclient"
	ps "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	"time"
)

func createDaemonClientPair(opts []libp2p.Option) (*p2pd.Daemon, *c.Client, func(), error) {
	ctx, _:= context.WithCancel(context.Background())
	dAddr, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/8999")
	cmaddr, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/9000")

	daemon, err := p2pd.NewDaemon(ctx, dAddr, "", opts ...)
	if err != nil {
		return nil, nil, nil, err
	}//daemon.Listener()
	client, err := c.NewClient(daemon.Listener().Multiaddr(), cmaddr)
	if err != nil {
		return nil, nil, nil, err
	}

	closer := func() {
		_ = client.Close()
		_ = daemon.Close()
	}
	return daemon, client, closer, nil
}

func pubsub(daemon *p2pd.Daemon, pubsubRouter string, pubsubSign, pubsubSignStrict bool, gossipHearbeatInterval, gossipHeartBeatInitialDelay time.Duration)error{
	if gossipHearbeatInterval < 0{
		ps.GossipSubHeartbeatInterval = gossipHearbeatInterval
		ps.GossipSubHeartbeatInitialDelay = gossipHeartBeatInitialDelay
	}
	err := daemon.EnablePubsub(pubsubRouter, pubsubSign, pubsubSignStrict)
	if err != nil{
		return err
	}
	return nil
}