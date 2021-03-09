package middleware

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func LoggerMiddleware(lf *middleware.DefaultLogFormatter, r *chi.Mux) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(lf)
}
