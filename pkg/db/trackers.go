package db

import (
	"github.com/google/uuid"
	"habiko-go/pkg/models"
	"habiko-go/pkg/utils"
)

type TrackerRepository struct {
	*TransactionProvider
}

func (r *TrackerRepository) FindByID(
	transaction models.Transaction, id uuid.UUID,
) (*models.Tracker, error) {
	var tracker models.Tracker
	return &tracker, dbConn(transaction, r.DBConn).First(&tracker, "id = ?", id).Error
}

func (r *TrackerRepository) FindByIDWithStreaks(
	transaction models.Transaction, id uuid.UUID, timezoneOffset int,
) (*models.Tracker, error) {
	var tracker models.Tracker

	operation := dbConn(transaction, r.DBConn).Raw(`
		WITH fulfilling_tracker_events_in_chronological_order AS (
			SELECT e.*
			FROM tracker_events e
			INNER JOIN (
				SELECT tracker_id, assigned_to_date, max(created_at) AS max_created_at
				FROM tracker_events
				GROUP BY tracker_id, assigned_to_date
			) te ON e.created_at = te.max_created_at AND e.tracker_id = te.tracker_id
			INNER JOIN trackers t ON t.id = e.tracker_id
			WHERE e.tracker_id = ? AND (
				e.tracker_fulfilled OR (
					e.type = 'timer_start' AND
					(
						CASE
							WHEN e.assigned_to_date = NOW()::date
							THEN EXTRACT(EPOCH FROM (NOW() - e.created_at)) + e.value
							ELSE EXTRACT(EPOCH FROM ((e.assigned_to_date + interval '1 day' - interval '1 second') - e.created_at)) - ? + e.value
						END
					) >= goal_value
				)
			)
			ORDER BY e.assigned_to_date
		),
		grouped_chronological_tracker_events AS (
			SELECT
				row_number() OVER (ORDER BY assigned_to_date),
				assigned_to_date,
				assigned_to_date - CAST(row_number() OVER (ORDER BY assigned_to_date) as INT) AS group_name,
				tracker_id
			FROM fulfilling_tracker_events_in_chronological_order
		),
		streaks AS (
			SELECT
				max(assigned_to_date) - min(assigned_to_date) + 1 AS streak,
				tracker_id, 
				max(assigned_to_date) as max_date
			FROM grouped_chronological_tracker_events
			GROUP BY tracker_id, group_name
			ORDER BY streak DESC
		),
		tracker_streak_stats AS (
			SELECT
				tracker_id,
				max(streak) max_streak,
				(
					SELECT s.streak 
					FROM streaks s 
					WHERE (s.max_date = CURRENT_DATE OR s.max_date = (CURRENT_DATE - INTERVAL '1 day')::date) AND s.tracker_id = streaks.tracker_id
				) current_streak
			FROM streaks
			GROUP BY tracker_id
		)
		SELECT t.*, s.max_streak, s.current_streak
		FROM trackers t
		LEFT JOIN tracker_streak_stats s ON s.tracker_id = t.id
		WHERE t.id = ?;
	`, id.String(), timezoneOffset, id.String()).Find(&tracker)

	return &tracker, operation.Error
}

func (r *TrackerRepository) FindAllForUserID(userID uuid.UUID, since int) ([]models.Tracker, error) {
	var trackers []models.Tracker

	transaction := r.DBConn.
		Where("trackers.user_id = ?", userID.String()).
		Where("trackers.updated_at >= to_timestamp(?)", since).
		Find(&trackers)

	return trackers, transaction.Error
}

func (r *TrackerRepository) FindAllForUserIDWithStreaks(userID uuid.UUID, since int, timezoneOffset int) ([]models.Tracker, error) {
	var trackers []models.Tracker

	transaction := r.DBConn.Raw(`
		WITH fulfilling_tracker_events_in_chronological_order AS (
			SELECT e.*
			FROM tracker_events e
			INNER JOIN (
				SELECT tracker_id, assigned_to_date, max(created_at) AS max_created_at
				FROM tracker_events
				GROUP BY tracker_id, assigned_to_date
			) te ON e.created_at = te.max_created_at AND e.tracker_id = te.tracker_id
			INNER JOIN trackers t ON t.id = e.tracker_id
			WHERE e.user_id = ? AND (
				e.tracker_fulfilled OR (
					e.type = 'timer_start' AND
					(
						CASE
							WHEN e.assigned_to_date = NOW()::date
							THEN EXTRACT(EPOCH FROM (NOW() - e.created_at)) + e.value
							ELSE EXTRACT(EPOCH FROM ((e.assigned_to_date + interval '1 day' - interval '1 second') - e.created_at)) - ? + e.value
						END
					) >= goal_value
				)
			)
			ORDER BY e.assigned_to_date
		),
		grouped_chronological_tracker_events AS (
			SELECT
				row_number() OVER (ORDER BY tracker_id, assigned_to_date),
				assigned_to_date,
				assigned_to_date - CAST(row_number() OVER (ORDER BY tracker_id, assigned_to_date) as INT) AS group_name,
				tracker_id
			FROM fulfilling_tracker_events_in_chronological_order
		),
		streaks AS (
			SELECT
				max(assigned_to_date) - min(assigned_to_date) + 1 AS streak,
				tracker_id, 
				max(assigned_to_date) as max_date
			FROM grouped_chronological_tracker_events
			GROUP BY tracker_id, group_name
			ORDER BY streak DESC
		),
		tracker_streak_stats AS (
			SELECT
				tracker_id,
				max(streak) max_streak,
				(
					SELECT s.streak 
					FROM streaks s 
					WHERE (s.max_date = CURRENT_DATE OR s.max_date = (CURRENT_DATE - INTERVAL '1 day')::date) AND s.tracker_id = streaks.tracker_id
				) current_streak
			FROM streaks
			GROUP BY tracker_id
		)
		SELECT t.*, s.max_streak, s.current_streak 
		FROM trackers t 
		LEFT JOIN tracker_streak_stats s ON s.tracker_id = t.id
		WHERE t.updated_at >= to_timestamp(?) AND t.user_id = ?;
	`, userID.String(), timezoneOffset, since, userID.String()).Scan(&trackers)

	return trackers, transaction.Error
}

func (r *TrackerRepository) FindAllPublicForUserIDWithStreaks(userID uuid.UUID, since int, timezoneOffset int) ([]models.Tracker, error) {
	var trackers []models.Tracker

	transaction := r.DBConn.Raw(`
		WITH fulfilling_tracker_events_in_chronological_order AS (
			SELECT e.*
			FROM tracker_events e
			INNER JOIN (
				SELECT tracker_id, assigned_to_date, max(created_at) AS max_created_at
				FROM tracker_events
				GROUP BY tracker_id, assigned_to_date
			) te ON e.created_at = te.max_created_at AND e.tracker_id = te.tracker_id
			INNER JOIN trackers t ON t.id = e.tracker_id
			WHERE e.user_id = ? AND (
				e.tracker_fulfilled OR (
					e.type = 'timer_start' AND
					(
						CASE
							WHEN e.assigned_to_date = NOW()::date
							THEN EXTRACT(EPOCH FROM (NOW() - e.created_at)) + e.value
							ELSE EXTRACT(EPOCH FROM ((e.assigned_to_date + interval '1 day' - interval '1 second') - e.created_at)) - ? + e.value
						END
					) >= goal_value
				)
			)
			ORDER BY e.assigned_to_date
		),
		grouped_chronological_tracker_events AS (
			SELECT
				row_number() OVER (ORDER BY tracker_id, assigned_to_date),
				assigned_to_date,
				assigned_to_date - CAST(row_number() OVER (ORDER BY tracker_id, assigned_to_date) as INT) AS group_name,
				tracker_id
			FROM fulfilling_tracker_events_in_chronological_order
		),
		streaks AS (
			SELECT
				max(assigned_to_date) - min(assigned_to_date) + 1 AS streak,
				tracker_id, 
				max(assigned_to_date) as max_date
			FROM grouped_chronological_tracker_events
			GROUP BY tracker_id, group_name
			ORDER BY streak DESC
		),
		tracker_streak_stats AS (
			SELECT
				tracker_id,
				max(streak) max_streak,
				(
					SELECT s.streak 
					FROM streaks s 
					WHERE (s.max_date = CURRENT_DATE OR s.max_date = (CURRENT_DATE - INTERVAL '1 day')::date) AND s.tracker_id = streaks.tracker_id
				) current_streak
			FROM streaks
			GROUP BY tracker_id
		)
		SELECT t.id, t.title, t.type, t.color, t.default_value, t.goal_value, t.default_change, t.is_infinite, t.rank, s.max_streak, s.current_streak 
		FROM trackers t 
		LEFT JOIN tracker_streak_stats s ON s.tracker_id = t.id
		WHERE t.updated_at >= to_timestamp(?) AND t.is_public = 'true' AND t.user_id = ?
		ORDER BY t.rank;
	`, userID.String(), timezoneOffset, since, userID.String()).Scan(&trackers)

	return trackers, transaction.Error
}

func (r *TrackerRepository) Insert(
	t models.Transaction, newTracker models.NewTracker,
) (*models.Tracker, error) {
	trackerType := utils.OneOf(newTracker.Type, models.TrackerTypes, models.DefaultTrackerType)
	tracker := models.Tracker{
		ID: uuid.New(),

		UserID:          &newTracker.UserID,
		ParentTrackerID: newTracker.ParentTrackerID,

		Title: &newTracker.Title,
		Type:  &trackerType,
		Color: &newTracker.Color,

		DefaultValue:  &newTracker.DefaultValue,
		GoalValue:     &newTracker.GoalValue,
		Measurement:   &newTracker.Measurement,
		DefaultChange: &newTracker.DefaultChange,
		IsInfinite:    &newTracker.IsInfinite,
		IsPublic:      &newTracker.IsPublic,
		Rank:          &newTracker.Rank,
	}
	return &tracker, dbConn(t, r.DBConn).Create(&tracker).Error
}

func (r *TrackerRepository) Patch(
	t models.Transaction, patchedTracker models.PatchedTracker,
) error {
	tracker := models.Tracker{
		ID: patchedTracker.ID,
	}

	if patchedTracker.ParentTrackerID != nil {
		tracker.ParentTrackerID = patchedTracker.ParentTrackerID
	}

	if patchedTracker.Title != nil {
		tracker.Title = patchedTracker.Title
	}
	if patchedTracker.Type != nil {
		trackerType := utils.OneOf(*patchedTracker.Type, models.TrackerTypes, models.DefaultTrackerType)
		tracker.Type = &trackerType
	}
	if patchedTracker.Color != nil {
		tracker.Color = patchedTracker.Color
	}

	if patchedTracker.DefaultValue != nil {
		tracker.DefaultValue = patchedTracker.DefaultValue
	}
	if patchedTracker.GoalValue != nil {
		tracker.GoalValue = patchedTracker.GoalValue
	}
	if patchedTracker.Measurement != nil {
		tracker.Measurement = patchedTracker.Measurement
	}
	if patchedTracker.DefaultChange != nil {
		tracker.DefaultChange = patchedTracker.DefaultChange
	}
	if patchedTracker.IsInfinite != nil {
		tracker.IsInfinite = patchedTracker.IsInfinite
	}
	if patchedTracker.IsPublic != nil {
		tracker.IsPublic = patchedTracker.IsPublic
	}
	if patchedTracker.Rank != nil {
		tracker.Rank = patchedTracker.Rank
	}

	return dbConn(t, r.DBConn).Model(&tracker).Updates(tracker).Error
}

func (r *TrackerRepository) Delete(t models.Transaction, deletedTracker models.DeletedTracker) error {
	return dbConn(t, r.DBConn).Where("id = ?", deletedTracker.ID).Delete(models.Tracker{}).Error
}
