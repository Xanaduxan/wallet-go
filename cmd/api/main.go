package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Xanaduxan/wallet-go/internal/config"
	"github.com/Xanaduxan/wallet-go/internal/http/handlers"
	"github.com/Xanaduxan/wallet-go/internal/http/router"
	"github.com/Xanaduxan/wallet-go/internal/service/auth"
	"github.com/Xanaduxan/wallet-go/internal/service/operations"
	"github.com/Xanaduxan/wallet-go/internal/storage"
)

func main() {
	config.LoadEnv(".env")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	db := storage.NewPostgres(dsn)
	userStorage := storage.NewUserStorage(db)

	authService := auth.NewService(userStorage, []byte(jwtSecret))
	handlers.SetAuthService(authService)
	operationStorage := storage.NewOperationStorage(db)
	operationsService := operations.NewService(operationStorage, userStorage)
	handlers.SetOperationService(operationsService)
	r := router.New([]byte(jwtSecret))
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
