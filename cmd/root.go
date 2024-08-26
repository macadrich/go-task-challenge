package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/macadrich/go-task-challenge/infra"
	"github.com/spf13/cobra"
)

var (
	customerRepository *infra.CustomerRepository
)

var rootCmd = &cobra.Command{
	Use:   "go-challege",
	Short: "go-challege CLI",
	Long:  "go-challege CLI example of using DDD pattern and concurrent programming to solve a problem",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if customerRepository == nil {
			customerRepository = infra.NewCustomerRepository()
		}
	},
}

func commandLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting...")
			break
		}

		cmdArgs := strings.Split(input, " ")
		rootCmd.SetArgs(cmdArgs)
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
	}
}

func Execute() {
	commandLoop()
}
