package sample_trackers

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

func (i *instrumentingService) FindSampleTrackersBySimilarTag(tag string) ([]models.SampleTracker, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "find_sample_trackers_by_similar_tag").Add(1)
		i.requestLatency.With(
			"method", "find_sample_trackers_by_similar_tag",
		).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return i.Service.FindSampleTrackersBySimilarTag(tag)
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{counter, latency, s}
}
