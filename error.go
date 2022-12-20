package sos

import "errors"

// Error ...
type Error interface {
	error
	Code() Code
	Details() map[string]string
	Message() string
	Operation() Op
	Reason() string
}

// New creates an new Err value for building out an error with desired details.
func New(code Code) *Err {
	return create(code, FallbackMessage(code), nil)
}

// Trace provides the ability to add trace the error without altering the error value.
//
// This is handy for debugging and log statements when the error is a pass-through but
// tracking the error through the application is desired.
//
// If the error provided is nil then the returned value is nil as well.
// And if the error provided does not satisfy the Error interface the Code will default to INTERNAL.
func Trace(err error) error {
	if err == nil {
		return nil
	}

	prev, ok := err.(*Err)
	if ok {
		if prev == nil || prev.op == nil {
			return nil
		}
		op := opParser(1)
		return prev.propagate(err, op)
	}

	return create(INTERNAL, err.Error(), err)
}

func create(code Code, msg string, err error) *Err {
	e := Err{
		code:    code,
		message: msg,
		reason:  string(code),
		op:      opParser(2),
		detail:  make(map[string]string),
	}

	if err != nil {
		e.err = err
	}

	return &e
}

// Is indicates whether the error provided implements the Error interface.
func Is(err error) bool {
	var i Error
	return errors.As(err, &i)
}

// As converts the error provided into an Err value if the error implements the Error interface.
//
// If the error does not implement the Error interface then the returned value is nil.
func As(err error) *Err {
	if e := new(Err); errors.As(err, &e) {
		return e
	}
	return nil
}

// Kind extracts the Code from the error provided if it satisfies the Error interface.
//
// If the error does not satisfy the Error interface or is nil then an empty Code value is returned..
func Kind(err error) Code {
	if err != nil {
		if e := Error(nil); errors.As(err, &e) {
			return e.Code()
		}
	}
	return ""
}
