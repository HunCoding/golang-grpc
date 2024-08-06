package client

import (
	"context"
	"log"
	"time"

	"github.com/HunCoding/golang-grpc/pb"

	"google.golang.org/grpc"
)

func RunClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AddUser(ctx, &pb.AddUserRequest{Name: "Alice", Age: 30})
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}
	log.Printf("User ID: %s added successfully", r.GetId())

	user, err := c.GetUser(ctx, &pb.GetUserRequest{Id: r.GetId()})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("User: %s, Age: %d", user.GetName(), user.GetAge())
}
