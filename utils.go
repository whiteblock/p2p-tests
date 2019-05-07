package main

import (
	"os"
	"sync"
	"runtime"
	"os/signal"
	"io/ioutil"
	"path/filepath"
	ma "github.com/multiformats/go-multiaddr"
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
