package errors

import "fmt"

type multierror []error

func (e *multierror) Error() string {
	if e == nil {
		return "nil"
	}
	var msg string
	for _, _e := range *e {
		msg += ("\n" + _e.Error())
	}
	return msg
}

func (e *multierror) String() string {
	if e == nil {
		return "nil"
	}
	var msg string
	for _, _e := range *e {
		if str, ok := _e.(fmt.Stringer); ok {
			msg += ("\n" + str.String())
		} else {
			msg += ("\n" + _e.Error())
		}
	}
	return msg
}

func NewMultiError(errs ...error) error {
	e := multierror(errs)
	return &e
}
