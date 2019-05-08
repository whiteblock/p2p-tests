package main

import(
	"fmt"
	"time"
	// "encoding/json"
)

type Message struct {
	Origin      string            `json:"origin"`
	Destination string            `json:"destination"`
	Timestamp   time.Time         `json:"timestamp"`
	Data        string            `json:"data"`
}

type Response struct {
	Message Message
}

type Request struct {
	Name string
}

type Handler struct {}

var (
	request  = Request{Name: "ok"}
	response = new(Response)
)

func (h *Handler) Execute(req Request, res *Response) (err error) {
	if req.Name == "" {
		fmt.Println("A name must be specified")
		return nil
	}
	res.Message = Message{Origin:"0x0"}
	return nil
}