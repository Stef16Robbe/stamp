package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/adr"
	"github.com/stef16robbe/stamp/internal/config"
	"github.com/stef16robbe/stamp/internal/ui"
)

var showCmd = &cobra.Command{
	Use:   "show <number>",
	Short: "Show an ADR",
	Long:  `Displays an Architecture Decision Record with rendered markdown.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid ADR number: %s", args[0])
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

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(80),
		)
		if err != nil {
			return fmt.Errorf("failed to create renderer: %w", err)
		}

		output, err := renderer.Render(a.ToMarkdown())
		if err != nil {
			return fmt.Errorf("failed to render ADR: %w", err)
		}

		// Trim trailing newlines from glamour output
		output = strings.TrimRight(output, "\n")

		// Create title for the frame
		title := fmt.Sprintf(" ADR %04d ", a.Number)
		titleStyle := lipgloss.NewStyle().
			Background(ui.Cyan).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Padding(0, 1)

		statusBadge := ui.RenderStatus(a.Status)

		header := titleStyle.Render(title) + " " + statusBadge

		// Frame the content
		frameStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ui.Gray).
			Padding(0, 1)

		framedContent := frameStyle.Render(output)

		fmt.Println(header)
		fmt.Println(framedContent)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
