package infra

import (
	"context"
	"errors"
	"sync"

	"github.com/macadrich/go-task-challenge/domain"
	"github.com/macadrich/go-task-challenge/external"
)

type KYCAdapter struct {
	externalService *external.ExternalKYCService
}

func NewKYCAdapter(externalService *external.ExternalKYCService) *KYCAdapter {
	return &KYCAdapter{externalService: externalService}
}

func (a *KYCAdapter) ValidateKYC(ctx context.Context, customer *domain.Customer) error {
	// Map the domain customer to the external service request format.
	request := &external.ExternalKYCRequest{
		FullName: customer.FirstName + " " + customer.LastName,
		Email:    customer.Email,
		Phone:    customer.Phone,
		Address:  customer.Address,
	}

	// Call the external service.
	response, err := a.externalService.Validate(request)
	if err != nil {
		return err
	}

	// Update the customer KYC status based on the external service response.
	if response.Status == "pending" {
		customer.KYCStatus = "pending"
	} else {
		return domain.ErrKYCFailed
	}

	return nil
}

func (a *KYCAdapter) VerifyCustomerKYC(ctx context.Context, numRequest int, customer *domain.Customer) error {
	// Map the domain customer to the external service request format.
	request := &external.ExternalKYCRequest{
		FullName: customer.FirstName + " " + customer.LastName,
		Email:    customer.Email,
		Phone:    customer.Phone,
		Address:  customer.Address,
	}

	results := make(chan *external.ExternalKYCResponse, numRequest)
	errorsChan := make(chan error, numRequest)

	var wg sync.WaitGroup

	for i := 0; i < numRequest; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Simulate receiving verification result from external API
			response, err := a.externalService.Verify(request)
			if err != nil {
				errorsChan <- err
				return
			}
			results <- response
		}()
	}

	go func() {
		wg.Wait()
		close(results)
		close(errorsChan)
	}()

	var finalStatus string
	for response := range results {
		if response.Status == "approved" {
			finalStatus = "approved"
		} else {
			finalStatus = "rejected"
		}
	}

	for err := range errorsChan {
		if err != nil {
			return err
		}
	}

	if finalStatus == "approved" {
		customer.KYCStatus = "approved"
		return nil
	}

	return errors.New("KYC validation failed")
}
