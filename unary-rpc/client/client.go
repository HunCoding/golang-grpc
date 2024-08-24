package client

import (
	"context"
	"fmt"
	"github.com/HunCoding/golang-grpc/unary-rpc/pb"
	"google.golang.org/grpc"
)

func Run() {
	dial, err := grpc.Dial("fdjksal:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer dial.Close()

	userClient := pb.NewUserClient(dial)

	user, err := userClient.AddUser(context.Background(), &pb.AddUserRequest{
		Id:   "1",
		Age:  40,
		Name: "Test",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("First user created: %v\n", user)

	user, err = userClient.AddUser(context.Background(), &pb.AddUserRequest{
		Id:   "2",
		Age:  42,
		Name: "Test 2",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Second user created: %v\n", user)

	getUserResponse, err := userClient.GetUser(context.Background(), &pb.GetUserRequest{Id: "2"})
	if err != nil {
		return
	}
	fmt.Printf("User returned from GetUser method: %v\n", getUserResponse)
}
