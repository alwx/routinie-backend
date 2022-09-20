package users

import (
	"net/http"

	"github.com/google/uuid"

	"habiko-go/pkg/models"
	"habiko-go/pkg/services"
)

type Service interface {
	GetUser(id uuid.UUID) (*models.User, error)
	GetUserByCredentials(email string, password string) (*models.User, error)
	GetPublicUserDataByName(name string) (*models.PublicUserData, error)

	FindTrackersWithStreaksForUserID(userID uuid.UUID, since int, timezoneOffset int) ([]models.Tracker, error)
	FindPublicTrackersWithStreaksForUserID(userID uuid.UUID, since int, timezoneOffset int) ([]models.Tracker, error)
	FindPublicTrackerEventsForUserID(userID uuid.UUID, since int, until int) ([]models.TrackerEvent, error)

	InsertUser(newUser models.NewUser) (*models.User, error)
	PatchUser(patchedUser *models.PatchedUser, id uuid.UUID) (*models.User, error)
	RemindPassword(email string) (*models.User, error)
	SetPassword(email string, remindPasswordToken string, password string) (*models.User, error)

	InitSession(w http.ResponseWriter, r *http.Request, user *models.User) error
}

type service struct {
	providers *services.Providers
}

func (s *service) GetUser(id uuid.UUID) (*models.User, error) {
	user, err := s.providers.UserRepository.FindByID(nil, id)
	if err != nil {
		return nil, err
	}

	user.IsEmailConfirmed = s.providers.UserRepository.IsEmailConfirmed(*user)

	return user, nil
}

func (s *service) GetUserByCredentials(email string, password string) (*models.User, error) {
	return s.providers.UserRepository.FindByCredentials(nil, email, password)
}

func (s *service) GetPublicUserDataByName(name string) (*models.PublicUserData, error) {
	return s.providers.UserRepository.FindPublicDataByName(nil, name)
}

func (s *service) FindTrackersWithStreaksForUserID(userID uuid.UUID, since int, timezoneOffset int) ([]models.Tracker, error) {
	trackers, err := s.providers.TrackerRepository.FindAllForUserIDWithStreaks(userID, since, timezoneOffset)
	if err != nil {
		return nil, models.ErrDb
	}
	return trackers, nil
}

func (s *service) FindPublicTrackersWithStreaksForUserID(userID uuid.UUID, since int, timezoneOffset int) ([]models.Tracker, error) {
	trackers, err := s.providers.TrackerRepository.FindAllPublicForUserIDWithStreaks(userID, since, timezoneOffset)
	if err != nil {
		return nil, models.ErrDb
	}
	return trackers, nil
}

func (s *service) FindPublicTrackerEventsForUserID(userID uuid.UUID, since int, until int) ([]models.TrackerEvent, error) {
	trackerEvents, err := s.providers.TrackerEventRepository.FindAllForUserID(userID, since, until, true)
	if err != nil {
		return nil, models.ErrDb
	}
	return trackerEvents, nil
}

func (s *service) InsertUser(newUser models.NewUser) (*models.User, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, models.ErrDb
	}

	// insert
	user, err := s.providers.UserRepository.Insert(transaction, newUser)
	if err != nil {
		_ = transaction.Rollback()
		switch err {
		case models.ErrDbEmail:
			return nil, errIncorrectEmail
		case models.ErrDbPassword:
			return nil, errIncorrectPassword
		default:
			return nil, errEmailIsInUse
		}
	}

	// send email
	if user.EmailConfirmationToken != nil && *user.EmailConfirmationToken != "" {
		go s.providers.EmailHandler.SendUserEmailConfirmation(*user)
	}
	user.IsEmailConfirmed = false

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		_ = transaction.Rollback()
		return nil, models.ErrDb
	}

	return user, nil
}

func (s *service) PatchUser(
	patchedUser *models.PatchedUser, id uuid.UUID,
) (*models.User, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, models.ErrDb
	}

	// find previous user
	previousUser, err := s.providers.UserRepository.FindByID(transaction, id)
	if err != nil {
		_ = transaction.Rollback()
		return nil, errUserNotFound
	}

	// patch
	patchedUser.ID = id
	newUser, err := s.providers.UserRepository.Patch(transaction, *patchedUser, *previousUser)
	if err != nil {
		_ = transaction.Rollback()
		switch err {
		case models.ErrDbEmail:
			return nil, errIncorrectEmail
		case models.ErrDbLogin:
			return nil, errIncorrectLogin
		case models.ErrDbPassword:
			return nil, errIncorrectPassword
		case models.ErrIncorrectOldPassword:
			return nil, errIncorrectOldPassword
		case models.ErrIncorrectConfirmationToken:
			return nil, errIncorrectConfirmationToken
		case models.ErrIncorrectRemindPasswordToken:
			return nil, errIncorrectRemindPasswordToken
		default:
			return nil, errEmailIsInUse
		}
	}

	// send email
	if newUser.EmailConfirmationToken != nil && *newUser.EmailConfirmationToken != "" {
		go s.providers.EmailHandler.SendUserEmailConfirmation(*newUser)
	}

	// find user
	user, err := s.providers.UserRepository.FindByID(transaction, id)
	if err != nil {
		_ = transaction.Rollback()
		return nil, models.ErrDb
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, models.ErrDb
	}

	user.IsEmailConfirmed = s.providers.UserRepository.IsEmailConfirmed(*user)
	return user, nil
}

func (s *service) RemindPassword(email string) (*models.User, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, models.ErrDb
	}

	// find previous user
	previousUser, err := s.providers.UserRepository.FindByEmail(transaction, email)
	if err != nil {
		_ = transaction.Rollback()
		return nil, errUserNotFound
	}

	// patch
	patchedUser := models.PatchedUser{
		ID:                                previousUser.ID,
		ShouldGenerateRemindPasswordToken: true,
	}
	newUser, err := s.providers.UserRepository.Patch(transaction, patchedUser, *previousUser)
	if err != nil {
		_ = transaction.Rollback()
		return nil, models.ErrDbObject
	}

	// send email
	if newUser.RemindPasswordToken != nil && *newUser.RemindPasswordToken != "" {
		newUser.Email = email
		go s.providers.EmailHandler.SendUserRemindPassword(*newUser)
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, models.ErrDb
	}

	return newUser, nil
}

func (s *service) SetPassword(email string, remindPasswordToken string, password string) (*models.User, error) {
	// start transaction
	transaction, err := s.providers.TransactionProvider.Create()
	if err != nil {
		return nil, models.ErrDb
	}

	// find previous user
	previousUser, err := s.providers.UserRepository.FindByEmail(transaction, email)
	if err != nil {
		_ = transaction.Rollback()
		return nil, errUserNotFound
	}

	// patch
	patchedUser := models.PatchedUser{
		ID:                  previousUser.ID,
		RemindPasswordToken: &remindPasswordToken,
		Password:            &password,
	}
	newUser, err := s.providers.UserRepository.Patch(transaction, patchedUser, *previousUser)
	if err != nil {
		_ = transaction.Rollback()
		switch err {
		case models.ErrIncorrectRemindPasswordToken:
			return nil, errIncorrectRemindPasswordToken
		default:
			return nil, models.ErrDbObject
		}
	}

	// commit transaction
	err = transaction.Commit()
	if err != nil {
		return nil, models.ErrDb
	}

	return newUser, nil
}

func (s *service) InitSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	return s.providers.SessionRepository.SaveSession(w, r, user)
}

func NewService(providers *services.Providers) Service {
	return &service{providers: providers}
}
