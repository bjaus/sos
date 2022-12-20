package sos

import (
	"errors"
	"fmt"
	"strings"
)

func trace(e *Err) string {

	var t tracer

	t.add(e.message, e.code, e.op)

	w := errors.Unwrap(e)
	t.origin(w)

	for w != nil {
		if x := new(Err); errors.As(w, &x) {
			t.add(x.message, x.code, x.op)
		} else {
			t.origin(w)
		}

		w = errors.Unwrap(w)
	}

	return t.String()
}

type tracer struct {
	// o is the error of origin whchi is any non-nil error which does not implement the Error interface.
	o error
	// k is a slice of messages used to loop over the map in order.
	k []string
	// m is the message to file paths mapping used to produce the trace.
	m map[string][]string
	// s is a set to for deduplication.
	s map[string]struct{}
}

func (t *tracer) add(m string, code Code, op Op) {
	// Avoid unalloc panic
	if t.s == nil {
		t.s = make(map[string]struct{})
	}
	if t.m == nil {
		t.m = make(map[string][]string)
	}

	// Create the message key.
	k := fmt.Sprintf("[%s] %s", code, m)

	// Creat the filepath with line number.
	p := fmt.Sprintf("%s:%d", op.File(), op.Line())

	if _, ok := t.s[p]; !ok {
		t.m[k] = append(t.m[k], p)
		t.s[p] = struct{}{}
	}

	if _, ok := t.s[k]; !ok {
		t.k = append(t.k, k)
		t.s[k] = struct{}{}
	}
}

func (t *tracer) origin(err error) {
	if err == nil {
		return
	} else if x := As(err); x == nil {
		t.o = err
	}
}

func (t *tracer) String() string {

	// Loop backwards over all data collected to order properly.

	var b strings.Builder

	for i := len(t.k) - 1; i >= 0; i-- {
		msg := t.k[i]
		fmt.Fprint(&b, msg)

		for i := len(t.m[msg]) - 1; i >= 0; i-- {
			path := t.m[msg][i]
			fmt.Fprintf(&b, "\n\t%s", path)
		}

		fmt.Fprint(&b, "\n")
	}

	s := strings.TrimSpace(b.String())

	if t.o != nil && !strings.Contains(s, t.o.Error()) {
		return fmt.Sprintf("%s\n%s", t.o.Error(), s)
	}

	return s
}
