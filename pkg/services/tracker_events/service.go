package tracker_events

import (
	"github.com/google/uuid"
	"habiko-go/pkg/models"
	"habiko-go/pkg/services"
)

type Service interface {
	FindTrackerEventsForUserID(userID uuid.UUID, since int, until int) ([]models.TrackerEvent, error)

	InsertTrackerEvent(newTrackerEvent models.NewTrackerEvent, timezoneOffset int) (*models.TrackerEvent, *models.Tracker, error)
	PatchTrackerEvent(patchedTrackerEvent models.PatchedTrackerEvent, timezoneOffset int) (*models.TrackerEvent, *models.Tracker, error)
	DeleteTrackerEvent(deletedTrackerEvent models.DeletedTrackerEvent, timezoneOffset int) (*models.Tracker, error)
}

type service struct {
	providers *services.Providers
}

func (s *service) FindTrackerEventsForUserID(userID uuid.UUID, since int, until int) ([]models.TrackerEvent, error) {
	trackerEvents, err := s.providers.TrackerEventRepository.FindAllForUserID(userID, since, until, false)
	if err != nil {
		return nil, models.ErrDb
	}
	return trackerEvents, nil
}

func (s *service) InsertTrackerEvent(newTrackerEvent models.NewTrackerEvent, timezoneOffset int) (*models.TrackerEvent, *models.Tracker, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, nil, models.ErrDb
	}

	// check permissions
	tracker, err := s.providers.TrackerRepository.FindByID(transaction, newTrackerEvent.TrackerID)
	if err != nil || *tracker.UserID != newTrackerEvent.UserID {
		_ = transaction.Rollback()
		return nil, nil, errTrackerAccessDenied
	}
	user, err := s.providers.UserRepository.FindByID(transaction, newTrackerEvent.UserID)
	if err != nil {
		_ = transaction.Rollback()
		return nil, nil, errUserAccessDenied
	}

	// insert tracker event
	trackerEvent, err := s.providers.TrackerEventRepository.Insert(transaction, *user, *tracker, newTrackerEvent)
	if err != nil {
		_ = transaction.Rollback()
		switch err {
		case models.ErrTooLateToAddTrackerEvent:
			return nil, nil, errTooLateToAddTrackerEvent
		default:
			return nil, nil, models.ErrDbObject
		}
	}

	// get tracker with streaks
	tracker, err = s.providers.TrackerRepository.FindByIDWithStreaks(transaction, newTrackerEvent.TrackerID, timezoneOffset)
	if err != nil {
		_ = transaction.Rollback()
		return nil, nil, err
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, nil, models.ErrDb
	}

	return trackerEvent, tracker, nil
}

func (s *service) PatchTrackerEvent(patchedTrackerEvent models.PatchedTrackerEvent, timezoneOffset int) (*models.TrackerEvent, *models.Tracker, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, nil, models.ErrDb
	}

	// get tracker event
	trackerEvent, err := s.providers.TrackerEventRepository.FindByID(transaction, patchedTrackerEvent.ID)
	if err != nil {
		return nil, nil, errNotFound
	}

	// check permissions — users can patch only their own tracker events
	if *trackerEvent.UserID != patchedTrackerEvent.UserID {
		return nil, nil, errAccessDenied
	}
	existingTracker, err := s.providers.TrackerRepository.FindByID(transaction, trackerEvent.TrackerID)
	if err != nil || *existingTracker.UserID != *trackerEvent.UserID {
		_ = transaction.Rollback()
		return nil, nil, errTrackerAccessDenied
	}
	user, err := s.providers.UserRepository.FindByID(transaction, *trackerEvent.UserID)
	if err != nil {
		_ = transaction.Rollback()
		return nil, nil, errUserAccessDenied
	}

	// patch
	err = s.providers.TrackerEventRepository.Patch(transaction, *user, *existingTracker, patchedTrackerEvent)
	if err != nil {
		_ = transaction.Rollback()
		switch err {
		case models.ErrTooLateToAddTrackerEvent:
			return nil, nil, errTooLateToAddTrackerEvent
		default:
			return nil, nil, models.ErrDbObject
		}
	}

	// find the patched tracker
	trackerEvent, err = s.providers.TrackerEventRepository.FindByID(transaction, patchedTrackerEvent.ID)
	if err != nil {
		_ = transaction.Rollback()
		return nil, nil, models.ErrDb
	}

	// get tracker with streaks
	tracker, err := s.providers.TrackerRepository.FindByIDWithStreaks(transaction, trackerEvent.TrackerID, timezoneOffset)
	if err != nil {
		_ = transaction.Rollback()
		return nil, nil, err
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, nil, models.ErrDb
	}

	return trackerEvent, tracker, nil
}

func (s *service) DeleteTrackerEvent(deletedTrackerEvent models.DeletedTrackerEvent, timezoneOffset int) (*models.Tracker, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, models.ErrDb
	}

	// get tracker event
	trackerEvent, err := s.providers.TrackerEventRepository.FindByID(transaction, deletedTrackerEvent.ID)
	if err != nil {
		return nil, errNotFound
	}

	// check permissions — users can delete only their own tracker events
	if *trackerEvent.UserID != deletedTrackerEvent.UserID {
		return nil, errAccessDenied
	}

	// delete
	err = s.providers.TrackerEventRepository.Delete(transaction, deletedTrackerEvent)
	if err != nil {
		return nil, models.ErrDb
	}

	// get tracker with streaks
	tracker, err := s.providers.TrackerRepository.FindByIDWithStreaks(transaction, trackerEvent.TrackerID, timezoneOffset)
	if err != nil {
		_ = transaction.Rollback()
		return nil, err
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, models.ErrDb
	}

	return tracker, nil
}

func NewService(providers *services.Providers) Service {
	return &service{providers: providers}
}
