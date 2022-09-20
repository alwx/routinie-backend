package users

import (
	"net/http"
	"time"

	kitLog "github.com/go-kit/kit/log"
	"github.com/google/uuid"

	"habiko-go/pkg/models"
)

type loggingService struct {
	logger kitLog.Logger
	Service
}

func (s *loggingService) GetUser(id uuid.UUID) (*models.User, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "get_user",
			"id", id.String(),
		)
	}(time.Now())
	return s.Service.GetUser(id)
}

func (s *loggingService) GetUserByCredentials(email string, password string) (*models.User, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "get_user_by_credentials",
			"email", email,
		)
	}(time.Now())
	return s.Service.GetUserByCredentials(email, password)
}

func (s *loggingService) InsertUser(newUser models.NewUser) (*models.User, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "insert_user",
		)
	}(time.Now())
	return s.Service.InsertUser(newUser)
}

func (s *loggingService) PatchUser(patchedUser *models.PatchedUser, id uuid.UUID) (*models.User, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "patch_user",
			"id", id.String(),
		)
	}(time.Now())
	return s.Service.PatchUser(patchedUser, id)
}

func (s *loggingService) RemindPassword(email string) (*models.User, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "remind_password",
			"email", email,
		)
	}(time.Now())
	return s.Service.RemindPassword(email)
}

func (s *loggingService) SetPassword(
	email string, remindPasswordToken string, password string,
) (*models.User, error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "set_password",
			"email", email,
		)
	}(time.Now())
	return s.Service.SetPassword(email, remindPasswordToken, password)
}

func (s *loggingService) InitSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "init_session",
		)
	}(time.Now())
	return s.Service.InitSession(w, r, user)
}

func NewLoggingService(logger kitLog.Logger, s Service) Service {
	return &loggingService{logger, s}
}
