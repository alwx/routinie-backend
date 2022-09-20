package trackers

import (
	"github.com/google/uuid"
	"habiko-go/pkg/models"
	"habiko-go/pkg/services"
)

type Service interface {
	FindTrackerByID(userID uuid.UUID, id uuid.UUID) (*models.Tracker, error)

	InsertTracker(newTracker models.NewTracker) (*models.Tracker, error)
	PatchTracker(patchedTracker models.PatchedTracker) (*models.Tracker, error)
	DeleteTracker(deletedTracker models.DeletedTracker) error
}

type service struct {
	providers *services.Providers
}

func (s *service) FindTrackerByID(userID uuid.UUID, id uuid.UUID) (*models.Tracker, error) {
	tracker, err := s.providers.TrackerRepository.FindByID(nil, id)
	if err != nil {
		return nil, errNotFound
	}
	if *tracker.UserID != userID {
		return nil, errAccessDenied
	}
	return tracker, nil
}

func (s *service) InsertTracker(newTracker models.NewTracker) (*models.Tracker, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, models.ErrDb
	}

	// check permissions
	user, err := s.providers.UserRepository.FindByID(transaction, newTracker.UserID)
	if err != nil {
		_ = transaction.Rollback()
		return nil, errUserAccessDenied
	}
	existingTrackers, err := s.providers.TrackerRepository.FindAllForUserID(user.ID, 0)
	if err != nil {
		_ = transaction.Rollback()
		return nil, errUserAccessDenied
	}
	if user.SubscribedAt == nil && len(existingTrackers) >= 3 {
		return nil, errNoPremiumSubscription
	}

	// insert tracker
	tracker, err := s.providers.TrackerRepository.Insert(transaction, newTracker)
	if err != nil {
		_ = transaction.Rollback()
		return nil, models.ErrDbObject
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, models.ErrDb
	}

	return tracker, nil
}

func (s *service) PatchTracker(patchedTracker models.PatchedTracker) (*models.Tracker, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, models.ErrDb
	}

	// get tracker
	tracker, err := s.providers.TrackerRepository.FindByID(transaction, patchedTracker.ID)
	if err != nil {
		return nil, errNotFound
	}
	if *tracker.UserID != patchedTracker.UserID {
		return nil, errAccessDenied
	}

	// patch
	err = s.providers.TrackerRepository.Patch(transaction, patchedTracker)
	if err != nil {
		_ = transaction.Rollback()
		return nil, models.ErrDbObject
	}

	// find user
	tracker, err = s.providers.TrackerRepository.FindByID(transaction, patchedTracker.ID)
	if err != nil {
		_ = transaction.Rollback()
		return nil, models.ErrDb
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, models.ErrDb
	}

	return tracker, nil
}

func (s *service) DeleteTracker(deletedTracker models.DeletedTracker) error {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return models.ErrDb
	}

	// get tracker
	tracker, err := s.providers.TrackerRepository.FindByID(transaction, deletedTracker.ID)
	if err != nil {
		return errNotFound
	}
	if *tracker.UserID != deletedTracker.UserID {
		return errAccessDenied
	}

	// delete
	err = s.providers.TrackerRepository.Delete(transaction, deletedTracker)
	if err != nil {
		return models.ErrDb
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return models.ErrDb
	}

	return nil
}

func NewService(providers *services.Providers) Service {
	return &service{providers: providers}
}
