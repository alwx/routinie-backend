package tracker_events

import (
	kitLog "github.com/go-kit/kit/log"
	"habiko-go/pkg/models"
	"time"
)

type loggingService struct {
	logger kitLog.Logger
	Service
}

func (l *loggingService) InsertTrackerEvent(
	newTrackerEvent models.NewTrackerEvent, timezoneOffset int,
) (*models.TrackerEvent, *models.Tracker, error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "insert_tracker_event",
			"user_id", newTrackerEvent.UserID.String(),
		)
	}(time.Now())
	return l.Service.InsertTrackerEvent(newTrackerEvent, timezoneOffset)
}

func (l *loggingService) PatchTrackerEvent(
	patchedTrackerEvent models.PatchedTrackerEvent, timezoneOffset int,
) (*models.TrackerEvent, *models.Tracker, error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "patch_tracker_event",
			"id", patchedTrackerEvent.ID.String(),
			"user_id", patchedTrackerEvent.UserID.String(),
		)
	}(time.Now())
	return l.Service.PatchTrackerEvent(patchedTrackerEvent, timezoneOffset)
}

func (l *loggingService) DeleteTrackerEvent(
	deletedTrackerEvent models.DeletedTrackerEvent, timezoneOffset int,
) (*models.Tracker, error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "delete_tracker_event",
			"id", deletedTrackerEvent.ID.String(),
			"user_id", deletedTrackerEvent.UserID.String(),
		)
	}(time.Now())
	return l.Service.DeleteTrackerEvent(deletedTrackerEvent, timezoneOffset)
}

func NewLoggingService(logger kitLog.Logger, s Service) Service {
	return &loggingService{logger, s}
}
