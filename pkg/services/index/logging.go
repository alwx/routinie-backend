package index

import (
	"time"

	kitLog "github.com/go-kit/kit/log"
)

type loggingService struct {
	logger kitLog.Logger
	Service
}

func (s *loggingService) GetIndexPage() []byte {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "get_index_page",
		)
	}(time.Now())
	return s.Service.GetIndexPage()
}

func NewLoggingService(logger kitLog.Logger, s Service) Service {
	return &loggingService{logger, s}
}
