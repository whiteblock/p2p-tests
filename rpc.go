package main

import(
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/http"
)

// rpc server
type Server struct {
	Listener net.Listener
}

func (s *Server) Close() (err error) {
	if s.Listener != nil {
		err = s.Listener.Close()
		return err
	}
	return nil
}

func (s *Server) Start() (err error) {
	handler := new(Handler)
	rpc.Register(handler)
	s.Listener, err = net.Listen("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatal("Listen error: ", err)
		return err
	}
	fmt.Printf("%#v",s)
	rpc.HandleHTTP()
	http.Serve(s.Listener, nil)
	return nil
}

// client
type Client struct {
	Client  *rpc.Client
}

func (c *Client) Init() (err error) {
	c.Client, err = rpc.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatal("Connection error: ", err)
		return err
	}
	fmt.Printf("%#v",c)
	return nil
}

func (c *Client) Close() (err error) {
	if c.Client != nil {
		err = c.Client.Close()
		return err
	}
	return nil
}

func (c *Client) ClientExec(name string) (msg Message, err error) {
	var (
		request  = Request{Name: name}
		response = new(Response)
	)
	err = c.Client.Call("Handler.Execute", request, response)
	if err != nil {
		return msg, err
	}
	msg = response.Message
	return msg, nil
}