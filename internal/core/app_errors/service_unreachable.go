package app_errors

type ServiceUnreachable struct{}

func (e *ServiceUnreachable) Error() string {
	return "Service unreachable"
}

func (e *ServiceUnreachable) Is(target error) bool {
	_, ok := target.(*ServiceUnreachable)
	return ok
}
