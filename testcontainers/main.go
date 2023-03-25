package main

// main for two services: run app exchange or app account
// use cobra to choose app (using cli arg)

import (
	"github.com/lejabque/software-design-2022/testcontainers/internal/account"
	"github.com/lejabque/software-design-2022/testcontainers/internal/exchange"
	"github.com/lejabque/software-design-2022/testcontainers/internal/lib"
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
		exchange.Run(&cliArgs)
	},
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Run account app",
	Run: func(cmd *cobra.Command, args []string) {
		account.Run(&cliArgs)
	},
}

var cliArgs lib.CliArgs

func main() {
	rootCmd.PersistentFlags().Uint16VarP(&cliArgs.Port, "port", "p", 8080, "port")

	accountCmd.PersistentFlags().StringVarP(&cliArgs.ExchangeEndpoint, "exchange", "e", "localhost:8081", "exchange endpoint")

	rootCmd.AddCommand(exchangeCmd, accountCmd)
	rootCmd.Execute()
}
