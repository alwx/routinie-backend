package services

import (
	"habiko-go/pkg/models"
)

type Providers struct {
	TransactionProvider models.TransactionProvider
	StripeProvider      models.StripeProvider
	EmailHandler        models.EmailHandler

	SessionRepository       models.SessionRepository
	TrackerRepository       models.TrackerRepository
	TrackerEventRepository  models.TrackerEventRepository
	SampleTrackerRepository models.SampleTrackerRepository
	UserRepository          models.UserRepository
}
