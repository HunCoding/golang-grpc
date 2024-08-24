package main

import (
	"github.com/HunCoding/golang-grpc/unary-rpc/client"
	"github.com/HunCoding/golang-grpc/unary-rpc/server"
)

func main() {
	go server.Run()
	client.Run()
}
