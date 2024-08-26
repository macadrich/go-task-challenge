package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/macadrich/go-task-challenge/application"
	"github.com/macadrich/go-task-challenge/domain"
	"github.com/macadrich/go-task-challenge/infra"
	"github.com/macadrich/go-task-challenge/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterCommand(t *testing.T) {
	mockKYC := new(mocks.MockKYCService)
	mockKYC.On("ValidateKYC", mock.Anything, mock.Anything).Return(nil)

	customerRepository := infra.NewCustomerRepository()
	customerService := application.NewCustomerService(mockKYC, customerRepository)

	output := new(bytes.Buffer)

	cmd := &cobra.Command{
		Use: "register",
		RunE: func(cmd *cobra.Command, args []string) error {
			customer := &domain.Customer{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Phone:     phone,
				Address:   address,
			}
			ctx := context.Background()

			if err := customerService.RegisterCustomer(ctx, customer); err != nil {
				return err
			}

			cmd.Printf("Customer registered successfully: %s %s\n", customer.FirstName, customer.LastName)

			return nil
		},
	}

	cmd.Flags().StringVar(&firstName, "first-name", "John", "Customer's first name")
	cmd.Flags().StringVar(&lastName, "last-name", "Doe", "Customer's last name")
	cmd.Flags().StringVar(&email, "email", "john.doe@example.com", "Customer's email")
	cmd.Flags().StringVar(&phone, "phone", "1234567890", "Customer's phone")
	cmd.Flags().StringVar(&address, "address", "123 Main St", "Customer's address")

	cmd.SetOut(output)
	err := cmd.Execute()

	assert.NoError(t, err)
	mockKYC.AssertCalled(t, "ValidateKYC", mock.Anything, mock.Anything)

	expectedOutput := "Customer registered successfully: John Doe\n"
	assert.Equal(t, expectedOutput, output.String())
}
