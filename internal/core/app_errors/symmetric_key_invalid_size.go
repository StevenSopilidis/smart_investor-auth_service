package app_errors

import (
	"fmt"

	"github.com/aead/chacha20poly1305"
)

type SymmetricKeyInvalidSize struct {
	validKeySize int
}

func NewSymmetricKeyInvalidSizeError(validKeySize int) error {
	return &SymmetricKeyInvalidSize{
		validKeySize: validKeySize,
	}
}

func (e *SymmetricKeyInvalidSize) Error() string {
	return fmt.Sprintf("Symmetric key must be of size: %d", chacha20poly1305.KeySize)
}

func (e *SymmetricKeyInvalidSize) Is(target error) bool {
	_, ok := target.(*SymmetricKeyInvalidSize)
	return ok
}
