package gapi

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/handlers/grpc/generated"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/services/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInvalidValidSessionReturnsUnauthenticated(t *testing.T) {
	userClientCtrl := gomock.NewController(t)
	userClientMock := mocks.NewMockIUserServiceClient(userClientCtrl)
	defer userClientCtrl.Finish()

	auth_service, config := getAuthService(userClientMock)
	server := NewServer(auth_service, config)

	_, err := server.VerifySession(context.Background(), &generated.VerifySessionRequest{
		Session: "invalid-sesssion",
	})
	require.Error(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unauthenticated, status.Code())
}
