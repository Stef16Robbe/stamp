package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/adr"
	"github.com/stef16robbe/stamp/internal/config"
)

var editCmd = &cobra.Command{
	Use:   "edit <number>",
	Short: "Edit an ADR in your editor",
	Long:  `Opens an Architecture Decision Record in $EDITOR (or $VISUAL).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ADR number: %s", args[0])
		}

		editor := os.Getenv("VISUAL")
		if editor == "" {
			editor = os.Getenv("EDITOR")
		}
		if editor == "" {
			return fmt.Errorf("no editor configured (set $EDITOR or $VISUAL)")
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

		filePath := filepath.Join(dir, a.Filename)
		editorCmd := exec.Command(editor, filePath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		return editorCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
