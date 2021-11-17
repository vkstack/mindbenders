package errors

import "runtime"

type withMessage struct {
	msg   string
	cause error

	code interface{}
	// see the comments in https://stackoverflow.com/questions/40807281/is-there-any-performance-cost-in-using-runtime-caller
	// stack []uintptr
}

func callers() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

var DefaultSeparator string = "\n"

func New(msg string) error {
	return &withMessage{msg: msg}
}

func NewWithCode(msg string, code Code) error {
	return &withMessage{msg: msg, code: code}
}

func WrapCode(err error, code Code) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*withMessage); ok {
		return &withMessage{msg: err.Error(), cause: e, code: code}
	}
	return &withMessage{cause: err, code: code}
}

// func NewwithMessageFromNumbermsg(number int, msg string) error {
// 	return &withMessage{msg: msg, errs: []string{msg}, errorNumber: number}
// }

// Replace you errors.New("") to N
// func Wrap(err error, msg string) error {
// 	if e, ok := err.(*withMessage); ok {
// 		return &withMessage{msg: msg, cause: e}
// 	}
// 	return &withMessage{msg: msg, errs: []error{err}}
// }

// this will give entire error trace
// It is a thread unsafe method

func (e *withMessage) String() string {
	return e.msg
}

func (e *withMessage) Error() string {
	return e.msg + DefaultSeparator + e.cause.Error()
}

func (e *withMessage) Cause() error {
	return e.cause
}

func (e *withMessage) Code() Code {
	return e.code
}

func Cause(err error) error {
	if causer, ok := err.(interface{ Cause() error }); ok && causer.Cause() != nil {
		return Cause(causer.Cause())
	}
	return err
}
