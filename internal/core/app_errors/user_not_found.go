package app_errors

type UserNotFound struct{}

func (e *UserNotFound) Error() string {
	return "User was not found"
}

func (e *UserNotFound) Is(target error) bool {
	_, ok := target.(*UserNotFound)
	return ok
}
