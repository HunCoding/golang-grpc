package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/HunCoding/golang-grpc/unary-rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func Run() {
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})

	dial, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}

	defer dial.Close()

	userClient := pb.NewUserClient(dial)

	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		loginResp, err := userClient.Login(ctx, &pb.LoginRequest{Username: "user", Password: "password"})
		if err != nil {
			log.Fatalf("Falha ao fazer login: %v", err)
		}
		token := loginResp.Token
		log.Printf("Token JWT recebido: %s", token)

		md := metadata.New(map[string]string{"authorization": "Bearer " + token})
		ctx = metadata.NewOutgoingContext(ctx, md)

		user, err := userClient.AddUser(ctx, &pb.AddUserRequest{
			Id:   "1",
			Age:  40,
			Name: "Test",
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("First user created: %v\n", user)

		user, err = userClient.AddUser(ctx, &pb.AddUserRequest{
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

		time.Sleep(500 * time.Millisecond) // Intervalo de 500ms entre as tentativas
	}
}
