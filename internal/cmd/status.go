package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/adr"
	"github.com/stef16robbe/stamp/internal/config"
	"github.com/stef16robbe/stamp/internal/ui"
)

var statusCmd = &cobra.Command{
	Use:   "status <number> <status>",
	Short: "Update the status of an ADR",
	Long: `Updates the status of an Architecture Decision Record.

Valid statuses: draft, proposed, accepted, deprecated, superseded, rejected`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ADR number: %s", args[0])
		}

		newStatus, err := adr.ParseStatus(args[1])
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		dir, err := cfg.ADRDirectory()
		if err != nil {
			return err
		}

		store := adr.NewStore(dir)

		a, err := store.FindByNumber(num)
		if err != nil {
			return fmt.Errorf("ADR %04d not found", num)
		}

		oldStatus := a.Status
		a.Status = newStatus

		if err := store.Save(a); err != nil {
			return fmt.Errorf("failed to save ADR: %w", err)
		}

		fmt.Println(ui.Success(fmt.Sprintf("Updated ADR %04d: ", num)) + ui.RenderStatusTransition(oldStatus, newStatus))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
