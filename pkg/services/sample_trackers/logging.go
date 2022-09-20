package sample_trackers

import (
	kitLog "github.com/go-kit/kit/log"
	"habiko-go/pkg/models"
	"time"
)

type loggingService struct {
	logger kitLog.Logger
	Service
}

func (l *loggingService) FindSampleTrackersBySimilarTag(tag string) ([]models.SampleTracker, error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "find_sample_trackers_by_similar_tag",
			"tag", tag,
		)
	}(time.Now())
	return l.Service.FindSampleTrackersBySimilarTag(tag)
}

func NewLoggingService(logger kitLog.Logger, s Service) Service {
	return &loggingService{logger, s}
}
