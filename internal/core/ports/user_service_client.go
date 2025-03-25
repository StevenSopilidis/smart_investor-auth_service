package ports

import "gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"

type IUserServiceClient interface {
	FindUserByEmail(email string) (domain.User, error)
}
