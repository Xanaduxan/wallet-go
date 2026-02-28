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
		"POST /operation",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.CreateOperation)),
	)

	return mux
}
