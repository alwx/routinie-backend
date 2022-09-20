package models

import (
	"github.com/stripe/stripe-go/v72"
)

type StripeProvider interface {
	Initialize() error
	EventFromBody(body []byte, signature string) (stripe.Event, error)
	CreateCheckoutSession(user User, priceId string, coupon string) (*stripe.CheckoutSession, error)
	CreateBillingPortalSession(sessionID string) (*stripe.BillingPortalSession, error)
}
