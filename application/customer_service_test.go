package application

import (
	"context"
	"testing"

	"github.com/macadrich/go-task-challenge/constants"
	"github.com/macadrich/go-task-challenge/domain"
	"github.com/macadrich/go-task-challenge/infra"
	"github.com/macadrich/go-task-challenge/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterCustomer(t *testing.T) {
	mockKYC := new(mocks.MockKYCService)
	mockKYC.On("ValidateKYC", mock.Anything, mock.Anything).Return(nil)

	customerRepository := infra.NewCustomerRepository()
	customerService := NewCustomerService(mockKYC, customerRepository)

	customer := &domain.Customer{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Address:   "123 Main St",
	}

	ctx := context.Background()
	err := customerService.RegisterCustomer(ctx, customer)

	assert.NoError(t, err)
	assert.Equal(t, "pending", customer.KYCStatus)
	mockKYC.AssertCalled(t, "ValidateKYC", mock.Anything, customer)
}

func TestVerifyCustomer(t *testing.T) {
	mockKYC := new(mocks.MockKYCService)
	mockKYC.On("VerifyCustomerKYC", mock.Anything, mock.Anything).Return(nil)

	customerRepository := infra.NewCustomerRepository()
	customerService := NewCustomerService(mockKYC, customerRepository)

	customer := &domain.Customer{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Address:   "123 Main St",
	}

	ctx := context.Background()
	err := customerService.VerifyRegisteredCustomer(ctx, constants.NumberOfRoutines, customer)

	assert.NoError(t, err)
	assert.Equal(t, "approved", customer.KYCStatus)
	mockKYC.AssertCalled(t, "VerifyCustomerKYC", mock.Anything, customer)
}
