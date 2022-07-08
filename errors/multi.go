package errors

import (
	"fmt"
	"strings"
)

type MultiError struct{ errs []error }

func DefaultMultiError() *MultiError { return &MultiError{} }

func NewMultiError(errs ...error) *MultiError {
	var me MultiError
	me.AddErrors(errs...)
	if me.IsNil() {
		return nil
	}
	return &me
}

func (e *MultiError) Error() string {
	if e == nil {
		return "nil"
	}
	var msg string
	for _, _e := range e.errs {
		if e != nil {
			msg += ("\n" + _e.Error())
		}
	}
	return strings.Trim(msg, "\n")
}

func (e *MultiError) String() string {
	if e == nil {
		return "nil"
	}
	var msg string
	for _, _e := range e.errs {
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

func (e *MultiError) AddErrors(errs ...error) *MultiError {
	if e != nil {
		for _, err := range errs {
			if err != nil {
				e.errs = append(e.errs, err)
			}
		}
	}
	return e
}

func (e *MultiError) IsNil() bool {
	return e == nil || len(e.errs) == 0
}

func (e *MultiError) Err() error {
	if e.IsNil() {
		return nil
	}
	return e
}
