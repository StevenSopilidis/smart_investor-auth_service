package app_errors

type SymmetricKeyInvalidSize struct {
	validKeySize int
}

func NewSymmetricKeyInvalidSizeError(validKeySize int) error {
	return &SymmetricKeyInvalidSize{
		validKeySize: validKeySize,
	}
}

func (e *SymmetricKeyInvalidSize) Error() string {
	return "Symmetric key must be of size"
}

func (e *SymmetricKeyInvalidSize) Is(target error) bool {
	_, ok := target.(*SymmetricKeyInvalidSize)
	return ok
}
