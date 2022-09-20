package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"habiko-go/pkg/models"
)

var errCannotSaveSession = errors.New("cannot save the session")
var errSessionIsAlreadyInitialized = errors.New("session is already initialized")

var errUserNotFound = errors.New("user not found")
var errMissingData = errors.New("no user data provided")
var errIncorrectEmail = errors.New("incorrect e-mail")
var errIncorrectLogin = errors.New("incorrect login: it must contain at least 3 characters")
var errIncorrectPassword = errors.New("incorrect password: it must contain at least 6 characters")
var errIncorrectOldPassword = errors.New("incorrect old password")
var errIncorrectConfirmationToken = errors.New("incorrect code")
var errIncorrectRemindPasswordToken = errors.New("incorrect code")
var errWrongCredentials = errors.New("user not found or the provided credentials are wrong")
var errEmailIsInUse = errors.New("the specified e-mail is already in use")

type patchedUserRequest struct {
	*models.PatchedUser
}

func (request *patchedUserRequest) Bind(r *http.Request) error {
	if request.PatchedUser == nil {
		return errMissingData
	}
	return nil
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (request *loginUserRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Email) == "" {
		return errIncorrectEmail
	}
	if strings.TrimSpace(request.Password) == "" {
		return errIncorrectPassword
	}
	return nil
}

type remindPasswordRequest struct {
	Email string `json:"email"`
}

func (request *remindPasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Email) == "" {
		return errIncorrectEmail
	}
	return nil
}

type createUserRequest struct {
	*models.NewUser
}

func (request *createUserRequest) Bind(r *http.Request) error {
	login := strings.TrimSpace(request.Login)
	if login == "" || len(login) < 3 {
		return errIncorrectLogin
	}
	if strings.TrimSpace(request.Email) == "" {
		return errIncorrectEmail
	}
	password := strings.TrimSpace(request.Password)
	if password == "" || len(password) < 6 {
		return errIncorrectPassword
	}
	return nil
}

type setPasswordRequest struct {
	Email               string `json:"email"`
	RemindPasswordToken string `json:"remind_password_token"`
	Password            string `json:"password"`
}

func (request *setPasswordRequest) Bind(r *http.Request) error {
	email := strings.TrimSpace(request.Email)
	remindPasswordToken := strings.TrimSpace(request.RemindPasswordToken)
	password := strings.TrimSpace(request.Password)

	if email == "" {
		return errIncorrectEmail
	}
	if remindPasswordToken == "" || len(remindPasswordToken) < 6 {
		return errIncorrectRemindPasswordToken
	}
	if password == "" || len(password) < 6 {
		return errIncorrectPassword
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
	case errCannotSaveSession:
		w.WriteHeader(http.StatusInternalServerError)
	case errSessionIsAlreadyInitialized:
		w.WriteHeader(http.StatusForbidden)
	case errUserNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errMissingData:
		w.WriteHeader(http.StatusBadRequest)
	case errIncorrectEmail:
		w.WriteHeader(http.StatusBadRequest)
	case errIncorrectLogin:
		w.WriteHeader(http.StatusBadRequest)
	case errIncorrectPassword:
		w.WriteHeader(http.StatusBadRequest)
	case errIncorrectOldPassword:
		w.WriteHeader(http.StatusBadRequest)
	case errIncorrectConfirmationToken:
		w.WriteHeader(http.StatusBadRequest)
	case errIncorrectRemindPasswordToken:
		w.WriteHeader(http.StatusBadRequest)
	case errWrongCredentials:
		w.WriteHeader(http.StatusForbidden)
	case errEmailIsInUse:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": err.Error(),
	})
}
