package main

import (
	"fmt"
	"context"
	"github.com/libp2p/go-libp2p"
	p2pd "github.com/libp2p/go-libp2p-daemon"
	c "github.com/libp2p/go-libp2p-daemon/p2pclient"
	ps "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	"time"
)

func createDaemonClientPair(opts []libp2p.Option) (*p2pd.Daemon, *c.Client, func(), context.Context, error) {
	ctx, _:= context.WithCancel(context.Background())

	dAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d",bindIP,portStartPoint))
	cmaddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d",bindIP,portStartPoint+1))

	daemon, err := p2pd.NewDaemon(ctx, dAddr, "", opts ...)
	if err != nil {
		return nil, nil, nil, nil, err
	}//daemon.Listener()
	client, err := c.NewClient(daemon.Listener().Multiaddr(), cmaddr)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	closer := func() {
		_ = client.Close()
		_ = daemon.Close()
	}
	return daemon, client, closer, ctx, nil
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