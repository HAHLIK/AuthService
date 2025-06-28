package authgrpc

import (
	"context"

	ssov1 "github.com/HAHLIK/AuthService/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appId int64) (string, error)
	RegisterNewUser(ctx context.Context, email string, password string) (int64, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverAPI struct {
	auth Auth
	ssov1.UnimplementedAuthServer
}

func Register(grpc *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(grpc, &serverAPI{auth: auth})
}

const emptyValue = 0

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponce, error) {
	if err := s.validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internall error")
	}

	return &ssov1.LoginResponce{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponce, error) {
	if err := s.validateRegister(req); err != nil {
		return nil, err
	}

	usedId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internall error")
	}

	return &ssov1.RegisterResponce{
		UserId: usedId,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponce, error) {
	if req.GetUserId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "is_admin is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &ssov1.IsAdminResponce{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func (s *serverAPI) validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}
