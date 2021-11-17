package errors

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

func NewWithCode(msg string, code Code) error {
	return &base{msg: msg, code: code}
}

func WrapMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &base{msg: msg, cause: err}
}

func WrapMessageWithCode(err error, msg string, code Code) error {
	if err == nil {
		return nil
	}
	return &base{msg: msg, cause: err, code: code}
}

func (e *base) String() string {
	return e.msg
}

func (e *base) Error() string {
	if e == nil {
		return ""
	}
	if e.cause == nil {
		return e.msg
	}
	return e.msg + DefaultSeparator + e.cause.Error()
}

func (e *base) Cause() error {
	return e.cause
}

func (e *base) Code() Code {
	return e.code
}

func (e *base) UnWrap() error {
	return e.cause
}

func Cause(err error) error {
	if causer, ok := err.(interface{ Cause() error }); ok && causer.Cause() != nil {
		return Cause(causer.Cause())
	}
	return err
}
