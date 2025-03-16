package app_errors

type InvalidSessionId struct {
	id string
}

func NewInvalidSessionIdError(id string) error {
	return &InvalidSessionId{
		id: id,
	}
}

func (e *InvalidSessionId) Error() string {
	return "Invalid session id provided: " + e.id
}

func (e *InvalidSessionId) Is(target error) bool {
	_, ok := target.(*InvalidSessionId)
	return ok
}
