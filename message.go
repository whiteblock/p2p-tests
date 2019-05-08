package main

import(
	"fmt"
	"time"
	"encoding/json"
)

type Message struct {
	Origin      string            `json:"origin"`
	Destination string            `json:"destination"`
	Timestamp   time.Time         `json:"timestamp"`
	Data        string            `json:"data"`
}

func (m *Message) Relay() ([]byte, error){
	out, err := json.Marshal(m)
	if err!= nil {
		return nil, err
	}
	fmt.Println(out)
	return out, nil
}