package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Xanaduxan/wallet-go/internal/storage"
)

type Service struct {
	users *storage.UserStorage
}

var JwtSecret = []byte("secret")

var (
	ErrValidationFailed   = errors.New("validation failed")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("email and password required")
)

func NewService(users *storage.UserStorage) *Service {
	return &Service{users: users}
}

func (s *Service) Registration(email, password string) (string, error) {

	if email == "" || password == "" {
		return "", ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := storage.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hash),
	}

	err = s.users.Create(user)
	if err != nil {
		return "", err
	}

	return s.generateToken(email)
}

func (s *Service) Login(email, password string) (string, error) {

	if email == "" || password == "" {
		return "", ErrInvalidInput
	}

	user, err := s.users.GetByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateToken(email)
}

func (s *Service) DeleteByEmail(email string) error {

	user, err := s.users.GetByEmail(email)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(user.ID)
	if err != nil {
		return err
	}

	err = s.users.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) generateToken(email string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
