package trackers

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

func (i *instrumentingService) InsertTracker(newTracker models.NewTracker) (*models.Tracker, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "insert_tracker").Add(1)
		i.requestLatency.With("method", "insert_tracker").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.Service.InsertTracker(newTracker)
}

func (i *instrumentingService) PatchTracker(patchedTracker models.PatchedTracker) (*models.Tracker, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "patch_tracker").Add(1)
		i.requestLatency.With("method", "patch_tracker").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.Service.PatchTracker(patchedTracker)
}

func (i *instrumentingService) DeleteTracker(deletedTracker models.DeletedTracker) error {
	defer func(begin time.Time) {
		i.requestCount.With("method", "delete_tracker").Add(1)
		i.requestLatency.With("method", "delete_tracker").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.Service.DeleteTracker(deletedTracker)
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{counter, latency, s}
}
