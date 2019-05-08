package main

import (
	"os"
	"log"
	"sync"
	"os/signal"
	// "syscall"
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

func runServer() {
	server := new(Server)
	defer server.Close()
	// handleSignals()
	err := server.Start()
	if err != nil {
		panic(err)
	}
}

func runClient() {
	client := new(Client)
	defer client.Close()
	err := client.Init()
	if err != nil {
		panic(err)
	}
	for {
		response, err := client.ClientExec("peepeepoopookak")
		if err != nil {
			panic(err)
		}
		log.Println(response)
	}
}

// func handleSignals() {
// 	signals := make(chan os.Signal, 1)

// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
// 	<-signals
// 	log.Println("signal received")
// }