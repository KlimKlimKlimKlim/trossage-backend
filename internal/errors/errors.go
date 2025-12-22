package errors

type Error struct {
	Message string
	Code    int
}

func New(status int, msg string) *Error {
	return &Error{
		Code:    status,
		Message: msg,
	}
}

func (e *Error) Error() string {
	return e.Message
}
