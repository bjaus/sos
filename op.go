package sos

import (
	"fmt"
	"runtime"
	"strings"
)

// Op provides details regarding the operation where the error originates.
type Op interface {
	// Package exposes the package which hosts the caller.
	Package() string
	// Func exposes the calling function or method.
	Caller() string
	// File exposes the file path.
	File() string
	// Line exposes the line number.
	Line() int

	fmt.Stringer
}

var (
	_ Op = new(op)
)

type op struct {
	pkg  string
	fn   string
	file string
	line int
}

func (o op) Package() string {
	return o.pkg
}

func (o op) Caller() string {
	return o.fn
}

func (o op) File() string {
	return o.file
}

func (o op) Line() int {
	return o.line
}

func (o op) String() string {
	return fmt.Sprintf("%s:%d", o.file, o.line)
}

var opReplacer = *strings.NewReplacer(
	"(*", "",
	")", "",
	".go", "",
)

func opParser(skip int) *op {
	o := op{
		pkg:  "unknown",
		fn:   "unknown",
		file: "unknown",
		line: -1,
	}

	if skip < 0 {
		skip = 0
	}
	skip++

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return &o
	}

	f := runtime.FuncForPC(pc)
	parts := strings.Split(f.Name(), "/")

	if len(parts) == 0 {
		return &o
	}

	caller := opReplacer.Replace(parts[len(parts)-1])
	parts = strings.Split(caller, ".func")

	caller = parts[0]
	parts = strings.Split(caller, ".")

	switch len(parts) {
	case 0:
		return &o
	case 2:
		o.fn = strings.Join(parts[1:], ".")
		fallthrough
	default:
		o.pkg = parts[0]
	}

	o.file = file
	o.line = line

	return &o
}
