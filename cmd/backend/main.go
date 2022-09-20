//go:generate go install github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=../../pkg/templates

package main

import (
	"fmt"
	"habiko-go/pkg/utils/email"
	"habiko-go/pkg/utils/stripe"
	"net/http"
	"os"

	kitLog "github.com/go-kit/kit/log"
	"github.com/spf13/viper"

	"habiko-go/pkg/config"
	"habiko-go/pkg/db"
	"habiko-go/pkg/initializers"
	"habiko-go/pkg/services"
	"habiko-go/pkg/utils/redis"
)

func main() {
	config.ReadConfig()

	// basic
	dbConn := db.Connect()
	redisConn := redis.Connect()

	// email sender
	emailSender := &email.SendgridEmailSender{APIKey: viper.GetString("email.sendgrid_api_key")}

	// stripe
	stripeProvider := stripe.StripeProvider{
		APIKey:         viper.GetString("stripe.key"),
		CallbackDomain: viper.GetString("stripe.callback_domain"),
		WebhookSecret:  viper.GetString("stripe.webhook_secret"),
	}

	// providers
	transactionProvider := &db.TransactionProvider{
		DBConn:          dbConn.Connection,
		TimescaleDBConn: dbConn.TimescaleConnection,
	}
	providers := services.Providers{
		TransactionProvider: transactionProvider,
		StripeProvider:      &stripeProvider,
		EmailHandler:        &email.EmailHandler{EmailSender: emailSender},

		SessionRepository:       &redis.SessionRepository{Redis: redisConn},
		TrackerRepository:       &db.TrackerRepository{TransactionProvider: transactionProvider},
		TrackerEventRepository:  &db.TrackerEventRepository{TransactionProvider: transactionProvider},
		SampleTrackerRepository: &db.SampleTrackerRepository{TransactionProvider: transactionProvider},
		UserRepository:          &db.UserRepository{TransactionProvider: transactionProvider},
	}

	// port
	port := fmt.Sprintf(":%s", viper.GetString("web.port"))

	fmt.Printf("Listening on port %s...\n", port)

	// logger
	logger := kitLog.NewLogfmtLogger(os.Stderr)
	logger = kitLog.With(logger, "listen", port, "caller", kitLog.DefaultCaller)

	// combine and run
	_ = http.ListenAndServe(port, initializers.CreateBackendRouter(&providers, logger))
}
