package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-daemon"
	"github.com/libp2p/go-libp2p-daemon/p2pclient"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
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

func createDaemon(daemonAddr ma.Multiaddr) (*p2pd.Daemon, func(), error) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	daemon, err := p2pd.NewDaemon(ctx, daemonAddr, "")
	if err != nil {
		return nil, nil, err
	}
	err = daemon.EnablePubsub("gossipsub", false, false)
	if err != nil {
		return nil, nil, err
	}
	return daemon, cancelCtx, nil
}

func createClient(daemonAddr ma.Multiaddr, clientAddr ma.Multiaddr) (*p2pclient.Client, func(), error) {
	client, err := p2pclient.NewClient(daemonAddr, clientAddr)
	if err != nil {
		return nil, nil, err
	}
	closer := func() {
		client.Close()
	}
	return client, closer, nil
}

func createDaemonClientPair() (*p2pd.Daemon, *p2pclient.Client, func(), error) {
	dmaddr, cmaddr, dirCloser, err := getEndpointsMaker()()
	daemon, closeDaemon, err := createDaemon(dmaddr)
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

func main() {
	_, c1, closer, err := createDaemonClientPair()
	if err != nil {
		panic(err)
	}
	defer closer()

	testprotos := []string{"/test"}

	err = c1.NewStreamHandler(testprotos, func(info *p2pclient.StreamInfo, conn io.ReadWriteCloser) {
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
