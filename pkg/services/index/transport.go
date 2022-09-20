package index

import (
	"errors"
	"net/http"
)

var errCannotApplyTemplate = errors.New("cannot apply template")
var errNotFound = errors.New("not found")

func encodeResponse(w http.ResponseWriter, statusCode int, response []byte) {
	w.WriteHeader(statusCode)
	_, _ = w.Write(response)
}

func encodeError(err error, w http.ResponseWriter) {
	switch err {
	case errCannotApplyTemplate:
		w.WriteHeader(http.StatusInternalServerError)
	case errNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, _ = w.Write([]byte(err.Error()))
}
