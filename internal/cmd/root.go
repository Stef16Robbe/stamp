package cmd

import (
	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "stamp",
	Short: "Manage Architecture Decision Records",
	Long:  `Stamp is a CLI tool for creating and managing Architecture Decision Records (ADRs).`,
}

func init() {
	rootCmd.Version = Version
}

func Execute() error {
	return rootCmd.Execute()
}
