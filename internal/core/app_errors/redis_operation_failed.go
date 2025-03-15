package app_errors

type RedisOperationError struct {
	msg string
}

func NewRedisOperationError(err error) error {
	return &RedisOperationError{
		msg: err.Error(),
	}
}

func (e *RedisOperationError) Error() string {
	return "Could not create token: " + e.msg
}

func (e *RedisOperationError) Is(target error) bool {
	_, ok := target.(*RedisOperationError)
	return ok
}
