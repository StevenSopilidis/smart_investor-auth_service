package session_manager

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/config"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/token_maker"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"
)

var sessionManager *RedisSessionManager[token_maker.Claims]

func TestMain(m *testing.M) {
	config := config.Config{
		RedisAddress:  "127.0.0.1:6379",
		RedisPassword: "",
		RedisDB:       0,
		SymmetricKey:  "47956164296358211962664839075870",
		TokenDuration: time.Hour,
	}

	tokenMaker, err := token_maker.NewPasetoTokenMaker(config.SymmetricKey)
	if err != nil {
		log.Fatal("---> Could not create token maker: ", err)
	}

	manager, err := NewRedisSessionManager(config, tokenMaker)
	if err != nil {
		log.Fatal("---> Could not create session manager: ", err)
	}

	sessionManager = manager
}

func ValidSessionReturnsClaims(t *testing.T) {
	email := "test@test.com"

	sessionId, err := sessionManager.CreateSession(domain.User{
		Email: email,
	})
	require.NoError(t, err)

	claims, err := sessionManager.VerifySession(sessionId)
	require.NoError(t, err)
	require.Equal(t, email, claims.Email)
}

func InvalidSessionIdReturnsInvalidSessionIdError(t *testing.T) {
	sessionId := "invalid-id"
	claims, err := sessionManager.VerifySession(sessionId)
	require.ErrorIs(t, err, &app_errors.InvalidSessionId{})
	require.Nil(t, claims)
}
