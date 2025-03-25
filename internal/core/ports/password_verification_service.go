package ports

type IPasswordVerificationService interface {
	VerifyPassword(password string, hash string) error
}
