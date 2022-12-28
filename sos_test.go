package sos_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/bjaus/sos"
	"github.com/google/go-cmp/cmp"
)

func TestError(t *testing.T) {

	// NOTE: This test relies on the package name and function name below.
	//       If either of these should change so should the vars below.
	pkg := "sos_test"
	caller := "TestError"

	var err error

	err = sos.New(sos.INTERNAL)
	if e := sos.As(err); e != nil {
		e.WithMessage("error message")
		e.WithError(fmt.Errorf("external error"))
		e.WithCode(sos.NOTIMPLEMENTED)
		e.WithReason("hello-world")
		err = e
	}

	err = sos.Trace(err)

	code := sos.CONFLICT
	reason := "testing"
	message := "new error message"
	details := map[string]string{
		"hello":   "world",
		"testing": "123",
	}

	if e := sos.As(err); e != nil {
		var k, v string = "something", "else"
		err = e.WithCode(code).
			WithError(err).
			WithMessage(message).
			WithReason(reason).
			WithDetail(k, v).
			WithDetails(details)

		details[k] = v

		if v := sos.As(err); v != nil {
			if diff := cmp.Diff(v.Details(), details); diff != "" {
				t.Error(diff)
			}
		} else {
			t.Fatal("not sos error")
		}
	}

	err = sos.Trace(err)
	_, file, line, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller not okay")
	}

	if e := sos.As(err); e != nil {
		t.Log(e.Error())
		if e.Message() != message {
			t.Errorf("message: got %s, want %s", e.Message(), message)
		}
		if kind := sos.Kind(err); kind != code {
			t.Errorf("code: got %s, want %s", e.Code(), code)
		}
		if e.Reason() != reason {
			t.Errorf("reason: got %s, want %s", e.Reason(), reason)
		}
		if diff := cmp.Diff(e.Details(), details); diff != "" {
			t.Error(diff)
		}
		op := e.Operation()
		if op.Package() != pkg {
			t.Errorf("op: pkg: got %s, want %s", op.Package(), pkg)
		}
		if op.Caller() != caller {
			t.Errorf("op: caller: got %s, want %s", op.Caller(), caller)
		}
		if op.File() != file {
			t.Errorf("op: file: got %s, want %s", op.File(), file)
		}
		lineno := line - 1
		if op.Line() != lineno {
			t.Errorf("op: line: got %d, want %d", op.Line(), lineno)
		}
		opstr := fmt.Sprintf("%s:%d", file, lineno)
		if op.String() != opstr {
			t.Errorf("op: got %s, want %s", op.String(), opstr)
		}
	} else {
		t.Fatal("not sos error")
	}
}

func TestNew(t *testing.T) {

	// NOTE: This test is sequential on purpose. Changing the order could break the test.
	//       It is intended to display error changes using the builder. Change with care :)

	/* ======================================================================
	Check Default on New
	====================================================================== */

	code := sos.NOTIMPLEMENTED // Initial error code.
	err := sos.New(code)
	if err.Code() != code {
		t.Errorf("code: got %q, want %q", err.Code(), code)
	}
	// Fallback message should change with the code since no message as been set.
	if msg := sos.FallbackMessage(code); err.Message() != msg {
		t.Errorf("message: got %q, want %q", err.Message(), msg)
	}
	if err.Operation() == nil {
		t.Fatalf("operation should be registered")
	}
	if err.Details() == nil {
		t.Fatalf("details map should be allocated")
	}
	if err.Reason() != string(code) {
		t.Fatalf("reason: got %q, want %q", err.Reason(), code)
	}

	/* ======================================================================
	Change Error Code
	====================================================================== */

	code = sos.NOTFOUND // Change the error code.
	err = err.WithCode(code)
	if err.Code() != code {
		t.Errorf("code: got %q, want %q", err.Code(), code)
	}
	// Fallback message should change with the code since no message as been set.
	if msg := sos.FallbackMessage(code); err.Message() != msg {
		t.Errorf("message: got %q, want %q", err.Message(), msg)
	}

	/* ======================================================================
	Check Error Message
	====================================================================== */

	msg := "this is the error message I expect"
	code = sos.FORBIDDEN // Change the error code.
	err = err.WithMessage(msg).WithCode(code)
	if err.Code() != code {
		t.Errorf("code: got %q, want %q", err.Code(), code)
	}
	// Fallback message should not effect changes when message is explicitly set.
	if err.Message() != msg {
		t.Errorf("message: got %q, want %q", err.Message(), msg)
	}

	/* ======================================================================
	Add Error Reason
	====================================================================== */

	reason := "testing"
	code = sos.INTERNAL
	err = err.WithReason(reason).WithCode(code)
	if err.Code() != code {
		t.Errorf("code: got %q, want %q", err.Code(), code)
	}
	// The reason should not change with the code because it has been explicitly set.
	if err.Reason() != reason {
		t.Errorf("reason: got %q, want %q", err.Reason(), reason)
	}
	err = err.WithResetReason()
	reason = string(code)
	// After reason reset it should go back to the default which is based on the error code.
	if err.Reason() != reason {
		t.Errorf("reason: got %q, want %q", err.Reason(), reason)
	}

	/* ======================================================================
	Add Single Detail
	====================================================================== */
	var k, v string = "hello", "world"
	err = err.WithDetail(k, v)
	if got := err.Details()[k]; got != v {
		t.Errorf("detail: key %s: got %s, want %s", k, got, v)
	}
	// err = err.WithResetDetails()
	// if len(err.Details()) > 0 {
	//     t.Fatal("did not reset details map")
	// }

	/* ======================================================================
	Add Multiple Details
	====================================================================== */
	deets := map[string]string{
		"hello":   "world",
		"testing": "123",
	}
	err = err.WithDetails(deets)
	if diff := cmp.Diff(err.Details(), deets); diff != "" {
		t.Error(diff)
	}
	err = err.WithResetDetails()
	if len(err.Details()) > 0 {
		t.Fatal("did not reset details map")
	}

	/* ======================================================================
	Check Error Trace
	====================================================================== */

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller failed")
	}
	if !strings.Contains(err.Error(), file) {
		t.Errorf("error trace did not include file: %s", file)
		t.Errorf("got %s", err.Error())
	}

}

func TestTrace(t *testing.T) {

	t.Run("from sos error", func(t *testing.T) {
		code := sos.NOTFOUND
		message := sos.FallbackMessage(code)
		reason := string(code)

		var err error = sos.New(code)
		err = sos.Trace(err)
		e := sos.As(err)
		if e == nil {
			t.Fatal("should not be nil error")
		}

		if e.Code() != code {
			t.Errorf("code: got %q, want %q", e.Code(), code)
		}
		if e.Reason() != reason {
			t.Errorf("reason: got %q, want %q", e.Reason(), reason)
		}
		if e.Operation() == nil {
			t.Error("operation: should be present")
		}
		if e.Details() == nil {
			t.Error("details: should be allocated")
		}

		err = sos.Trace(err)
		if e := sos.As(err); err != nil {
			// Prove that we can change the characteristic of an error.
			code = sos.EXPIRED
			message = "changed error code to expired"
			reason = "testing"
			err = e.WithCode(code).WithMessage(message).WithReason(reason)
		}

		e = sos.As(err)
		if e == nil {
			t.Fatal("should not be nil error")
		}

		if e.Code() != code {
			t.Errorf("code: got %q, want %q", e.Code(), code)
		}
		if e.Message() != message {
			t.Errorf("message: got %q, want %q", e.Message(), message)
		}
		if e.Reason() != reason {
			t.Errorf("reason: got %q, want %q", e.Reason(), reason)
		}
		if e.Operation() == nil {
			t.Error("operation: should be present")
		}
		if e.Details() == nil {
			t.Error("details: should be allocated")
		}
	})

	t.Run("from external error", func(t *testing.T) {
		err := fmt.Errorf("not sos error")

		err = sos.Trace(err)

		e := sos.As(err)
		if e == nil {
			t.Fatal("should not be nil error")
		}

		code := sos.INTERNAL
		if e.Code() != code {
			t.Errorf("code: got %q, want %q", e.Code(), code)
		}
		if e.Operation() == nil {
			t.Error("operation: should be present")
		}
		if e.Details() == nil {
			t.Error("details: should be allocated")
		}
	})

	t.Run("from nil error", func(t *testing.T) {
		var err error

		err = sos.Trace(err)
		if err != nil {
			t.Log("trace on nil should remain nil")
		}

		err = sos.Trace(sos.Error(nil))
		if err != nil {
			t.Log("trace on nil should remain nil")
		}

		var impl *sos.Err
		err = sos.Trace(impl)
		if err != nil {
			t.Log("trace on nil should remain nil")
		}

		impl = new(sos.Err)
		err = sos.Trace(impl)
		if err != nil {
			t.Log("trace on nil should remain nil")
		}
	})
}

func TestIs(t *testing.T) {

	type testcase struct {
		err  error
		want bool
	}

	cases := map[string]testcase{
		"nil": {
			err:  nil,
			want: false,
		},
		"error": {
			err:  fmt.Errorf("external error"),
			want: false,
		},
		"test error": {
			err:  &testerror{},
			want: false,
		},
		"custom": {
			err:  sos.New(sos.Code("custom")),
			want: true,
		},
	}

	for _, code := range sos.Codes {
		cases[string(code)] = testcase{
			err:  sos.New(code),
			want: true,
		}
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := sos.Is(tc.err)
			if got != tc.want {
				t.Errorf("%s: got %t, want %t", name, got, tc.want)
			}
		})
	}
}

func TestAs(t *testing.T) {

	type testcase struct {
		err  error
		want bool
	}

	cases := map[string]testcase{
		"nil": {
			err:  nil,
			want: false,
		},
		"error": {
			err:  fmt.Errorf("external error"),
			want: false,
		},
		"test error": {
			err:  &testerror{},
			want: false,
		},
		"custom": {
			err:  sos.New(sos.Code("custom")),
			want: true,
		},
	}

	for _, code := range sos.Codes {
		cases[string(code)] = testcase{
			err:  sos.New(code),
			want: true,
		}
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := sos.As(tc.err)
			if got := (err != nil); got != tc.want {
				t.Errorf("%s: got %t, want %t", name, got, tc.want)
			}
		})
	}
}

func TestMust(t *testing.T) {

	type testcase struct {
		err  error
		want bool
	}

	cases := map[string]testcase{
		"nil": {
			err:  nil,
			want: true,
		},
		"error": {
			err:  fmt.Errorf("external error"),
			want: true,
		},
		"test error": {
			err:  &testerror{},
			want: true,
		},
		"custom": {
			err:  sos.New(sos.Code("custom")),
			want: false,
		},
	}

	for _, code := range sos.Codes {
		cases[string(code)] = testcase{
			err:  sos.New(code),
			want: false,
		}
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if (recover() != nil) != tc.want {
					t.Fatalf("%s: panic state mismatch: got %t, want %t", name, !tc.want, tc.want)
				}
			}()
			sos.Must(tc.err)
		})
	}
}

func TestKind(t *testing.T) {

	type testcase struct {
		err  error
		want sos.Code
	}

	cases := map[string]testcase{
		"nil": {
			err:  nil,
			want: sos.Code(""),
		},
		"error": {
			err:  fmt.Errorf("external error"),
			want: sos.Code(""),
		},
		"custom": {
			err:  sos.New(sos.Code("custom")),
			want: sos.Code("custom"),
		},
	}

	for _, code := range sos.Codes {
		cases[string(code)] = testcase{
			err:  sos.New(code),
			want: code,
		}
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := sos.Kind(tc.err)
			if got != tc.want {
				t.Errorf("%s: got %q, want %q", name, got, tc.want)
			}
		})
	}
}

type testerror struct{}

func (ce testerror) Error() string {
	return "custom error"
}
