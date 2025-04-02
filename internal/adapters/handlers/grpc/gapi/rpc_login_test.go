package gapi

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/config"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/handlers/grpc/generated"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/password_verification"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/session_manager"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/token_maker"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/ports"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/services"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/services/mocks"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func getAuthService(user_client ports.IUserServiceClient) (
	*services.AuthService[token_maker.Claims], config.Config,
) {
	config := config.Config{
		SymmetricKey:    "RmlXfdOwoj8mPUj4xH5VvVRKsCJRdQNZ",
		RedisAddress:    "localhost:6379",
		RedisPassword:   "",
		RedisDB:         0,
		TokenDuration:   time.Minute,
		UserServiceAddr: "localhost:8080",
	}

	token_maker, err := token_maker.NewPasetoTokenMaker(config.SymmetricKey)
	if err != nil {
		log.Fatal("---> Could not create token maker: ", err)
	}

	session_manager, err := session_manager.NewRedisSessionManager(config, token_maker)
	if err != nil {
		log.Fatal("---> Could not create session_manager: ", err)
	}

	password_service := password_verification.NewBcryptPasswordHashService()

	auth_service := services.NewAuthService(
		session_manager,
		password_service,
		user_client,
	)

	return auth_service, config
}

func TestValidLoginReturnsOK(t *testing.T) {
	userClientCtrl := gomock.NewController(t)
	userClientMock := mocks.NewMockIUserServiceClient(userClientCtrl)
	defer userClientCtrl.Finish()

	req := &generated.LoginRequest{
		Email:    "test@test.com",
		Password: "pass",
	}

	password_hash, err := hashPassword(req.GetPassword())
	require.NoError(t, err)

	userClientMock.EXPECT().FindUserByEmail(req.Email).Times(1).
		Return(
			domain.User{
				Email:         req.Email,
				EmaiLVerified: true,
				Password:      password_hash,
			}, nil,
		)

	auth_service, config := getAuthService(userClientMock)

	server := NewServer(auth_service, config)
	res, err := server.Login(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	claims, err := server.auth_service.ValidateSession(res.Session)
	require.NoError(t, err)
	require.Equal(t, req.Email, claims.Email)
}

func TestUserServiceUnreachableReturnsServiceUnavailable(t *testing.T) {
	userClientCtrl := gomock.NewController(t)
	userClientMock := mocks.NewMockIUserServiceClient(userClientCtrl)
	defer userClientCtrl.Finish()

	req := &generated.LoginRequest{
		Email:    "test@test.com",
		Password: "pass",
	}

	userClientMock.EXPECT().FindUserByEmail(req.Email).Times(1).
		Return(domain.User{}, &app_errors.ServiceUnreachable{})

	auth_service, config := getAuthService(userClientMock)

	server := NewServer(auth_service, config)
	_, err := server.Login(context.Background(), req)
	require.Error(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unavailable, status.Code())
}

func TestUserNotFoundReturnsUnathentiacted(t *testing.T) {
	userClientCtrl := gomock.NewController(t)
	userClientMock := mocks.NewMockIUserServiceClient(userClientCtrl)
	defer userClientCtrl.Finish()

	req := &generated.LoginRequest{
		Email:    "test@test.com",
		Password: "pass",
	}

	userClientMock.EXPECT().FindUserByEmail(req.Email).Times(1).
		Return(domain.User{}, &app_errors.UserNotFound{})

	auth_service, config := getAuthService(userClientMock)

	server := NewServer(auth_service, config)
	_, err := server.Login(context.Background(), req)
	require.Error(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unauthenticated, status.Code())
}

func TestEmailNotVerifiedReturnsUnathenticated(t *testing.T) {
	userClientCtrl := gomock.NewController(t)
	userClientMock := mocks.NewMockIUserServiceClient(userClientCtrl)
	defer userClientCtrl.Finish()

	req := &generated.LoginRequest{
		Email:    "test@test.com",
		Password: "pass",
	}

	password_hash, err := hashPassword(req.GetPassword())
	require.NoError(t, err)

	userClientMock.EXPECT().FindUserByEmail(req.Email).Times(1).
		Return(domain.User{
			Email:         req.Email,
			EmaiLVerified: false,
			Password:      password_hash,
		}, nil)

	auth_service, config := getAuthService(userClientMock)

	server := NewServer(auth_service, config)
	_, err = server.Login(context.Background(), req)
	require.Error(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unauthenticated, status.Code())
}

func TestInvalidPasswordReturnsUnauthenticated(t *testing.T) {
	userClientCtrl := gomock.NewController(t)
	userClientMock := mocks.NewMockIUserServiceClient(userClientCtrl)
	defer userClientCtrl.Finish()

	req := &generated.LoginRequest{
		Email:    "test@test.com",
		Password: "pass",
	}

	userClientMock.EXPECT().FindUserByEmail(req.Email).Times(1).
		Return(domain.User{
			Email:         req.Email,
			EmaiLVerified: true,
			Password:      "invalid",
		}, nil)

	auth_service, config := getAuthService(userClientMock)

	server := NewServer(auth_service, config)
	_, err := server.Login(context.Background(), req)
	require.Error(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unauthenticated, status.Code())
}
