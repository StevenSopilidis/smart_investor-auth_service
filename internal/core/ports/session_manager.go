package ports

import "gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"

type ISessionManager[T any] interface {
	CreateSession(user domain.User) (string, error)
	VerifySession(sessionId string) (*T, error)
}
