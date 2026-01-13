package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stamp",
	Short: "Manage Architecture Decision Records",
	Long:  `Stamp is a CLI tool for creating and managing Architecture Decision Records (ADRs).`,
}

func Execute() error {
	return rootCmd.Execute()
}
