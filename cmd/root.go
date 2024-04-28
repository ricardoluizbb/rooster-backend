// cmd/root.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "time-tracker",
	Short: "Time Tracker CLI",
	Long:  "A CLI for time tracking",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'help' to see available commands.")
	},
}

func Execute() {
	viper.SetEnvPrefix("CLOWN")
	viper.AutomaticEnv() // read in environment variables that match

	if err := rootCmd.Execute(); err != nil {
		// retornar erros
	}
}
