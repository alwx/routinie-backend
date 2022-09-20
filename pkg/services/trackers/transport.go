package trackers

import (
	"encoding/json"
	"errors"
	"habiko-go/pkg/models"
	"net/http"
)

var errInvalidParams = errors.New("invalid parameters provided")
var errNotFound = errors.New("tracker not found")
var errAccessDenied = errors.New("access to this tracker denied")
var errMissingData = errors.New("no tracker data provided")
var errTitleIncorrect = errors.New("title is too short")
var errUserAccessDenied = errors.New("user access denied")
var errNoPremiumSubscription = errors.New("cannot do this without premium subscription")

type newTrackerRequest struct {
	*models.NewTracker
}

func (request *newTrackerRequest) Bind(r *http.Request) error {
	if request.NewTracker == nil {
		return errMissingData
	}
	if len(request.Title) < 3 {
		return errTitleIncorrect
	}
	return nil
}

type patchedTrackerRequest struct {
	*models.PatchedTracker
}

func (request *patchedTrackerRequest) Bind(r *http.Request) error {
	if request.PatchedTracker == nil {
		return errMissingData
	}
	if request.Title != nil && len(*request.Title) < 3 {
		return errTitleIncorrect
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
	case errAccessDenied:
		w.WriteHeader(http.StatusForbidden)
	case errNoPremiumSubscription:
		w.WriteHeader(http.StatusForbidden)
	case errMissingData:
		w.WriteHeader(http.StatusBadRequest)
	case errTitleIncorrect:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": err.Error(),
	})
}
