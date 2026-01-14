package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of stamp.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("stamp %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
