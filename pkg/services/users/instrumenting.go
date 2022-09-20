package users

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/google/uuid"

	"habiko-go/pkg/models"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func (s *instrumentingService) GetUser(id uuid.UUID) (*models.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get_user").Add(1)
		s.requestLatency.With("method", "get_user").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.GetUser(id)
}

func (s *instrumentingService) GetUserByCredentials(email string, password string) (*models.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get_user_by_credentials").Add(1)
		s.requestLatency.With("method", "get_user_by_credentials").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.GetUserByCredentials(email, password)
}

func (s *instrumentingService) InsertUser(newUser models.NewUser) (*models.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "insert_user").Add(1)
		s.requestLatency.With("method", "insert_user").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.InsertUser(newUser)
}

func (s *instrumentingService) PatchUser(patchedUser *models.PatchedUser, id uuid.UUID) (*models.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "patch_user").Add(1)
		s.requestLatency.With("method", "patch_user").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.PatchUser(patchedUser, id)
}

func (s *instrumentingService) RemindPassword(email string) (*models.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "remind_password").Add(1)
		s.requestLatency.With("method", "remind_password").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.RemindPassword(email)
}

func (s *instrumentingService) SetPassword(
	email string, remindPasswordToken string, password string,
) (*models.User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "set_password").Add(1)
		s.requestLatency.With("method", "set_password").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.SetPassword(email, remindPasswordToken, password)
}

func (s *instrumentingService) InitSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "init_session").Add(1)
		s.requestLatency.With("method", "init_session").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.InitSession(w, r, user)
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{counter, latency, s}
}
