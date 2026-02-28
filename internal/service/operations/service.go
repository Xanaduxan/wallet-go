package operations

import (
	"database/sql"
	"errors"

	"github.com/Xanaduxan/wallet-go/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	operations *storage.OperationStorage
	users      *storage.UserStorage
}

func NewService(operations *storage.OperationStorage, users *storage.UserStorage) *Service {
	return &Service{operations: operations, users: users}
}

func (s *Service) CreateOperation(email, name, opType string, amount float64) (uuid.UUID, error) {
	if email == "" || name == "" || amount <= 0 {
		return uuid.Nil, ErrInvalidInput
	}
	if opType != "income" && opType != "expense" {
		return uuid.Nil, ErrInvalidInput
	}

	user, err := s.users.GetByEmail(email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrInvalidInput
		}
		return uuid.Nil, err
	}

	op := storage.Operation{
		ID:     uuid.New(),
		Name:   name,
		Type:   opType,
		Amount: amount,
		UserID: user.ID,
	}

	if err := s.operations.Create(op); err != nil {
		return uuid.Nil, err
	}

	return op.ID, nil
}
