package cmd

import (
	"fmt"
	"time"

	rediscache "github.com/macadrich/go-task-challenge/redis-cache"
	"github.com/spf13/cobra"
)

var ttl int

var c = rediscache.NewRedisCache(5)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a key-value pair in the cache",
	Long:  "This command allows you to set a key-value pair in the cache with an optional TTL (Time-To-Live) in seconds.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Error: 'set' requires a key and a value")
			return
		}

		key := args[0]
		value := args[1]

		var ttlDuration time.Duration
		if ttl > 0 {
			ttlDuration = time.Duration(ttl) * time.Second
		} else {
			ttlDuration = 0
		}

		c.Set(key, value, ttlDuration)
		fmt.Printf("Key '%s' set successfully with value '%s'\n", key, value)
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the value of a key from the cache",
	Long:  "This command allows you to get the value of a key stored in the cache.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: 'get' requires a key")
			return
		}

		key := args[0]
		value := c.Get(key)

		if value != nil {
			fmt.Printf("Value for key '%s': %v\n", key, value)
		} else {
			fmt.Printf("Key '%s' not found or expired\n", key)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().IntVarP(&ttl, "ttl", "t", 0, "Time-to-live for the key in seconds")

	rootCmd.AddCommand(getCmd)
}
