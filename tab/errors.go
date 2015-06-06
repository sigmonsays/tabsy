package tab

import (
	"fmt"
)

type Error struct {
	Message string
}

func NewError(message string) *Error {
	return &Error{
		Message: message,
	}
}

func (e *Error) Copy() *Error {
	e2 := *e
	return &e2
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Errorf(s string, args ...interface{}) *Error {
	res := e.Copy()
	extra := fmt.Sprintf(s, args...)
	res.Message = fmt.Sprintf("%s: %s", e.Message, extra)
	return res
}
