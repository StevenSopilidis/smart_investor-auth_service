package session_manager

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/adapters/config"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/app_errors"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/domain"
	"gitlab.com/stevensopi/smart_investor/auth_service/internal/core/ports"
)

type SessionPayload struct {
	Email string
}

type RedisSessionManager[T any] struct {
	redisClient *redis.Client
	tokenMaker  ports.ITokenMaker[T]
	config      config.Config
}

func NewRedisSessionManager[T any](
	config config.Config,
	tokenMaker ports.ITokenMaker[T],
) (*RedisSessionManager[T], error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, app_errors.NewRedisConnectionError(err)
	}

	return &RedisSessionManager[T]{
		redisClient: redisClient,
		tokenMaker:  tokenMaker,
		config:      config,
	}, nil
}

func (m *RedisSessionManager[T]) CreateSession(user domain.User) (string, error) {
	token, err := m.tokenMaker.GenerateToken(user, m.config.TokenDuration)
	if err != nil {
		return "", err
	}

	sessionId := uuid.New().String()
	err = m.redisClient.Set(context.Background(), sessionId, token, 0).Err()
	if err != nil {
		return "", app_errors.NewRedisOperationError(err)
	}

	return sessionId, nil
}

func (m *RedisSessionManager[T]) VerifySession(sessionId string) (*T, error) {
	val, err := m.redisClient.Get(context.Background(), sessionId).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, app_errors.NewInvalidSessionIdError(sessionId)
		} else {
			return nil, app_errors.NewRedisError(err)
		}
	}

	payload, err := m.tokenMaker.VerifyToken(val)
	if err != nil {
		return nil, err
	}

	return payload, err
}
