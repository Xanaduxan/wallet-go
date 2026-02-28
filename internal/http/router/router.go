package router

import (
	"net/http"

	"github.com/Xanaduxan/wallet-go/internal/http/handlers"
	"github.com/Xanaduxan/wallet-go/internal/http/handlers/middleware"
)

func New(jwtSecret []byte) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"GET /me",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.Me)),
	)

	mux.HandleFunc("POST /login", handlers.Login)
	mux.HandleFunc("POST /registration", handlers.Registration)

	mux.Handle(
		"DELETE /me",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.DeleteMe)),
	)

	mux.Handle(
		"GET /operation/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.GetOperation)),
	)
	mux.Handle(
		"POST /operation",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.CreateOperation)),
	)
	mux.Handle(
		"PUT /operation/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.UpdateOperation)),
	)
	mux.Handle(
		"DELETE /operation/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.DeleteOperation)),
	)

	return mux
}
