package domain

import "errors"

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountExists     = errors.New("account already exists")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
