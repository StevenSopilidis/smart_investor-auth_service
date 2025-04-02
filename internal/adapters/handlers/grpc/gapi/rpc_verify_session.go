package gapi

import (
	"context"

	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/handlers/grpc/generated"
)

func (s *Server) VerifySession(ctx context.Context, req *generated.VerifySessionRequest) (
	*generated.VerifySessionResponse, error,
) {
	data, err := s.auth_service.ValidateSession(req.GetSession())
	if err != nil {
		return nil, getErrorResponseFromAuthServiceError(err)
	}

	return &generated.VerifySessionResponse{
		Email: data.Email,
	}, nil
}
