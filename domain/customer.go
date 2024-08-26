package domain

import (
	"context"
	"errors"
)

var ErrKYCFailed = errors.New("KYC validation failed")

type Customer struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Address   string
	KYCStatus string
}

type KYCService interface {
	ValidateKYC(context.Context, *Customer) error
	VerifyCustomerKYC(context.Context, int, *Customer) error
}
