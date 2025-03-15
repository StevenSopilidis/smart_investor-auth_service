package ports

import (
	"time"

	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"
)

type ITokenMaker[T any] interface {
	GenerateToken(user domain.User, duration time.Duration) (string, error)
	VerifyToken(token string) (T, error)
}
