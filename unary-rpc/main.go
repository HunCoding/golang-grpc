package main

import (
	"github.com/HunCoding/golang-grpc/client"
	"github.com/HunCoding/golang-grpc/server"
)

func main() {
	go server.RunServer()
	client.RunClient()
}
