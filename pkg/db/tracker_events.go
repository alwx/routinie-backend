package db

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"habiko-go/pkg/models"
	"habiko-go/pkg/utils"
	"time"
)

type TrackerEventRepository struct {
	*TransactionProvider
}

func (r *TrackerEventRepository) FindByID(
	transaction models.Transaction, id uuid.UUID,
) (*models.TrackerEvent, error) {
	var trackerEvent models.TrackerEvent
	return &trackerEvent, dbConn(transaction, r.DBConn).First(&trackerEvent, "id = ?", id).Error
}

func (r *TrackerEventRepository) FindAllForUserID(
	userID uuid.UUID, since int, until int, isPublic bool,
) ([]models.TrackerEvent, error) {
	var trackerEvents []models.TrackerEvent

	var selectQuery string
	if isPublic {
		selectQuery = "e.id, e.tracker_id, e.value, e.type, e.assigned_to_date, e.tracker_fulfilled, e.created_at"
	} else {
		selectQuery = "e.*"
	}

	transaction := r.DBConn.Raw(`
		SELECT `+selectQuery+`
		FROM tracker_events e
		INNER JOIN (
			SELECT tracker_id, assigned_to_date, max(created_at) AS max_created_at
			FROM tracker_events
			WHERE 
				user_id = ? 
				AND assigned_to_date >= to_timestamp(?)::date 
				AND assigned_to_date <= to_timestamp(?)::date
			GROUP BY tracker_id, assigned_to_date
		) te ON e.created_at = te.max_created_at AND e.tracker_id = te.tracker_id
		ORDER BY e.assigned_to_date
	`, userID.String(), since, until).Find(&trackerEvents)

	return trackerEvents, transaction.Error
}

func (r *TrackerEventRepository) getLatestTrackerEventForTracker(
	transaction models.Transaction, trackerID uuid.UUID, assignedToDate datatypes.Date,
) (*models.TrackerEvent, error) {
	var trackerEvent models.TrackerEvent

	db := dbConn(transaction, r.DBConn).
		Where("tracker_events.tracker_id = ? AND tracker_events.assigned_to_date = ?", trackerID.String(), assignedToDate).
		Order("tracker_events.created_at DESC").
		First(&trackerEvent)

	return &trackerEvent, db.Error
}

func (r *TrackerEventRepository) isTrackerFulfilled(tracker models.Tracker, value int) bool {
	isGoingUp := (*tracker.GoalValue - *tracker.DefaultValue) >= 0
	return (isGoingUp && value >= *tracker.GoalValue) || (!isGoingUp && value <= *tracker.GoalValue)
}

func (r *TrackerEventRepository) isAllowedToAddAtDate(addAt time.Time, user models.User) bool {
	now := time.Now()

	var numberOfDays int
	if user.SubscribedAt != nil {
		numberOfDays = 4
	} else {
		numberOfDays = 2
	}

	earliestAllowedDate := time.Date(
		now.Year(),
		now.Month(),
		now.Day()-numberOfDays,
		0,
		0,
		0,
		0,
		addAt.Location(),
	)
	return earliestAllowedDate.Before(addAt)
}

func (r *TrackerEventRepository) Insert(
	t models.Transaction,
	user models.User,
	existingTracker models.Tracker,
	newTrackerEvent models.NewTrackerEvent,
) (*models.TrackerEvent, error) {
	if !r.isAllowedToAddAtDate(time.Time(newTrackerEvent.AssignedToDate), user) {
		return nil, models.ErrTooLateToAddTrackerEvent
	}
	isFulfilled := r.isTrackerFulfilled(existingTracker, newTrackerEvent.Value)
	trackerEventType := utils.OneOf(newTrackerEvent.Type, models.TrackerEventTypes, models.DefaultTrackerEventType)
	trackerEvent := models.TrackerEvent{
		ID:               uuid.New(),
		UserID:           &newTrackerEvent.UserID,
		TrackerID:        existingTracker.ID,
		Type:             &trackerEventType,
		Value:            &newTrackerEvent.Value,
		AssignedToDate:   &newTrackerEvent.AssignedToDate,
		TrackerFulfilled: &isFulfilled,
	}
	return &trackerEvent, dbConn(t, r.DBConn).Create(&trackerEvent).Error
}

func (r *TrackerEventRepository) Patch(
	t models.Transaction,
	user models.User,
	existingTracker models.Tracker,
	patchedTrackerEvent models.PatchedTrackerEvent,
) error {
	trackerEvent := models.TrackerEvent{
		ID: patchedTrackerEvent.ID,
	}

	// TODO(alwx): some changes to the logic might be required
	if patchedTrackerEvent.AssignedToDate != nil {
		if !r.isAllowedToAddAtDate(time.Time(*patchedTrackerEvent.AssignedToDate), user) {
			return models.ErrTooLateToAddTrackerEvent
		}

		isFulfilled := r.isTrackerFulfilled(existingTracker, *patchedTrackerEvent.Value)

		trackerEvent.Value = patchedTrackerEvent.Value
		trackerEvent.AssignedToDate = patchedTrackerEvent.AssignedToDate
		trackerEvent.TrackerFulfilled = &isFulfilled
	}

	return dbConn(t, r.DBConn).Model(&trackerEvent).Updates(trackerEvent).Error
}

func (r *TrackerEventRepository) Delete(
	t models.Transaction, deletedTrackerEvent models.DeletedTrackerEvent,
) error {
	return dbConn(t, r.DBConn).Where("id = ?", deletedTrackerEvent.ID).Delete(models.TrackerEvent{}).Error
}
