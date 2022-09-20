package redis

import (
	"encoding/hex"
	"net/http"

	"github.com/boj/redistore"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"

	"habiko-go/pkg/models"
	"habiko-go/pkg/utils"
)

type SessionRepository struct {
	Redis *redis.Client
	store *redistore.RediStore
}

func (service *SessionRepository) InitSessionStore() error {
	store, err := redistore.NewRediStore(
		10,
		"tcp",
		service.Redis.Options().Addr,
		"",
		[]byte(viper.GetString("web.session_key")),
	)
	if err != nil {
		return err
	}

	cookieDomain := viper.GetString("web.cookie_domain")
	if cookieDomain != "" {
		store.Options.Domain = cookieDomain
	}
	service.store = store

	return nil
}

func (service *SessionRepository) SaveSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	userSession, err := service.store.Get(r, "routinie")
	if err != nil {
		return err
	}

	sessionID := hex.EncodeToString(utils.NewEntropy(256))
	userSession.Values["session_id"] = sessionID
	userSession.Values["user_id"] = user.ID.String()
	_ = service.store.Save(r, w, userSession)
	return nil
}

func (service *SessionRepository) GetSession(r *http.Request) (*sessions.Session, error) {
	return service.store.Get(r, "routinie")
}
