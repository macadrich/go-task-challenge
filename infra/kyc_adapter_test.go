package infra

import (
	"context"
	"sync"
	"testing"

	"github.com/macadrich/go-task-challenge/constants"
	"github.com/macadrich/go-task-challenge/domain"
	"github.com/macadrich/go-task-challenge/external"
	"github.com/stretchr/testify/assert"
)

func TestKYCAdapter(t *testing.T) {
	externalService := new(external.ExternalKYCService)
	adapter := NewKYCAdapter(externalService)

	customer := &domain.Customer{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Address:   "123 Main St",
	}

	ctx := context.Background()
	err := adapter.ValidateKYC(ctx, customer)

	assert.NoError(t, err)
	assert.Equal(t, "pending", customer.KYCStatus)
}

func TestSimulateKYCValidation(t *testing.T) {
	externalService := &external.ExternalKYCService{}
	adapter := NewKYCAdapter(externalService)

	var wg sync.WaitGroup
	numRoutines := constants.NumberOfRoutines

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			customer := &domain.Customer{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Phone:     "1234567890",
				Address:   "123 Main St",
			}

			ctx := context.Background()
			err := adapter.VerifyCustomerKYC(ctx, numRoutines, customer)

			assert.NoError(t, err)
			assert.Equal(t, "approved", customer.KYCStatus)
		}()
	}

	wg.Wait()
}
