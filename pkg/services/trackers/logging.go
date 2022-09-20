package trackers

import (
	kitLog "github.com/go-kit/kit/log"
	"habiko-go/pkg/models"
	"time"
)

type loggingService struct {
	logger kitLog.Logger
	Service
}

func (l *loggingService) InsertTracker(newTracker models.NewTracker) (*models.Tracker, error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "insert_tracker",
			"user_id", newTracker.UserID.String(),
		)
	}(time.Now())
	return l.Service.InsertTracker(newTracker)
}

func (l *loggingService) PatchTracker(patchedTracker models.PatchedTracker) (*models.Tracker, error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "patch_tracker",
			"id", patchedTracker.ID.String(),
			"user_id", patchedTracker.UserID.String(),
		)
	}(time.Now())
	return l.Service.PatchTracker(patchedTracker)
}

func (l *loggingService) DeleteTracker(deletedTracker models.DeletedTracker) error {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "delete_tracker",
			"id", deletedTracker.ID.String(),
			"user_id", deletedTracker.UserID.String(),
		)
	}(time.Now())
	return l.Service.DeleteTracker(deletedTracker)
}

func NewLoggingService(logger kitLog.Logger, s Service) Service {
	return &loggingService{logger, s}
}
