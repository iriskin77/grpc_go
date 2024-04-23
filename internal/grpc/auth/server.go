package auth

import (
	"context"

	grpc_go "github.com/iriskin77/grpc_go/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	grpc_go.UnimplementedAuthServer
}

// Это функция, которая регистрирует auth сервис
func Register(gRPC *grpc.Server) {
	grpc_go.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *grpc_go.LoginRequest) (*grpc_go.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *grpc_go.RegisterRequest) (*grpc_go.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *grpc_go.IsAdminRequest) (*grpc_go.IsAdminResponse, error) {
	panic("implement me")
}
