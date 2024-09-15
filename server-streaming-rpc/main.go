package main

import (
	"github.com/HunCoding/server-streaming-rpc/client"
	"github.com/HunCoding/server-streaming-rpc/server"
)

func main() {
	go server.Run()
	client.Run()
}
