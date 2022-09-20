package sample_trackers

import (
	"github.com/go-chi/chi"
	"net/http"
)

func MakeHandler(service Service) http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")

		sampleTrackers, err := service.FindSampleTrackersBySimilarTag(tag)
		if err != nil {
			encodeError(err, w)
			return
		}

		encodeResponse(w, http.StatusOK, h{
			"sample_trackers": sampleTrackers,
		})
	})

	return router
}
