package app_errors

type EmailNotVerified struct{}

func (e *EmailNotVerified) Error() string {
	return "Service unreachable"
}

func (e *EmailNotVerified) Is(target error) bool {
	_, ok := target.(*EmailNotVerified)
	return ok
}
