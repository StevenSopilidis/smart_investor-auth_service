package app_errors

type TokenExpired struct {
	msg string
}

func NewTokenExpiredError() error {
	return &TokenExpired{}
}

func (e *TokenExpired) Error() string {
	return "Token has expired"
}

func (e *TokenExpired) Is(target error) bool {
	_, ok := target.(*TokenExpired)
	return ok
}
