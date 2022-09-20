package tracker_events

import (
	"github.com/go-kit/kit/metrics"
	"habiko-go/pkg/models"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func (i *instrumentingService) InsertTrackerEvent(
	newTrackerEvent models.NewTrackerEvent, timezoneOffset int,
) (*models.TrackerEvent, *models.Tracker, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "insert_tracker_event").Add(1)
		i.requestLatency.With("method", "insert_tracker_event").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.Service.InsertTrackerEvent(newTrackerEvent, timezoneOffset)
}

func (i *instrumentingService) PatchTrackerEvent(
	patchedTrackerEvent models.PatchedTrackerEvent, timezoneOffset int,
) (*models.TrackerEvent, *models.Tracker, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "patch_tracker_event").Add(1)
		i.requestLatency.With("method", "patch_tracker_event").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.Service.PatchTrackerEvent(patchedTrackerEvent, timezoneOffset)
}

func (i *instrumentingService) DeleteTrackerEvent(
	deletedTrackerEvent models.DeletedTrackerEvent, timezoneOffset int,
) (*models.Tracker, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "delete_tracker_event").Add(1)
		i.requestLatency.With("method", "delete_tracker_event").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.Service.DeleteTrackerEvent(deletedTrackerEvent, timezoneOffset)
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{counter, latency, s}
}
