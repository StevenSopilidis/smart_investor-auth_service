package app_errors

type CreateTokenOperationFailed struct {
	msg string
}

func NewCreateTokenOperationFailed(err error) error {
	return &CreateTokenOperationFailed{
		msg: err.Error(),
	}
}

func (e *CreateTokenOperationFailed) Error() string {
	return "Could not create token: " + e.msg
}

func (e *CreateTokenOperationFailed) Is(target error) bool {
	_, ok := target.(*CreateTokenOperationFailed)
	return ok
}
