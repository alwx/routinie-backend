package index

import (
	"net/http"

	"github.com/go-chi/chi"
)

func MakeHandler(service Service) http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		bytes := service.GetIndexPage()
		encodeResponse(w, http.StatusOK, bytes)
	})

	_ = service.ServeStatic(router, "/*")

	return router
}
