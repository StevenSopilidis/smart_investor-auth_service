package app_errors

type RedisError struct {
	msg string
}

func NewRedisError(err error) error {
	return &RedisError{
		msg: err.Error(),
	}
}

func (e *RedisError) Error() string {
	return "Could not create token: " + e.msg
}

func (e *RedisError) Is(target error) bool {
	_, ok := target.(*RedisError)
	return ok
}
