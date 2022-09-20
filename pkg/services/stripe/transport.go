package stripe

import (
	"encoding/json"
	"errors"
	"net/http"
)

var errInvalidFormData = errors.New("invalid form data")
var errInvalidSession = errors.New("invalid session")
var errUserNotFound = errors.New("user not found")
var errIncorrectStripeSignature = errors.New("incorrect stripe signature")

func encodeResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func encodeError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case errInvalidFormData:
		w.WriteHeader(http.StatusBadRequest)
	case errInvalidSession:
		w.WriteHeader(http.StatusBadRequest)
	case errUserNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errIncorrectStripeSignature:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": err.Error(),
	})
}
