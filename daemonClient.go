package main

import(
	"time"
	"context"
	libp2p "github.com/libp2p/go-libp2p"
	ps "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	p2pd "github.com/libp2p/go-libp2p-daemon"
	c "github.com/libp2p/go-libp2p-daemon/p2pclient"
)

func createDaemon(daemonAddr ma.Multiaddr, opts []libp2p.Option, pubsubRouter string, pubsubSign, pubsubSignStrict bool, gossipsubHeartbeatInterval, gossipsubHeartbeatInitialDelay time.Duration) (*p2pd.Daemon, func(), error) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	daemon, err := p2pd.NewDaemon(ctx, daemonAddr, "", opts...)
	if err != nil {
		return nil, nil, err
	}
	
	if gossipsubHeartbeatInterval > 0 {
		ps.GossipSubHeartbeatInterval = gossipsubHeartbeatInterval
		ps.GossipSubHeartbeatInitialDelay = gossipsubHeartbeatInitialDelay
	}
	

	err = daemon.EnablePubsub(pubsubRouter, pubsubSign, pubsubSignStrict)
	if err != nil {
		return nil, nil, err
	}
	return daemon, cancelCtx, nil
}

func createClient(daemonAddr ma.Multiaddr, clientAddr ma.Multiaddr) (*c.Client, func(), error) {
	client, err := c.NewClient(daemonAddr, clientAddr)
	if err != nil {
		return nil, nil, err
	}
	closer := func() {
		client.Close()
	}
	return client, closer, nil
}

func createDaemonClientPair(opts []libp2p.Option, pubsubRouter string, pubsubSign, pubsubSignStrict bool, gossipsubHeartbeatInterval, gossipsubHeartbeatInitialDelay time.Duration) (*p2pd.Daemon, *c.Client, func(), error) {
	dmaddr, cmaddr, dirCloser, err := getEndpointsMaker()()
	
	daemon, closeDaemon, err := createDaemon(dmaddr, opts, pubsubRouter, pubsubSign, pubsubSignStrict, gossipsubHeartbeatInterval, gossipsubHeartbeatInitialDelay)
	if err != nil {
		return nil, nil, nil, err
	}
	client, closeClient, err := createClient(daemon.Listener().Multiaddr(), cmaddr)
	if err != nil {
		return nil, nil, nil, err
	}

	closer := func() {
		closeDaemon()
		closeClient()
		dirCloser()
	}
	return daemon, client, closer, nil
}