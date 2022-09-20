package index

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func (s *instrumentingService) GetIndexPage() []byte {
	defer func(begin time.Time) {
		s.requestCount.With(
			"method", "get_index_page",
			"params", "",
		).Add(1)
		s.requestLatency.With(
			"method", "get_index_page",
			"params", "",
		).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.GetIndexPage()
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{counter, latency, s}
}
