package app_errors

type RedisConnectionError struct {
	msg string
}

func NewRedisConnectionError(err error) error {
	return &RedisConnectionError{
		msg: err.Error(),
	}
}

func (e *RedisConnectionError) Error() string {
	return "Could not create token: " + e.msg
}

func (e *RedisConnectionError) Is(target error) bool {
	_, ok := target.(*RedisConnectionError)
	return ok
}
