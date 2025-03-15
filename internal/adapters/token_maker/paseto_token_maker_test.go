package token_maker

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"
)

var pasetoMaker *PasetoTokenMaker

func TestMain(m *testing.M) {
	maker, err := NewPasetoTokenMaker("47956164296358211962664839075870")
	if err != nil {
		log.Fatalf("---> %s\n", err.Error())
	}
	pasetoMaker = maker
}

func ValidTokenVerificationReturnsClaims(t *testing.T) {
	user := domain.User{
		Email: "test@test.com",
	}

	token, err := pasetoMaker.GenerateToken(user, time.Minute)
	require.NoError(t, err)

	claims, err := pasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, user.Email, claims.Email)
}

func ExpiredTokenVerificationReturnsError(t *testing.T) {
	user := domain.User{
		Email: "test@test.com",
	}

	token, err := pasetoMaker.GenerateToken(user, -time.Minute)
	claims, err := pasetoMaker.VerifyToken(token)
	require.ErrorIs(t, err, &app_errors.TokenExpired{})
	require.Nil(t, claims)
}
