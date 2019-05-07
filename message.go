package main

import(
	"time"
)

type Message struct {
	from string
	to string
	timestamp time.Time
	data string
}