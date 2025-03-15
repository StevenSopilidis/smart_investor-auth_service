package app_errors

type TokenDecryptionFailed struct {
	msg string
}

func NewTokenDecryptionFailedError(err error) error {
	return &TokenDecryptionFailed{
		msg: err.Error(),
	}
}

func (e *TokenDecryptionFailed) Error() string {
	return "Could decrypt token: " + e.msg
}

func (e *TokenDecryptionFailed) Is(target error) bool {
	_, ok := target.(*TokenDecryptionFailed)
	return ok
}
