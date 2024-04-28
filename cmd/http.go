// cmd/http.go
package cmd

import (
	"fmt"
	"time-tracker-backend/routes"

	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the HTTP server...")
		routes.SetupRoutes()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
