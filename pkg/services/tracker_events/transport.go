package tracker_events

import (
	"encoding/json"
	"errors"
	"habiko-go/pkg/models"
	"net/http"
)

var errInvalidParams = errors.New("invalid parameters provided")
var errNotFound = errors.New("tracker not found")
var errTrackerAccessDenied = errors.New("access to this tracker denied")
var errTooLateToAddTrackerEvent = errors.New("too late to add tracker event")
var errAccessDenied = errors.New("access to this tracker event denied")
var errMissingData = errors.New("no tracker data provided")
var errTypeIncorrect = errors.New("event type is too short")
var errUserAccessDenied = errors.New("user access denied")

type newTrackerEventRequest struct {
	*models.NewTrackerEvent
}

func (request *newTrackerEventRequest) Bind(r *http.Request) error {
	if request.NewTrackerEvent == nil {
		return errMissingData
	}
	return nil
}

type patchedTrackerEventRequest struct {
	*models.PatchedTrackerEvent
}

func (request *patchedTrackerEventRequest) Bind(r *http.Request) error {
	if request.PatchedTrackerEvent == nil {
		return errMissingData
	}
	return nil
}

type h map[string]interface{}

func encodeResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func encodeError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case models.ErrDb:
		w.WriteHeader(http.StatusInternalServerError)
	case models.ErrDbObject:
		w.WriteHeader(http.StatusInternalServerError)
	case models.ErrInvalidSession:
		w.WriteHeader(http.StatusForbidden)
	case errInvalidParams:
		w.WriteHeader(http.StatusBadRequest)
	case errNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errTrackerAccessDenied:
		w.WriteHeader(http.StatusForbidden)
	case errTooLateToAddTrackerEvent:
		w.WriteHeader(http.StatusForbidden)
	case errUserAccessDenied:
		w.WriteHeader(http.StatusForbidden)
	case errAccessDenied:
		w.WriteHeader(http.StatusForbidden)
	case errMissingData:
		w.WriteHeader(http.StatusBadRequest)
	case errTypeIncorrect:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": err.Error(),
	})
}
