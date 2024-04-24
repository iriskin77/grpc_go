package auth

import (
	"context"
	"fmt"

	grpc_go "github.com/iriskin77/grpc_go/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	grpc_go.UnimplementedAuthServer
	auth Auth
}

// Это функция, которая регистрирует auth сервис
func Register(gRPC *grpc.Server, auth Auth) {
	grpc_go.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *grpc_go.LoginRequest) (*grpc_go.LoginResponse, error) {

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app is required")
	}

	// TODO: implements login via auth service

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &grpc_go.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *grpc_go.RegisterRequest) (*grpc_go.RegisterResponse, error) {

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &grpc_go.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *grpc_go.IsAdminRequest) (*grpc_go.IsAdminResponse, error) {

	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	return &grpc_go.IsAdminResponse{IsAdmin: true}, nil

}
