package server

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var jwtKey = []byte("minhaChaveSecreta")

func validateJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), nil
	}
	return "", errors.New("token inválido")
}

func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if info.FullMethod == "/pb.User/Login" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadados ausentes")
	}

	token := md["authorization"]
	if len(token) == 0 {
		return nil, errors.New("token ausente")
	}

	username, err := validateJWT(token[0][7:])
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, "username", username)
	return handler(newCtx, req)
}
