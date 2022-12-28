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
	if o.file == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", o.file, o.line)
}

var opReplacer = *strings.NewReplacer(
	"(*", "",
	")", "",
	".go", "",
)

func opParser(skip int) *op {
	if skip < 0 {
		skip = 0
	}
	skip++ // Add one to skip this func.

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil
	}

	f := runtime.FuncForPC(pc)
	parts := strings.Split(f.Name(), "/")

	if len(parts) == 0 {
		return nil
	}

	// Clean up some of the cruft provided by the runtime package.
	caller := opReplacer.Replace(parts[len(parts)-1])
	parts = strings.Split(caller, ".func")
	caller = parts[0]
	parts = strings.Split(caller, ".")

	var pkg, fn string

	switch len(parts) {
	case 0:
		return nil
	case 2:
		fn = strings.Join(parts[1:], ".")
		fallthrough
	default:
		pkg = parts[0]
	}

	return &op{
		pkg:  pkg,
		fn:   fn,
		file: file,
		line: line,
	}
}
