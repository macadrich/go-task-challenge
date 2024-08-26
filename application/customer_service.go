package application

import (
	"context"
	"errors"

	"github.com/macadrich/go-task-challenge/domain"
)

var (
	ErrCustomerExists = errors.New("customer already exists")
)

type CustomerRepository interface {
	Save(context.Context, *domain.Customer) error
	FindByEmail(context.Context, string) (*domain.Customer, error)
}

type CustomerService struct {
	kycService         domain.KYCService
	customerRepository CustomerRepository
}

func NewCustomerService(kycService domain.KYCService, customerRepository CustomerRepository) *CustomerService {
	return &CustomerService{
		kycService:         kycService,
		customerRepository: customerRepository,
	}
}

func (s *CustomerService) RegisterCustomer(ctx context.Context, customer *domain.Customer) error {
	existingCustomer, _ := s.customerRepository.FindByEmail(ctx, customer.Email)
	if existingCustomer != nil {
		return ErrCustomerExists
	}

	customer.KYCStatus = "pending"

	if err := s.kycService.ValidateKYC(ctx, customer); err != nil {
		return err
	}

	if err := s.customerRepository.Save(ctx, customer); err != nil {
		return err
	}

	return nil
}

func (s *CustomerService) VerifyRegisteredCustomer(ctx context.Context, numRequest int, customer *domain.Customer) error {
	// Validate multiple customer KYC
	if err := s.kycService.VerifyCustomerKYC(ctx, numRequest, customer); err != nil {
		return err
	}

	return nil
}
