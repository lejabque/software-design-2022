package main

// main for two services: run app exchange or app account
// use cobra to choose app (using cli arg)

import (
	"github.com/lejabque/software-design-2022/testcontainers/account"
	"github.com/lejabque/software-design-2022/testcontainers/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Run app",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var exchangeCmd = &cobra.Command{
	Use:   "exchange",
	Short: "Run exchange app",
	Run: func(cmd *cobra.Command, args []string) {
		account.Run(&cliArgs)
	},
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Run account app",
	Run: func(cmd *cobra.Command, args []string) {
		account.Run(&cliArgs)
	},
}

var cliArgs app.CliArgs

func main() {
	rootCmd.PersistentFlags().Uint16VarP(&cliArgs.Port, "port", "p", 8080, "port")

	rootCmd.AddCommand(exchangeCmd, accountCmd)
	rootCmd.Execute()
}
