package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Xanaduxan/wallet-go/internal/config"
	"github.com/Xanaduxan/wallet-go/internal/http/handlers"
	"github.com/Xanaduxan/wallet-go/internal/http/router"
	"github.com/Xanaduxan/wallet-go/internal/service/auth"
	"github.com/Xanaduxan/wallet-go/internal/storage"
)

func main() {
	config.LoadEnv(".env")

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	db := storage.NewPostgres(dsn)
	userStorage := storage.NewUserStorage(db)
	authService := auth.NewService(userStorage)
	handlers.SetAuthService(authService)
	r := router.New()
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
