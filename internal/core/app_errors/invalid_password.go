package app_errors

type InvalidPassword struct{}

func (e *InvalidPassword) Error() string {
	return "Invalid Password Provided"
}

func (e *InvalidPassword) Is(target error) bool {
	_, ok := target.(*InvalidPassword)
	return ok
}
