package initializers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	kitLog "github.com/go-kit/kit/log"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"habiko-go/pkg/services"
	"habiko-go/pkg/services/index"
	"habiko-go/pkg/services/sample_trackers"
	"habiko-go/pkg/services/stripe"
	"habiko-go/pkg/services/tracker_events"
	"habiko-go/pkg/services/trackers"
	"habiko-go/pkg/services/users"
)

func prepareCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{viper.GetString("web.allowed_origin")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

func apiRouter(providers *services.Providers, logger kitLog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(prepareCORS().Handler)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userSession, err := providers.SessionRepository.GetSession(r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			userID := userSession.Values["user_id"]
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	usersService := users.NewService(providers)
	usersService = users.NewLoggingService(logger, usersService)
	usersService = users.NewInstrumentingService(
		kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "users_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "users_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
		usersService,
	)
	router.Mount("/users", users.MakeHandler(usersService))

	trackersService := trackers.NewService(providers)
	trackersService = trackers.NewLoggingService(logger, trackersService)
	trackersService = trackers.NewInstrumentingService(
		kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "trackers_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "trackers_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
		trackersService,
	)
	router.Mount("/trackers", trackers.MakeHandler(trackersService))

	trackerEventsService := tracker_events.NewService(providers)
	trackerEventsService = tracker_events.NewLoggingService(logger, trackerEventsService)
	trackerEventsService = tracker_events.NewInstrumentingService(
		kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "tracker_events_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "tracker_events_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
		trackerEventsService,
	)
	router.Mount("/tracker_events", tracker_events.MakeHandler(trackerEventsService))

	sampleTrackersService := sample_trackers.NewService(providers)
	sampleTrackersService = sample_trackers.NewLoggingService(logger, sampleTrackersService)
	sampleTrackersService = sample_trackers.NewInstrumentingService(
		kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "sample_trackers_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "sample_trackers_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
		sampleTrackersService,
	)
	router.Mount("/sample_trackers", sample_trackers.MakeHandler(sampleTrackersService))

	return router
}

func indexRouter(providers *services.Providers, logger kitLog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeHTML))

	indexService := index.NewService(providers)
	indexService = index.NewLoggingService(logger, indexService)
	indexService = index.NewInstrumentingService(
		kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "index",
			Subsystem: "index_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method", "params"}),
		kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: "index",
			Subsystem: "index_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method", "params"}),
		indexService,
	)
	router.Mount("/", index.MakeHandler(indexService))

	return router
}

func stripeRouter(providers *services.Providers, logger kitLog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(prepareCORS().Handler)

	stripeService := stripe.NewService(providers)
	router.Mount("/", stripe.MakeHandler(stripeService))

	return router
}

func CreateBackendRouter(providers *services.Providers, logger kitLog.Logger) *chi.Mux {
	err := providers.SessionRepository.InitSessionStore()
	if err != nil {
		panic(err)
	}
	err = providers.StripeProvider.Initialize()
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.DefaultCompress)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Mount("/", indexRouter(providers, logger))
	router.Mount("/api", apiRouter(providers, logger))
	router.Mount("/stripe", stripeRouter(providers, logger))
	router.Mount("/metrics", promhttp.Handler())

	return router
}
