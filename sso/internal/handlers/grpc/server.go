package grpc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	// "fmt"

	"sso/internal/services/auth"
	"sso/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssov1 "github.com/GOeda-Co/proto-contract/gen/go/sso"
	"github.com/google/uuid"
	//use protos package
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
		name string,
	) (userID uuid.UUID, err error)
	IsAdmin(
		ctx context.Context,
		userId uuid.UUID,
	) (isAdmin bool, err error)
	RegisterApp(
		ctx context.Context,
		name string,
		secret string,
	) (appID int, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		// Ошибку auth.ErrInvalidCredentials мы создадим ниже
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to login: %v", err))
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	uid, err := s.auth.RegisterNewUser(ctx, in.Email, in.Password, in.Name)
	if err != nil {
		// Ошибку storage.ErrUserExists мы создадим ниже
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.RegisterResponse{UserId: string(uid.String())}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid uuid")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, userID)
	if err != nil {
		// Check if user not found
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to get user admin status")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func (s *serverAPI) RegisterApp(ctx context.Context, in *ssov1.RegisterAppRequest) (*ssov1.RegisterAppResponse, error) {
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if in.Secret == "" {
		return nil, status.Error(codes.InvalidArgument, "secret is required")
	}

	appID, err := s.auth.RegisterApp(ctx, in.Name, in.Secret)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register app")
	}

	return &ssov1.RegisterAppResponse{AppId: strconv.Itoa(appID)}, nil
}