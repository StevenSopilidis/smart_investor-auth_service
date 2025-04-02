package gapi

import (
	"context"

	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/handlers/grpc/generated"
)

func (s *Server) Login(ctx context.Context, req *generated.LoginRequest) (
	*generated.LoginResponse, error,
) {
	session, err := s.auth_service.Login(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, getErrorResponseFromAuthServiceError(err)
	}

	return &generated.LoginResponse{
		Session: session,
	}, nil
}
