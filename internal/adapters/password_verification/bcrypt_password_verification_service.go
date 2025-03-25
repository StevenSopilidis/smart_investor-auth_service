package password_verification

import (
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHashService struct{}

func NewBcryptPasswordHashService() *BcryptPasswordHashService {
	return &BcryptPasswordHashService{}
}

func (bs *BcryptPasswordHashService) VerifyPassword(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return &app_errors.InvalidPassword{}
	}

	return nil
}
