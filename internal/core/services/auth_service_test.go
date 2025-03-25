package services

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/ports"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/services/mocks"
)

type SessionData struct{}

func TestLogin(t *testing.T) {
	testCases := []struct {
		testName   string
		email      string
		password   string
		buildStubs func(
			mockSessionManager *mocks.MockITestSessionManager,
			mockPasswordVerificationService *mocks.MockIPasswordVerificationService,
			mockUserServiceClient *mocks.MockIUserServiceClient,
		)
		checkResponse func(err error)
	}{
		{
			testName: "OK",
			email:    "test@test.com",
			password: "pass1234",
			buildStubs: func(
				mockSessionManager *mocks.MockITestSessionManager,
				mockPasswordVerificationService *mocks.MockIPasswordVerificationService,
				mockUserServiceClient *mocks.MockIUserServiceClient) {
				mockUserServiceClient.EXPECT().FindUserByEmail("test@test.com").
					Times(1).
					Return(domain.User{EmaiLVerified: true}, nil)
				mockPasswordVerificationService.EXPECT().VerifyPassword("pass1234", gomock.Any()).
					Times(1).
					Return(nil)
				mockSessionManager.EXPECT().CreateSession(gomock.Any()).
					Times(1).
					Return("", nil)
			},
			checkResponse: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			testName: "EmailNotVerifiedReturnsError",
			email:    "test@test.com",
			password: "pass1234",
			buildStubs: func(
				mockSessionManager *mocks.MockITestSessionManager,
				mockPasswordVerificationService *mocks.MockIPasswordVerificationService,
				mockUserServiceClient *mocks.MockIUserServiceClient) {
				mockUserServiceClient.EXPECT().FindUserByEmail("test@test.com").
					Times(1).
					Return(domain.User{EmaiLVerified: false}, nil)
				mockPasswordVerificationService.EXPECT().VerifyPassword(gomock.Any(), gomock.Any()).
					Times(0)
				mockSessionManager.EXPECT().CreateSession(gomock.Any()).
					Times(0)
			},
			checkResponse: func(err error) {
				require.Error(t, err)
			},
		},
		{
			testName: "CouldNotCreateSessionReturnError",
			email:    "test@test.com",
			password: "pass1234",
			buildStubs: func(
				mockSessionManager *mocks.MockITestSessionManager,
				mockPasswordVerificationService *mocks.MockIPasswordVerificationService,
				mockUserServiceClient *mocks.MockIUserServiceClient) {
				mockUserServiceClient.EXPECT().FindUserByEmail("test@test.com").
					Times(1).
					Return(domain.User{EmaiLVerified: true}, nil)
				mockPasswordVerificationService.EXPECT().VerifyPassword("pass1234", gomock.Any()).
					Times(1).
					Return(nil)
				mockSessionManager.EXPECT().CreateSession(gomock.Any()).
					Times(1).
					Return("", fmt.Errorf("Error"))
			},
			checkResponse: func(err error) {
				require.Error(t, err)
			},
		},
		{
			testName: "InvalidPasswordReturnsError",
			email:    "test@test.com",
			password: "pass1234",
			buildStubs: func(
				mockSessionManager *mocks.MockITestSessionManager,
				mockPasswordVerificationService *mocks.MockIPasswordVerificationService,
				mockUserServiceClient *mocks.MockIUserServiceClient) {
				mockUserServiceClient.EXPECT().FindUserByEmail("test@test.com").
					Times(1).
					Return(domain.User{EmaiLVerified: true}, nil)
				mockPasswordVerificationService.EXPECT().VerifyPassword("pass1234", gomock.Any()).
					Times(1).
					Return(&app_errors.InvalidPassword{})
				mockSessionManager.EXPECT().CreateSession(gomock.Any()).
					Times(0)
			},
			checkResponse: func(err error) {
				require.ErrorIs(t, err, &app_errors.InvalidPassword{})
			},
		},
		{
			testName: "UserNotFoundReturnsError",
			email:    "test@test.com",
			password: "pass1234",
			buildStubs: func(
				mockSessionManager *mocks.MockITestSessionManager,
				mockPasswordVerificationService *mocks.MockIPasswordVerificationService,
				mockUserServiceClient *mocks.MockIUserServiceClient) {
				mockUserServiceClient.EXPECT().FindUserByEmail("test@test.com").
					Times(1).
					Return(domain.User{}, &app_errors.UserNotFound{})
				mockPasswordVerificationService.EXPECT().VerifyPassword(gomock.Any(), gomock.Any()).
					Times(0)
				mockSessionManager.EXPECT().CreateSession(gomock.Any()).
					Times(0)
			},
			checkResponse: func(err error) {
				require.ErrorIs(t, err, &app_errors.UserNotFound{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			sessionManagerCtrl := gomock.NewController(t)
			mockSessionManager := mocks.NewMockITestSessionManager(sessionManagerCtrl)
			defer sessionManagerCtrl.Finish()

			passwordCtrl := gomock.NewController(t)
			mockPassword := mocks.NewMockIPasswordVerificationService(passwordCtrl)
			defer passwordCtrl.Finish()

			userServiceClientCtrl := gomock.NewController(t)
			mockUserServiceClient := mocks.NewMockIUserServiceClient(userServiceClientCtrl)
			defer userServiceClientCtrl.Finish()

			tc.buildStubs(mockSessionManager, mockPassword, mockUserServiceClient)

			service := NewAuthService(mockSessionManager, mockPassword, mockUserServiceClient)
			_, err := service.Login(tc.email, tc.password)
			tc.checkResponse(err)
		})
	}
}

func ValidateSession(t *testing.T) {
	testCases := []struct {
		testName   string
		sessionId  string
		buildStubs func(
			mockSessionManager *mocks.MockITestSessionManager,
		)
		checkResponse func(err error)
	}{
		{
			testName:  "OK",
			sessionId: "id",
			buildStubs: func(mockSessionManager *mocks.MockITestSessionManager) {
				mockSessionManager.EXPECT().VerifySession("id").
					Times(1).
					Return(ports.TestData{}, nil)
			},
			checkResponse: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			testName:  "SessionManageReturnError",
			sessionId: "id",
			buildStubs: func(mockSessionManager *mocks.MockITestSessionManager) {
				mockSessionManager.EXPECT().VerifySession("id").
					Times(1).
					Return(ports.TestData{}, fmt.Errorf("error"))
			},
			checkResponse: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			sessionManagerCtrl := gomock.NewController(t)
			mockSessionManager := mocks.NewMockITestSessionManager(sessionManagerCtrl)
			defer sessionManagerCtrl.Finish()

			passwordCtrl := gomock.NewController(t)
			mockPassword := mocks.NewMockIPasswordVerificationService(passwordCtrl)
			defer passwordCtrl.Finish()

			userServiceClientCtrl := gomock.NewController(t)
			mockUserServiceClient := mocks.NewMockIUserServiceClient(userServiceClientCtrl)
			defer userServiceClientCtrl.Finish()

			tc.buildStubs(mockSessionManager)

			service := NewAuthService(mockSessionManager, mockPassword, mockUserServiceClient)
			_, err := service.ValidateSession(tc.sessionId)
			tc.checkResponse(err)
		})
	}
}
