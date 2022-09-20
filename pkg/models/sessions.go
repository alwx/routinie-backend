package models

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionRepository interface {
	InitSessionStore() error

	SaveSession(w http.ResponseWriter, r *http.Request, user *User) error
	GetSession(r *http.Request) (*sessions.Session, error)
}

var ErrInvalidSession = errors.New("invalid session")
