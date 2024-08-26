package cmd

import (
	"context"
	"fmt"

	"github.com/macadrich/go-task-challenge/application"
	"github.com/macadrich/go-task-challenge/constants"
	external "github.com/macadrich/go-task-challenge/external"
	"github.com/macadrich/go-task-challenge/infra"
	"github.com/spf13/cobra"
)

var (
	verifyEmail string
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Task2 verify a customer",
	Long:  "Task2 verify a customer information using an external service.",
	RunE: func(cmd *cobra.Command, args []string) error {
		externalService := &external.ExternalKYCService{}
		kycAdapter := infra.NewKYCAdapter(externalService)

		customerService := application.NewCustomerService(kycAdapter, customerRepository)

		ctx := context.Background()

		customer, err := customerRepository.FindByEmail(ctx, verifyEmail)
		if err != nil {
			return fmt.Errorf("failed to find customer: %w", err)
		}

		if err := customerService.VerifyRegisteredCustomer(ctx, constants.NumberOfRoutines, customer); err != nil {
			return fmt.Errorf("failed to verify customer: %w", err)
		}

		cmd.Printf("Customer verify successfully: %s %s\n", customer.FirstName, customer.LastName)

		return nil
	},
}

func init() {
	verifyCmd.Flags().StringVar(&verifyEmail, "email", "", "Customer email")
	verifyCmd.MarkFlagRequired("email")
	rootCmd.AddCommand(verifyCmd)
}
