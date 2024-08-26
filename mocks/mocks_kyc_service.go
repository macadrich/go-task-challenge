package mocks

import (
	"context"

	"github.com/macadrich/go-task-challenge/domain"
	"github.com/stretchr/testify/mock"
)

type MockKYCService struct {
	mock.Mock
}

func (m *MockKYCService) ValidateKYC(ctx context.Context, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	customer.KYCStatus = "pending"
	return args.Error(0)
}

func (m *MockKYCService) VerifyCustomerKYC(ctx context.Context, numRequest int, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	customer.KYCStatus = "approved"
	return args.Error(0)
}
