package errors

type Error struct {
	Message string
	Code    int
}

func New(msg string, status int) *Error {
	return &Error{
		Code:    status,
		Message: msg,
	}
}

func (e *Error) Error() string {
	return e.Message
}
