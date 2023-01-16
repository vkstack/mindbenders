package errors

import "fmt"

type causer interface {
	Cause() error
}

type BasicError struct {
	msg   string
	cause error

	code interface{}

	// see the comments in https://stackoverflow.com/questions/40807281/is-there-any-performance-cost-in-using-runtime-caller
	// stack []uintptr
}

// func callers() []uintptr {
// 	const depth = 32
// 	var pcs [depth]uintptr
// 	n := runtime.Callers(3, pcs[:])
// 	return pcs[0:n]
// }

var DefaultSeparator string = "\n"

func New(msg string) error {
	return &BasicError{msg: msg}
}

func NewWithError(err error) error {
	return &BasicError{msg: err.Error(), cause: err}
}

func NewWithCode(msg string, code interface{}) error {
	return &BasicError{msg: msg, code: code}
}

func WrapMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(interface{ Code() interface{} }); ok {
		return &BasicError{msg: msg, cause: err, code: e.Code()}
	}
	return &BasicError{msg: msg, cause: err}
}

func WrapMessageWithCode(err error, msg string, code interface{}) error {
	if err == nil {
		return nil
	}
	return &BasicError{msg: msg, cause: err, code: code}
}

func (e *BasicError) String() string {
	if e == nil {
		return ""
	}
	if e.cause == nil {
		return e.msg
	}
	if cause, ok := e.cause.(fmt.Stringer); ok {
		return e.msg + DefaultSeparator + cause.String()
	}
	return e.msg + DefaultSeparator + e.cause.Error()
}

func (e *BasicError) Error() string {
	return e.msg
}

func (e *BasicError) Cause() error {
	return e.cause
}

func (e *BasicError) Code() interface{} {
	return e.code
}

func Cause(err error) error {
	if causer, ok := err.(causer); ok && causer.Cause() != nil {
		return Cause(causer.Cause())
	}
	return err
}

func UnWrap(err error) error {
	if causer, ok := err.(causer); ok {
		return causer.Cause()
	}
	return err
}
