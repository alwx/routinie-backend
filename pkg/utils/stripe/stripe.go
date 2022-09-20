package stripe

import (
	"github.com/stripe/stripe-go/v72"
	portalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/webhook"
	"habiko-go/pkg/models"
)

type StripeProvider struct {
	APIKey         string
	CallbackDomain string
	WebhookSecret  string
}

func (provider *StripeProvider) Initialize() error {
	stripe.Key = provider.APIKey
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "routinie/checkout",
		Version: "1.0.0",
		URL:     "https://github.com/alwx/routinie-backend",
	})
	return nil
}

func (provider *StripeProvider) EventFromBody(body []byte, signature string) (stripe.Event, error) {
	return webhook.ConstructEvent(body, signature, provider.WebhookSecret)
}

func (provider *StripeProvider) CreateBillingPortalSession(sessionID string) (*stripe.BillingPortalSession, error) {
	s, err := session.Get(sessionID, nil)
	if err != nil {
		return nil, err
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(s.Customer.ID),
		ReturnURL: stripe.String(provider.CallbackDomain + "/settings"),
	}
	return portalsession.New(params)
}

func (provider *StripeProvider) CreateCheckoutSession(user models.User, priceId string, coupon string) (*stripe.CheckoutSession, error) {
	var allowPromotionCodes *bool
	var discounts []*stripe.CheckoutSessionDiscountParams
	if coupon == "" {
		allowPromotionCodes = stripe.Bool(true)
		discounts = nil
	} else {
		allowPromotionCodes = nil
		discounts = []*stripe.CheckoutSessionDiscountParams{
			{Coupon: stripe.String(coupon)},
		}
	}
	params := &stripe.CheckoutSessionParams{
		SuccessURL:          stripe.String(provider.CallbackDomain + "/premium/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:           stripe.String(provider.CallbackDomain + "/premium/cancelled"),
		ClientReferenceID:   stripe.String(user.ID.String()),
		CustomerEmail:       &user.Email,
		Mode:                stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		AllowPromotionCodes: allowPromotionCodes,
		Discounts:           discounts,
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			TrialPeriodDays: stripe.Int64(7),
			Metadata: map[string]string{
				"user_id": user.ID.String(),
			},
		},
		// AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
	}
	return session.New(params)
}
