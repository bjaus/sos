package sos

// Code indicates the error code.
type Code string

// Default error codes.
const (
	// INTERNAL indicates an error caused by internal failure.
	INTERNAL Code = "internal"
	// CONFLICT indicates an error caused by a conflict (i.e., update conflict, etc.).
	CONFLICT Code = "conflict"
	// EXPIRED indicates an error caused by resource expiration.
	EXPIRED Code = "expired"
	// FORBIDDEN indicates an error caused by forbidden actions.
	FORBIDDEN Code = "forbidden"
	// INVALID indicates an error caused by invalid actions.
	INVALID Code = "invalid"
	// NOTFOUND indicates an error caused by a resource not being found.
	NOTFOUND Code = "not found"
	// NOTIMPLEMENTED indicates an error caused by functionality not being implemented.
	NOTIMPLEMENTED Code = "not implemented"
	// TEMPORARY indicates an error caused by a temporary issue.
	TEMPORARY Code = "temporary"
	// TIMEOUT indiciates an error caused by something timing out.
	TIMEOUT Code = "timeout"
	// UNAUTHORIZED indicates an error caused by an unauthorized actor.
	UNAUTHORIZED Code = "unauthorized"
	// UNPROCESSABLE indicates an error caused by input that can't be processed.
	UNPROCESSABLE Code = "unprocessable"
)

// Codes is a slice of all built-in error Code values.
var Codes = []Code{
	INTERNAL,
	CONFLICT,
	EXPIRED,
	FORBIDDEN,
	INVALID,
	NOTFOUND,
	NOTIMPLEMENTED,
	TEMPORARY,
	TIMEOUT,
	UNAUTHORIZED,
	UNPROCESSABLE,
}
