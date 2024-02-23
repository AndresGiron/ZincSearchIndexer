package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Permitir solicitudes preflight OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Routes() *chi.Mux {
	mux := chi.NewMux()

	mux.Use(
		middleware.Logger,
		middleware.Recoverer,
		corsHandler,
	)
	mux.Post("/search", Search)
	mux.Post("/all", ListAll)
	mux.Get("/hello", HelloHandler)

	return mux
}
