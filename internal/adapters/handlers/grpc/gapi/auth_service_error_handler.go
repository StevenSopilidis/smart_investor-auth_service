package gapi

import (
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getErrorResponseFromAuthServiceError(err error) error {
	switch e := err.(type) {
	case *app_errors.EmailNotVerified,
		*app_errors.InvalidPassword,
		*app_errors.InvalidSessionId,
		*app_errors.TokenDecryptionFailed,
		*app_errors.TokenExpired,
		*app_errors.UserNotFound:
		return status.Errorf(codes.Unauthenticated, e.Error())
	case *app_errors.ServiceUnreachable:
		return status.Errorf(codes.Unavailable, e.Error())
	default:
		return status.Errorf(codes.Internal, e.Error())
	}
}
