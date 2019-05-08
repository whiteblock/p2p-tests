package main

import(
	"fmt"
	"time"
	"encoding/json"
)

type Message struct {
	Origin      string            `json:"origin"` //node from which the message was initially sent
	Destination string            `json:"destination"` //node to which the message is being sent
	Relayer     string			  `json:"relayer"` //last node which relayed the message
	Timestamp   time.Time         `json:"timestamp"` //time at which the message was first broadcast
	Data        string            `json:"data"` //message payload
	MessageID   string			  `json:"messageid"` //a unique identified associated with a message
	Nonce       int64             `json:"nonce"` //counter representing order in which message was sent from origin
}

func (m *Message) Relay() ([]byte, error){
	out, err := json.Marshal(m)
	if err!= nil {
		return nil, err
	}
	fmt.Println(out)
	return out, nil
}