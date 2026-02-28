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

func (s *Service) GetOperation(email string, operationID uuid.UUID) (storage.Operation, error) {
	if email == "" || operationID == uuid.Nil {
		return storage.Operation{}, ErrInvalidInput
	}

	user, err := s.users.GetByEmail(email)
	if err != nil {
		return storage.Operation{}, err
	}

	op, err := s.operations.GetById(operationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Operation{}, ErrNotFound
		}
		return storage.Operation{}, err
	}

	if op.UserID != user.ID {
		return storage.Operation{}, ErrForbidden
	}

	return op, nil
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
	err = s.operations.Create(op)
	if err != nil {
		return uuid.Nil, err
	}

	return op.ID, nil
}

func (s *Service) DeleteOperation(email string, operationID uuid.UUID) error {
	if email == "" || operationID == uuid.Nil {
		return ErrInvalidInput
	}

	user, err := s.users.GetByEmail(email)
	if err != nil {

		return ErrInvalidInput
	}

	op, err := s.operations.GetById(operationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if op.UserID != user.ID {
		return ErrForbidden
	}

	return s.operations.DeleteByID(operationID)
}

func (s *Service) UpdateOperation(email string, operationID uuid.UUID, name, opType string, amount float64) error {
	if email == "" || operationID == uuid.Nil || name == "" || amount <= 0 {
		return ErrInvalidInput
	}
	if opType != "income" && opType != "expense" {
		return ErrInvalidInput
	}

	user, err := s.users.GetByEmail(email)
	if err != nil {
		return err
	}

	op, err := s.operations.GetById(operationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if op.UserID != user.ID {
		return ErrForbidden
	}

	op.Name = name
	op.Type = opType
	op.Amount = amount

	return s.operations.Update(op)
}
