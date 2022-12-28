package sos

import (
	"encoding/json"
	"net/http"
)

// HTTPStatusMap is initial mapping from sos.Code values to HTTP status codes.
var HTTPStatusMap = map[Code]int{
	CONFLICT:       http.StatusBadRequest,
	EXPIRED:        http.StatusBadRequest,
	FORBIDDEN:      http.StatusForbidden,
	INTERNAL:       http.StatusInternalServerError,
	INVALID:        http.StatusNotAcceptable,
	NOTFOUND:       http.StatusNotFound,
	NOTIMPLEMENTED: http.StatusNotImplemented,
	TEMPORARY:      http.StatusInternalServerError,
	TIMEOUT:        http.StatusRequestTimeout,
	UNAUTHORIZED:   http.StatusUnauthorized,
	UNPROCESSABLE:  http.StatusUnprocessableEntity,
}

// MarshalJSON implements the json.Marshaler interface.
func (e *Err) MarshalJSON() ([]byte, error) {
	if e == nil {
		return json.Marshal(nil)
	}

	v := struct {
		Code    Code              `json:"code"`
		Message string            `json:"message"`
		Reason  string            `json:"reason"`
		Details map[string]string `json:"details"`
	}{
		Code:    e.Code(),
		Message: e.Message(),
		Reason:  e.Reason(),
		Details: e.Details(),
	}

	return json.Marshal(v)
}
