package services

import (
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/ports"
)

type AuthService[T any] struct {
	sessionManager              ports.ISessionManager[T]
	passwordVerificationService ports.IPasswordVerificationService
	userServiceClient           ports.IUserServiceClient
}

func NewAuthService[T any](
	sessionManager ports.ISessionManager[T],
	passwordVerificationService ports.IPasswordVerificationService,
	userServiceClient ports.IUserServiceClient,
) *AuthService[T] {
	return &AuthService[T]{
		sessionManager:              sessionManager,
		passwordVerificationService: passwordVerificationService,
		userServiceClient:           userServiceClient,
	}
}

func (s *AuthService[T]) Login(email string, password string) (string, error) {
	user, err := s.userServiceClient.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	if !user.EmaiLVerified {
		return "", &app_errors.EmailNotVerified{}
	}

	err = s.passwordVerificationService.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	session, err := s.sessionManager.CreateSession(user)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *AuthService[T]) ValidateSession(session string) (*T, error) {
	data, err := s.sessionManager.VerifySession(session)
	if err != nil {
		return nil, err
	}

	return data, nil
}
