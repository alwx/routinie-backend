package tracker_events

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"habiko-go/pkg/models"
	"net/http"
	"strconv"
	"time"
)

func MakeHandler(service Service) http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// get session
		ctx := r.Context()
		userID, ok := ctx.Value("user_id").(string)
		if !ok {
			encodeError(models.ErrInvalidSession, w)
			return
		}
		userUUID, _ := uuid.Parse(userID)

		// get since variable
		since := r.URL.Query().Get("since")
		sinceInt, err := strconv.Atoi(since)
		if err != nil {
			sinceInt = 0
		}

		// get until variable
		until := r.URL.Query().Get("until")
		untilInt, err := strconv.Atoi(until)
		if err != nil {
			untilInt = int(time.Now().Unix())
		}

		// find tracker events
		trackerEvents, err := service.FindTrackerEventsForUserID(userUUID, sinceInt, untilInt)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"tracker_events": trackerEvents,
		})
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		// get session
		ctx := r.Context()
		userID, ok := ctx.Value("user_id").(string)
		if !ok {
			encodeError(models.ErrInvalidSession, w)
			return
		}
		userUUID, _ := uuid.Parse(userID)

		// get timezoneOffset variable
		timezoneOffset := r.URL.Query().Get("timezoneOffset")
		timezoneOffsetInt, err := strconv.Atoi(timezoneOffset)
		if err != nil {
			timezoneOffsetInt = 0
		}

		// bind
		newTrackerEventReq := newTrackerEventRequest{}
		if err := render.Bind(r, &newTrackerEventReq); err != nil {
			encodeError(err, w)
			return
		}

		// insert
		newTrackerEventReq.UserID = userUUID
		trackerEvent, tracker, err := service.InsertTrackerEvent(*newTrackerEventReq.NewTrackerEvent, timezoneOffsetInt)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusCreated, h{
			"tracker_event": trackerEvent,
			"tracker":       tracker,
		})
	})

	router.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
		// get session
		ctx := r.Context()
		userID, ok := ctx.Value("user_id").(string)
		if !ok {
			encodeError(models.ErrInvalidSession, w)
			return
		}
		userUUID, _ := uuid.Parse(userID)

		// parse params
		id := chi.URLParam(r, "id")
		trackerUUID, err := uuid.Parse(id)
		if err != nil {
			encodeError(errInvalidParams, w)
			return
		}

		// get timezoneOffset variable
		timezoneOffset := r.URL.Query().Get("timezoneOffset")
		timezoneOffsetInt, err := strconv.Atoi(timezoneOffset)
		if err != nil {
			timezoneOffsetInt = 0
		}

		// bind
		patchedTrackerEventReq := patchedTrackerEventRequest{}
		if err := render.Bind(r, &patchedTrackerEventReq); err != nil {
			encodeError(err, w)
			return
		}

		// patch
		patchedTrackerEventReq.ID = trackerUUID
		patchedTrackerEventReq.UserID = userUUID
		trackerEvent, tracker, err := service.PatchTrackerEvent(*patchedTrackerEventReq.PatchedTrackerEvent, timezoneOffsetInt)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"tracker_event": trackerEvent,
			"tracker":       tracker,
		})
	})

	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		// get session
		ctx := r.Context()
		userID, ok := ctx.Value("user_id").(string)
		if !ok {
			encodeError(models.ErrInvalidSession, w)
			return
		}
		userUUID, _ := uuid.Parse(userID)

		// get timezoneOffset variable
		timezoneOffset := r.URL.Query().Get("timezoneOffset")
		timezoneOffsetInt, err := strconv.Atoi(timezoneOffset)
		if err != nil {
			timezoneOffsetInt = 0
		}

		// parse params
		id := chi.URLParam(r, "id")
		trackerUUID, err := uuid.Parse(id)
		if err != nil {
			encodeError(errInvalidParams, w)
			return
		}

		// delete
		tracker, err := service.DeleteTrackerEvent(models.DeletedTrackerEvent{
			ID:     trackerUUID,
			UserID: userUUID,
		}, timezoneOffsetInt)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{"tracker": tracker})
	})

	return router
}
