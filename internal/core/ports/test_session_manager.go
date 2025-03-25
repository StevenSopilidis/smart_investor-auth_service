package ports

import "gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"

type TestData struct {
}

// for testing because gomock does not support generic interfaces
type ITestSessionManager interface {
	CreateSession(user domain.User) (string, error)
	VerifySession(sessionId string) (*TestData, error)
}
