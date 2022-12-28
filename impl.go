package sos

var (
	_ Error = new(Err)
)

// Err is the concrete error value which implements the Error interface.
//
// The values are encapsulated in order to accurately obtain the
// runtime caller value and should be used an as builder.
type Err struct {
	reason  string
	code    Code
	message string
	err     error
	op      *op
	detail  map[string]string
}

// Code exposes the error Code value.
func (e *Err) Code() Code {
	return e.code
}

// Details exposes the error details map.
func (e *Err) Details() map[string]string {
	return e.detail
}

// Message exposes the most recent error message.
func (e *Err) Message() string {
	return e.message
}

// Operation exposes the most recent error Op value.
func (e *Err) Operation() Op {
	return e.op
}

// Reason exposes the most recent error reason code.
func (e *Err) Reason() string {
	return e.reason
}

// Error implements the error interface.
func (e *Err) Error() string {
	return trace(e)
}

func (e *Err) Unwrap() error {
	return e.err
}

// WithCode changes the error Code of the error value.
func (e *Err) WithCode(code Code) *Err {
	if e.message == FallbackMessage(e.code) {
		e.message = FallbackMessage(code)
	}
	if e.reason == string(e.code) {
		e.reason = string(code)
	}
	e.code = code
	return e
}

// WithError adds an error value to the error chain.
func (e *Err) WithError(err error) *Err {
	if v := As(err); v != nil {
		return e.propagate(v.err)
	}
	return e.propagate(err)
}

// WithMessage adds an error message.
func (e *Err) WithMessage(msg string, args ...interface{}) *Err {
	return e.propagate(message(sprintf(msg, args...)))
}

// WithReason adds an error reason code.
func (e *Err) WithReason(r string) *Err {
	return e.propagate(reason(r))
}

// WithDetail adds a single key-value error detail to the error details map.
func (e *Err) WithDetail(k string, v string) *Err {
	return e.propagate(details{k: v})
}

// WithDetails ...
func (e *Err) WithDetails(d map[string]string) *Err {
	return e.propagate(details(d))
}

// WithResetDetails empties the error detail map.
func (e *Err) WithResetDetails() *Err {
	e.detail = make(map[string]string)
	return e
}

// WithResetReason removes the error reason code.
func (e *Err) WithResetReason() *Err {
	e.reason = string(e.code)
	return e
}

type (
	message string
	reason  string
	details map[string]string
)

func (e *Err) propagate(args ...interface{}) *Err {

	if e.err != nil {
		p, ok := e.err.(*Err)
		if ok {

			for k, v := range p.detail {
				if _, ok := e.detail[k]; !ok {
					e.detail[k] = v
				}
			}

			if p.reason != "" {
				if p.code == e.code && e.reason == "" {
					e.reason = p.reason
				}
			}
		}

		if e.message == FallbackMessage(e.code) {
			if ok {
				e.message = p.message
			} else {
				e.message = e.err.Error()
			}
		}
	}

	for i := range args {
		switch v := args[i].(type) {
		case nil:
			continue
		case *op:
			e.op = v
		case reason:
			e.reason = string(v)
		case message:
			e.message = string(v)
		case details:
			if len(e.detail) == 0 {
				e.detail = v
			} else {
				for k, v := range v {
					e.detail[k] = v
				}
			}
		case Err:
			e.err = &v
		case *Err:
			if v == nil {
				continue
			}
			cp := *v
			e.err = &cp
		case error:
			e.err = v
			if !Is(v) && e.message == FallbackMessage(e.code) {
				e.message = v.Error()
			}
		}
	}

	return e
}
