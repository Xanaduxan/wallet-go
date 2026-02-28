package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type Operation struct {
	ID     uuid.UUID
	Name   string
	Type   string
	Amount float64
	UserID uuid.UUID
}

type OperationStorage struct {
	DB *sql.DB
}

func NewOperationStorage(db *sql.DB) *OperationStorage {
	return &OperationStorage{DB: db}
}

func (s *OperationStorage) Create(operation Operation) error {
	_, err := s.DB.Exec(`
		INSERT INTO operations (id, name, type, amount, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`, operation.ID, operation.Name, operation.Type, operation.Amount, operation.UserID)

	return err
}

func (s *OperationStorage) GetById(id uuid.UUID) (Operation, error) {
	var operation Operation

	err := s.DB.QueryRow(`
		SELECT id, name, type, amount, user_id
		FROM operations
		WHERE id = $1
	`, id).Scan(&operation.ID, &operation.Name, &operation.Type, &operation.Amount, &operation.UserID)

	return operation, err
}

func (s *OperationStorage) Update(operation Operation) error {
	_, err := s.DB.Exec(`
		UPDATE operations 
		SET name = $2, type = $3, amount = $4
		WHERE id = $1
	`, operation.ID, operation.Name, operation.Type, operation.Amount)

	return err
}

func (s *OperationStorage) DeleteByID(id uuid.UUID) error {

	_, err := s.DB.Exec(`
		DELETE FROM operations
		WHERE id = $1
	`, id)

	return err
}
