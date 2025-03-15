package token_maker

import (
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"
)

type Claims struct {
	Email     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type PasetoTokenMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoTokenMaker(symmetricKey string) (*PasetoTokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, app_errors.NewSymmetricKeyInvalidSizeError(chacha20poly1305.KeySize)
	}

	maker := &PasetoTokenMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (p *PasetoTokenMaker) GenerateToken(user domain.User, duration time.Duration) (string, error) {
	claims := &Claims{
		Email:     user.Email,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	token, err := p.paseto.Encrypt(p.symmetricKey, claims, nil)
	if err != nil {
		return "", app_errors.NewCreateTokenOperationFailed(err)
	}

	return token, nil
}

func (p *PasetoTokenMaker) VerifyToken(token string) (*Claims, error) {
	claims := &Claims{}

	err := p.paseto.Decrypt(token, p.symmetricKey, claims, nil)
	if err != nil {
		return nil, app_errors.NewTokenDecryptionFailedError(err)
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, app_errors.NewTokenExpiredError()
	}

	return claims, nil
}
