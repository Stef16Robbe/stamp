package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/adr"
	"github.com/stef16robbe/stamp/internal/config"
	"github.com/stef16robbe/stamp/internal/ui"
)

var openEditor bool

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Create a new ADR",
	Long:  `Creates a new Architecture Decision Record with the next available number.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		dir, err := cfg.ADRDirectory()
		if err != nil {
			return err
		}

		store := adr.NewStore(dir)

		nextNum, err := store.NextNumber()
		if err != nil {
			return fmt.Errorf("failed to determine next ADR number: %w", err)
		}

		newADR := adr.NewADR(nextNum, title)

		if err := store.Save(newADR); err != nil {
			return fmt.Errorf("failed to save ADR: %w", err)
		}

		fmt.Println(ui.Success("Created " + ui.Muted(newADR.Filename)))

		if openEditor {
			editor := os.Getenv("VISUAL")
			if editor == "" {
				editor = os.Getenv("EDITOR")
			}
			if editor == "" {
				return fmt.Errorf("no editor configured (set $EDITOR or $VISUAL)")
			}

			filePath := filepath.Join(dir, newADR.Filename)
			editorCmd := exec.Command(editor, filePath)
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr

			if err := editorCmd.Run(); err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}
		}

		return nil
	},
}

func init() {
	newCmd.Flags().BoolVarP(&openEditor, "editor", "e", false, "Open the new ADR in $EDITOR")
	rootCmd.AddCommand(newCmd)
}
