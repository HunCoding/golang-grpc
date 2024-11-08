package server

import (
	"context"
	"errors"
	"github.com/HunCoding/golang-grpc/unary-rpc/pb"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"
	"time"
)

type User struct {
	ID   string
	Name string
	Age  int32
}

func Run() {
	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Erro ao carregar certificados TLS: %v", err)
	}

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(ChainUnaryInterceptors(
			authInterceptor,
			rateLimitInterceptor,
			logInterceptor,
		)))
	pb.RegisterUserServer(s, NewUserService())
	reflection.Register(s)

	err = s.Serve(listen)
	if err != nil {
		panic(err)
	}
}

type UserService struct {
	pb.UnimplementedUserServer

	users map[string]*User
	mu    sync.Mutex
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
	}
}

func (us *UserService) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	user := &User{
		ID:   req.Id,
		Name: req.Name,
		Age:  req.Age,
	}

	us.users[user.ID] = user

	return &pb.AddUserResponse{
		Id:   user.ID,
		Age:  user.Age,
		Name: user.Name,
	}, nil
}

func (us *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, ok := us.users[req.Id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return &pb.GetUserResponse{
		Id:   user.ID,
		Age:  user.Age,
		Name: user.Name,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Username == "user" && req.Password == "password" {
		token, err := generateJWT(req.Username)
		if err != nil {
			return nil, err
		}
		return &pb.LoginResponse{Token: token}, nil
	}
	return nil, errors.New("usuário ou senha inválidos")
}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(jwtKey)
}

var limiter = rate.NewLimiter(0.5, 1) // Permite 1 requisição a cada 2 segundos
var mu sync.Mutex

func rateLimitInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	mu.Lock()
	defer mu.Unlock()
	if !limiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "limite de requisições excedido")
	}
	return handler(ctx, req)
}

func logInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("Recebendo requisição para método %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Erro: %v", err)
	}
	return resp, err
}

func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Função handler que chamará o próximo interceptador da cadeia
		currentHandler := handler

		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := currentHandler
			currentHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return interceptor(currentCtx, currentReq, info, next)
			}
		}

		return currentHandler(ctx, req)
	}
}
