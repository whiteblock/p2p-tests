package main

import (
	"os"
	"time"
	"context"
	"runtime"
	"io/ioutil"
	"path/filepath"
	libp2p "github.com/libp2p/go-libp2p"
	ps "github.com/libp2p/go-libp2p-pubsub"
	p2pd "github.com/libp2p/go-libp2p-daemon"
	ma "github.com/multiformats/go-multiaddr"
	c "github.com/libp2p/go-libp2p-daemon/p2pclient"
)

type makeEndpoints func() (daemon, client ma.Multiaddr, cleanup func(), err error)

func createTempDir() (string, string, func(), error) {
	root := os.TempDir()
	dir, err := ioutil.TempDir(root, "p2pd")
	if err != nil {
		return "", "", nil, err
	}
	daemonPath := filepath.Join(dir, "daemon.sock")
	clientPath := filepath.Join(dir, "client.sock")
	closer := func() {
		os.RemoveAll(dir)
	}
	return daemonPath, clientPath, closer, nil
}

func makeTcpLocalhostEndpoints() (daemon, client ma.Multiaddr, cleanup func(), err error) {
	daemon, err = ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	if err != nil {
		return nil, nil, nil, err
	}
	client, err = ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	if err != nil {
		return nil, nil, nil, err
	}
	cleanup = func() {}
	return daemon, client, cleanup, nil
}

func makeUnixEndpoints() (daemon, client ma.Multiaddr, cleanup func(), err error) {
	daemonPath, clientPath, cleanup, err := createTempDir()
	if err != nil {
		return nil, nil, nil, err
	}
	daemon, _ = ma.NewComponent("unix", daemonPath)
	client, _ = ma.NewComponent("unix", clientPath)
	return
}

func getEndpointsMaker() makeEndpoints {
	if runtime.GOOS == "windows" {
		return makeTcpLocalhostEndpoints
	} else {
		return makeUnixEndpoints
	}
}

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