package errors

import (
	"fmt"
	"strings"
)

type MultiError interface {
	error
	AddErrors(errs ...error) MultiError
	IsNil() bool
	Err() error
}

type multierror []error

func DefaultMultiError() MultiError { return &multierror{} }

func NewMultiError(errs ...error) MultiError {
	var me multierror
	me.AddErrors(errs...)
	if me.IsNil() {
		return nil
	}
	return &me
}

func (e *multierror) Error() string {
	if e == nil {
		return "nil"
	}
	var msg string
	for _, _e := range *e {
		if e != nil {
			msg += ("\n" + _e.Error())
		}
	}
	return strings.Trim(msg, "\n")
}

func (e *multierror) String() string {
	if e == nil {
		return "nil"
	}
	var msg string
	for _, _e := range *e {
		if e == nil {
			continue
		}
		if str, ok := _e.(fmt.Stringer); ok {
			msg += ("\n" + str.String())
		} else {
			msg += ("\n" + _e.Error())
		}
	}
	return strings.Trim(msg, "\n")
}

func (e *multierror) AddErrors(errs ...error) MultiError {
	if e != nil {
		for _, err := range errs {
			if err != nil {
				*e = append(*e, err)
			}
		}
	}
	return e
}

func (e *multierror) IsNil() bool {
	return e == nil || len(*e) == 0
}

func (e *multierror) Err() error {
	if e.IsNil() {
		return nil
	}
	return e
}
