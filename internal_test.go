package sos

import (
	"testing"
)

func TestSprintf(t *testing.T) {

	cases := map[string]struct {
		fmt  string
		args []interface{}
		want string
	}{
		"empty": {
			fmt:  "",
			args: []interface{}{1, 2, 3},
			want: "",
		},
		"no args": {
			fmt:  "testing 123",
			want: "testing 123",
		},
		"happy path": {
			fmt:  "testing %s %s",
			args: []interface{}{"hello", "world"},
			want: "testing hello world",
		},
		"too many args": {
			fmt:  "testing %d",
			args: []interface{}{123, "hello", "world"},
			want: "testing 123",
		},
		"too few args": {
			fmt:  "testing %d %s %s %s %s %s",
			args: []interface{}{123, "hello", "world"},
			want: "testing 123 hello world...",
		},
		"with literal %": {
			fmt:  "testing %d%%",
			args: []interface{}{100},
			want: "testing 100%",
		},
		"complex %% with too few args": {
			fmt:  "testing %d %% %d",
			args: []interface{}{5, 100, "hello", "world"},
			want: "testing 5 % 100",
		},
		"just literal %": {
			fmt:  "%%",
			args: []interface{}{5, 100, "hello", "world"},
			want: "%",
		},
		"invalid with string": {
			fmt:  "testing %d",
			args: []interface{}{"hello world"},
			want: "testing hello world",
		},
		"invalid with int": {
			fmt:  "testing %t",
			args: []interface{}{123},
			want: "testing 123",
		},
		"invalid with float": {
			fmt:  "testing %t",
			args: []interface{}{1.23},
			want: "testing 1.23",
		},
		"invalid with bool": {
			fmt:  "testing %f",
			args: []interface{}{true},
			want: "testing true",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := sprintf(tc.fmt, tc.args...)
			if got != tc.want {
				t.Errorf("%s: got %q, want %q", name, got, tc.want)
			}
		})
	}
}
