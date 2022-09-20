package sample_trackers

import (
	"habiko-go/pkg/models"
	"habiko-go/pkg/services"
)

type Service interface {
	FindSampleTrackersBySimilarTag(tag string) ([]models.SampleTracker, error)
}

type service struct {
	providers *services.Providers
}

func (s *service) FindSampleTrackersBySimilarTag(tag string) ([]models.SampleTracker, error) {
	sampleTrackers, err := s.providers.SampleTrackerRepository.FindAllBySimilarTag(tag, 5)
	if err != nil {
		return nil, models.ErrDb
	}
	return sampleTrackers, nil
}

func NewService(providers *services.Providers) Service {
	return &service{providers: providers}
}
