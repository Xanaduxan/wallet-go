package handlers

import "github.com/Xanaduxan/wallet-go/internal/service/operations"

func validateOperation(name, opType string, amount float64) error {
	if name == "" || amount <= 0 {
		return operations.ErrInvalidInput
	}
	if opType != "income" && opType != "expense" {
		return operations.ErrInvalidInput
	}
	return nil
}
