package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORS)
	r.Use(RateLimit)

	r.Route("/api", func(r chi.Router) {
		r.Get("/pack-sizes", handler.GetPackSizes)
		r.Post("/pack-sizes", handler.UpdatePackSizes)
		r.Post("/calculate", handler.CalculatePacks)
	})

	return r
}
