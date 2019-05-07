package main

import(
	"log"
	"net"
	"net/rpc"
	"net/http"
)

type Message struct {
	payload string
}

func CreateRpcClient() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	return client
}

func StartRpcServer() error {
	// Listen to TPC connections on port
	listener, e := net.Listen("tcp", ":9000")
	if e != nil {
		log.Fatal("Listen error: ", e)
		return e
	}
	log.Printf("Serving RPC server on port %d", 9000)
	// Start accept incoming HTTP connections
	err := http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
		return err
	}
	return nil
}

