package cmd

import (
	"context"

	"github.com/macadrich/go-task-challenge/application"
	"github.com/macadrich/go-task-challenge/domain"
	external "github.com/macadrich/go-task-challenge/external"
	"github.com/macadrich/go-task-challenge/infra"
	"github.com/spf13/cobra"
)

var (
	firstName string
	lastName  string
	email     string
	phone     string
	address   string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Task1 register a new customer",
	Long:  "Task1 register a new customer and validate their KYC information using an external service",
	RunE: func(cmd *cobra.Command, args []string) error {
		externalService := &external.ExternalKYCService{}
		kycAdapter := infra.NewKYCAdapter(externalService)

		customerService := application.NewCustomerService(kycAdapter, customerRepository)

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

func init() {
	registerCmd.Flags().StringVar(&firstName, "first-name", "", "Customer's first name")
	registerCmd.Flags().StringVar(&lastName, "last-name", "", "Customer's last name")
	registerCmd.Flags().StringVar(&email, "email", "", "Customer's email")
	registerCmd.Flags().StringVar(&phone, "phone", "", "Customer's phone number")
	registerCmd.Flags().StringVar(&address, "address", "", "Customer's address")

	registerCmd.MarkFlagRequired("first-name")
	registerCmd.MarkFlagRequired("last-name")
	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("phone")
	registerCmd.MarkFlagRequired("address")

	rootCmd.AddCommand(registerCmd)
}
