package infra

import (
	"context"
	"errors"
	"sync"

	"github.com/macadrich/go-task-challenge/domain"
)

// CustomerRepository to simulate database, in-memory customer repository
type CustomerRepository struct {
	mu        *sync.Mutex
	customers map[string]*domain.Customer
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{
		mu:        &sync.Mutex{},
		customers: make(map[string]*domain.Customer),
	}
}

func (r *CustomerRepository) GetCustomers() map[string]*domain.Customer {
	return r.customers
}

func (r *CustomerRepository) Save(ctx context.Context, customer *domain.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.customers[customer.Email] = customer
	return nil
}

func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	customer, exists := r.customers[email]
	if !exists {
		return nil, errors.New("customer not found")
	}

	return customer, nil
}
