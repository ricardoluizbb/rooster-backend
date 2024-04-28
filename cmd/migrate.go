// cmd/migrate.go
package cmd

import (
	"fmt"
	"time-tracker-backend/database"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		err := database.RunMigrations()
		if err != nil {
			fmt.Println("Error running migrations:", err)
		} else {
			fmt.Println("Migrations executed successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
