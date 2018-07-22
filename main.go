package main

import (
	// "runtime"
	"time"
)

const vtid = 150585302

func main() {
	// runtime.GOMAXPROCS(8)

	pingsCh := make(chan Ping, 10)

	go Bot(pingsCh)
	go Server(pingsCh)

	defer DBClose()

	for {
		time.Sleep(1 * time.Second)
	}
}
