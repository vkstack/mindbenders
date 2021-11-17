package errors

import "fmt"

type base struct {
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
	return &base{msg: msg}
}

func NewWithCode(msg string, code interface{}) error {
	return &base{msg: msg, code: code}
}

func WrapMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(interface{ Code() interface{} }); ok {
		return &base{msg: msg, cause: err, code: e.Code()}
	}
	return &base{msg: msg, cause: err}
}

func WrapMessageWithCode(err error, msg string, code interface{}) error {
	if err == nil {
		return nil
	}
	return &base{msg: msg, cause: err, code: code}
}

func (e *base) String() string {
	if e == nil {
		return ""
	}
	if e.cause == nil {
		return e.msg
	}
	return e.msg + DefaultSeparator + e.cause.Error()
}

func (e *base) Error() string {
	return e.msg
}

func (e *base) Cause() error {
	return e.cause
}

func (e *base) Code() interface{} {
	return e.code
}

/**
{
	"error":Cause(e),
	"errMsg":e.Error(),
}
*/

// e0->e1->e2->e3
// e0 and e3 is ground level and top level errors
// UnWrap(e3) -> e2
// Cause(e3) ->e0
//fmt.Println(e) -> print msg
//fmt.Println(e.Error()) ->prints msg-trace

func Cause(err error) error {
	fmt.Println(err)
	if causer, ok := err.(interface{ Cause() error }); ok && causer.Cause() != nil {
		return Cause(causer.Cause())
	}
	return err
}

func UnWrap(err error) error {
	if causer, ok := err.(interface{ Cause() error }); ok && causer.Cause() != nil {
		return causer.Cause()
	}
	return err
}
