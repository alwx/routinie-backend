package stripe

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"habiko-go/pkg/models"
	"habiko-go/pkg/services"
	"net/http"
)

type Service interface {
	CreateCheckoutSession(userID uuid.UUID, productId string, coupon string) (*stripe.CheckoutSession, error)
	CreateBillingPortalSession(sessionID string) (*stripe.BillingPortalSession, error)

	EventFromBody(body []byte, r *http.Request) (stripe.Event, error)
	ExtractCheckoutSession(event stripe.Event) (*stripe.CheckoutSession, error)
	ExtractSubscription(event stripe.Event) (*stripe.Subscription, error)
	GetUser(userId string) (*models.User, error)

	HandleSubscriptionCreated(user models.User, subscription stripe.Subscription) error
	HandleSubscriptionCanceled(user models.User, subscription stripe.Subscription) error
	HandleCheckoutSessionCompleted(user models.User, checkoutSession stripe.CheckoutSession) error
}

type service struct {
	providers *services.Providers
}

func (s *service) CreateCheckoutSession(userID uuid.UUID, priceId string, coupon string) (*stripe.CheckoutSession, error) {
	user, err := s.providers.UserRepository.FindByID(nil, userID)
	if err != nil {
		return nil, errUserNotFound
	}

	return s.providers.StripeProvider.CreateCheckoutSession(*user, priceId, coupon)
}

func (s *service) CreateBillingPortalSession(sessionID string) (*stripe.BillingPortalSession, error) {
	return s.providers.StripeProvider.CreateBillingPortalSession(sessionID)
}

func (s *service) EventFromBody(body []byte, r *http.Request) (stripe.Event, error) {
	return s.providers.StripeProvider.EventFromBody(body, r.Header.Get("Stripe-Signature"))
}

func (s *service) ExtractCheckoutSession(event stripe.Event) (*stripe.CheckoutSession, error) {
	var checkoutSession stripe.CheckoutSession
	err := json.Unmarshal(event.Data.Raw, &checkoutSession)
	if err != nil {
		return nil, err
	}
	return &checkoutSession, nil
}

func (s *service) ExtractSubscription(event stripe.Event) (*stripe.Subscription, error) {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (s *service) GetUser(userId string) (*models.User, error) {
	if userId == "" {
		return nil, errUserNotFound
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	user, err := s.providers.UserRepository.FindByID(nil, userUUID)
	if err != nil {
		return nil, errUserNotFound
	}

	return user, nil
}

func (s *service) HandleSubscriptionCreated(user models.User, subscription stripe.Subscription) error {
	patchedUser := models.PatchedUser{
		ID:           user.ID,
		IsSubscribed: stripe.Bool(true),
	}
	_, err := s.providers.UserRepository.PatchStripeData(nil, patchedUser)
	if err != nil {
		return models.ErrDb
	}
	return nil
}

func (s *service) HandleSubscriptionCanceled(user models.User, subscription stripe.Subscription) error {
	patchedUser := models.PatchedUser{
		ID:           user.ID,
		IsSubscribed: stripe.Bool(false),
	}
	_, err := s.providers.UserRepository.PatchStripeData(nil, patchedUser)
	if err != nil {
		return models.ErrDb
	}
	return nil
}

func (s *service) HandleCheckoutSessionCompleted(user models.User, checkoutSession stripe.CheckoutSession) error {
	patchedUser := models.PatchedUser{
		ID:              user.ID,
		StripeSessionId: &checkoutSession.ID,
	}
	_, err := s.providers.UserRepository.PatchStripeData(nil, patchedUser)
	if err != nil {
		return models.ErrDb
	}
	return nil
}

func NewService(providers *services.Providers) Service {
	return &service{providers: providers}
}
