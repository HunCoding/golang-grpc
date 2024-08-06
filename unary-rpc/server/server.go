package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/HunCoding/golang-grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type User struct {
	ID   string
	Name string
	Age  int32
}

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	users map[string]User
	mu    sync.Mutex
}

func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		users: make(map[string]User),
	}
}

func (s *UserServiceServer) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	user := User{
		ID:   id,
		Name: req.Name,
		Age:  req.Age,
	}

	s.users[id] = user

	return &pb.AddUserResponse{Id: id}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.Id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &pb.GetUserResponse{
		Name: user.Name,
		Age:  user.Age,
	}, nil
}

func generateID() string {
	return fmt.Sprintf("%d", rand.Int())
}

func RunServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, NewUserServiceServer())
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
