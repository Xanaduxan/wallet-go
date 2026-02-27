package router

import (
	"net/http"

	"github.com/Xanaduxan/wallet-go/internal/http/handlers"
	"github.com/Xanaduxan/wallet-go/internal/http/handlers/middleware"
)

func New() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(
		"GET /me",
		middleware.JWT(http.HandlerFunc(handlers.Me)),
	)
	mux.HandleFunc("POST /login", handlers.Login)
	mux.HandleFunc("POST /registration", handlers.Registration)
	mux.Handle(
		"DELETE /me",
		middleware.JWT(http.HandlerFunc(handlers.DeleteMe)),
	)
	return mux
}
