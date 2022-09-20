package stripe

import (
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

func MakeHandler(service Service) http.Handler {
	router := chi.NewRouter()

	router.Post("/create-checkout-session", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			encodeError(errInvalidFormData, w)
			return
		}
		userId := r.PostFormValue("userId")
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			encodeError(errInvalidFormData, w)
			return
		}

		priceId := r.PostFormValue("priceId")
		if priceId == "" {
			encodeError(errInvalidFormData, w)
			return
		}

		// create stripe checkout session
		session, err := service.CreateCheckoutSession(userUUID, priceId, r.PostFormValue("coupon"))
		if err != nil {
			encodeError(err, w)
			return
		}
		http.Redirect(w, r, session.URL, http.StatusSeeOther)
	})

	router.Post("/create-portal-session", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			encodeError(errInvalidFormData, w)
			return
		}
		sessionId := r.PostFormValue("sessionId")

		// get session
		session, err := service.CreateBillingPortalSession(sessionId)
		if err != nil {
			encodeError(errInvalidSession, w)
			return
		}
		http.Redirect(w, r, session.URL, http.StatusSeeOther)
	})

	router.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			encodeError(err, w)
			return
		}

		event, err := service.EventFromBody(body, r)
		if err != nil {
			encodeError(errIncorrectStripeSignature, w)
			return
		}

		switch event.Type {
		case "checkout.session.completed":
			checkoutSession, err := service.ExtractCheckoutSession(event)
			if err != nil {
				encodeError(err, w)
				return
			}
			user, err := service.GetUser(checkoutSession.ClientReferenceID)
			if err != nil {
				encodeError(err, w)
				return
			}
			err = service.HandleCheckoutSessionCompleted(*user, *checkoutSession)
			if err != nil {
				encodeError(err, w)
				return
			}

		case "customer.subscription.trial_will_end":
		case "customer.subscription.created":
			subscription, err := service.ExtractSubscription(event)
			if err != nil {
				encodeError(err, w)
				return
			}
			user, err := service.GetUser(subscription.Metadata["user_id"])
			if err != nil {
				encodeError(err, w)
				return
			}
			err = service.HandleSubscriptionCreated(*user, *subscription)
			if err != nil {
				encodeError(err, w)
				return
			}
		case "customer.subscription.deleted":
			subscription, err := service.ExtractSubscription(event)
			if err != nil {
				encodeError(err, w)
				return
			}
			user, err := service.GetUser(subscription.Metadata["user_id"])
			if err != nil {
				encodeError(err, w)
				return
			}
			err = service.HandleSubscriptionCanceled(*user, *subscription)
			if err != nil {
				encodeError(err, w)
				return
			}
		case "invoice.paid":
		case "invoice.payment_action_required":
		case "invoice.payment_failed":
		default:
		}

		encodeResponse(w, http.StatusOK, nil)
	})

	return router
}
