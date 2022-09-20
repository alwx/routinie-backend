package users

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
		until := time.Now().Unix()

		// get timezoneOffset variable
		timezoneOffset := r.URL.Query().Get("timezoneOffset")
		timezoneOffsetInt, err := strconv.Atoi(timezoneOffset)
		if err != nil {
			timezoneOffsetInt = 0
		}

		// get user
		user, err := service.GetUser(userUUID)
		if err != nil {
			encodeError(errUserNotFound, w)
			return
		}

		trackers, err := service.FindTrackersWithStreaksForUserID(userUUID, sinceInt, timezoneOffsetInt)
		if err != nil {
			encodeError(err, w)
			return
		}

		// return everything
		encodeResponse(w, http.StatusOK, h{
			"user":     user,
			"trackers": trackers,
			"until":    until,
		})
	})

	router.Get("/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		// get timezoneOffset variable
		timezoneOffset := r.URL.Query().Get("timezoneOffset")
		timezoneOffsetInt, err := strconv.Atoi(timezoneOffset)
		if err != nil {
			timezoneOffsetInt = 0
		}

		user, err := service.GetPublicUserDataByName(name)
		if err != nil || user.Login == "" {
			encodeError(errUserNotFound, w)
			return
		}

		t := time.Now()
		until := int(t.Unix())
		since := int(t.AddDate(0, 0, -12).Unix())

		trackers, err := service.FindPublicTrackersWithStreaksForUserID(user.ID, since, timezoneOffsetInt)
		if err != nil {
			encodeError(err, w)
			return
		}

		// find tracker events
		trackerEvents, err := service.FindPublicTrackerEventsForUserID(user.ID, since, until)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"user":           user,
			"trackers":       trackers,
			"tracker_events": trackerEvents,
		})
	})

	router.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		// make sure session is not created yet
		ctx := r.Context()
		_, ok := ctx.Value("session_id").(string)
		if ok {
			encodeError(errSessionIsAlreadyInitialized, w)
			return
		}

		// bind
		createUserReq := createUserRequest{}
		if err := render.Bind(r, &createUserReq); err != nil {
			encodeError(err, w)
			return
		}

		// create a new user
		user, err := service.InsertUser(*createUserReq.NewUser)
		if err != nil {
			encodeError(err, w)
			return
		}
		until := time.Now().Unix()

		// initialize session
		err = service.InitSession(w, r, user)
		if err != nil {
			encodeError(errCannotSaveSession, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"user":  user,
			"until": until,
		})
	})

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		// bind
		loginUserReq := loginUserRequest{}
		if err := render.Bind(r, &loginUserReq); err != nil {
			encodeError(err, w)
			return
		}

		// try to find a user
		user, err := service.GetUserByCredentials(loginUserReq.Email, loginUserReq.Password)
		if err != nil {
			encodeError(errWrongCredentials, w)
			return
		}

		// initialize session
		err = service.InitSession(w, r, user)
		if err != nil {
			encodeError(errCannotSaveSession, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"user": user,
		})
	})

	router.Post("/remind-password", func(w http.ResponseWriter, r *http.Request) {
		// get session
		ctx := r.Context()
		_, ok := ctx.Value("user_id").(string)
		if ok {
			encodeError(errSessionIsAlreadyInitialized, w)
			return
		}

		// bind
		req := remindPasswordRequest{}
		if err := render.Bind(r, &req); err != nil {
			encodeError(err, w)
			return
		}

		// remind password
		if _, err := service.RemindPassword(req.Email); err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"result": http.StatusText(http.StatusOK),
		})
	})

	router.Post("/set-password", func(w http.ResponseWriter, r *http.Request) {
		// get session
		ctx := r.Context()
		_, ok := ctx.Value("user_id").(string)
		if ok {
			encodeError(errSessionIsAlreadyInitialized, w)
			return
		}

		// bind
		req := setPasswordRequest{}
		if err := render.Bind(r, &req); err != nil {
			encodeError(err, w)
			return
		}

		// set password
		if _, err := service.SetPassword(req.Email, req.RemindPasswordToken, req.Password); err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"result": http.StatusText(http.StatusOK),
		})
	})

	router.Patch("/", func(w http.ResponseWriter, r *http.Request) {
		// get session
		ctx := r.Context()
		userID, ok := ctx.Value("user_id").(string)
		if !ok {
			encodeError(models.ErrInvalidSession, w)
			return
		}
		userUUID, _ := uuid.Parse(userID)

		// bind
		req := patchedUserRequest{}
		if err := render.Bind(r, &req); err != nil {
			encodeError(err, w)
			return
		}

		// try to patch
		patchedUser, err := service.PatchUser(req.PatchedUser, userUUID)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"user": patchedUser,
		})
	})

	return router
}
