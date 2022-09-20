package trackers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"habiko-go/pkg/models"
	"net/http"
)

func MakeHandler(service Service) http.Handler {
	router := chi.NewRouter()

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		// find tracker
		tracker, err := service.FindTrackerByID(userUUID, trackerUUID)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{"tracker": tracker})
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

		// bind
		newTrackerReq := newTrackerRequest{}
		if err := render.Bind(r, &newTrackerReq); err != nil {
			encodeError(err, w)
			return
		}

		// insert
		newTrackerReq.UserID = userUUID
		tracker, err := service.InsertTracker(*newTrackerReq.NewTracker)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusCreated, h{"tracker": tracker})
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

		// bind
		patchedTrackerReq := patchedTrackerRequest{}
		if err := render.Bind(r, &patchedTrackerReq); err != nil {
			encodeError(err, w)
			return
		}

		// patch
		patchedTrackerReq.ID = trackerUUID
		patchedTrackerReq.UserID = userUUID
		tracker, err := service.PatchTracker(*patchedTrackerReq.PatchedTracker)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{"tracker": tracker})
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

		// parse params
		id := chi.URLParam(r, "id")
		trackerUUID, err := uuid.Parse(id)
		if err != nil {
			encodeError(errInvalidParams, w)
			return
		}

		// delete
		err = service.DeleteTracker(models.DeletedTracker{
			ID:     trackerUUID,
			UserID: userUUID,
		})
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{"result": http.StatusText(http.StatusOK)})
	})

	return router
}
