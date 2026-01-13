package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/adr"
	"github.com/stef16robbe/stamp/internal/config"
	"github.com/stef16robbe/stamp/internal/ui"
)

var initDirectory string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new ADR directory",
	Long:  `Creates the ADR directory, .stamp.yaml configuration file, and an initial ADR explaining the practice.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		cfg := config.DefaultConfig()
		if initDirectory != "" {
			cfg.Directory = initDirectory
		}

		adrPath := fmt.Sprintf("%s/%s", cwd, cfg.Directory)
		if err := os.MkdirAll(adrPath, 0755); err != nil {
			return fmt.Errorf("failed to create ADR directory: %w", err)
		}

		if err := cfg.Save(cwd); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		store := adr.NewStore(adrPath)
		initialADR := &adr.ADR{
			Number: 1,
			Title:  "Record architecture decisions",
			Date:   time.Now(),
			Status: adr.StatusAccepted,
			Context: `We need to record the architectural decisions made on this project so that future
team members (and our future selves) can understand the reasoning behind our choices.

Without documented decisions, teams often revisit the same discussions, forget why
certain approaches were chosen, or make changes that conflict with earlier decisions.`,
			Decision: `We will use Architecture Decision Records (ADRs) as described by Michael Nygard
in his article [Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions).

See [adr.github.io](https://adr.github.io/) for more information on ADRs.`,
			Consequences: `- All significant architecture decisions will be documented in this directory
- ADRs are numbered sequentially and never deleted (superseded instead)
- Each ADR describes the context, decision, and consequences
- Team members should review existing ADRs before proposing conflicting changes

---

*This project uses [stamp](https://github.com/stef16robbe/stamp) for managing ADRs.*`,
		}

		if err := store.Save(initialADR); err != nil {
			return fmt.Errorf("failed to create initial ADR: %w", err)
		}

		fmt.Println(ui.Success("Initialized ADR directory at " + ui.Bold(cfg.Directory)))
		fmt.Println(ui.Success("Configuration saved to " + ui.Muted(config.ConfigFileName)))
		fmt.Println(ui.Success("Created " + ui.Muted(initialADR.Filename)))

		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&initDirectory, "directory", "d", "", "ADR directory (default: docs/adr)")
	rootCmd.AddCommand(initCmd)
}
