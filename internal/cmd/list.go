package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/adr"
	"github.com/stef16robbe/stamp/internal/config"
	"github.com/stef16robbe/stamp/internal/ui"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ADRs",
	Long:  `Lists all Architecture Decision Records with their status and date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		dir, err := cfg.ADRDirectory()
		if err != nil {
			return err
		}

		store := adr.NewStore(dir)

		adrs, err := store.List()
		if err != nil {
			return fmt.Errorf("failed to list ADRs: %w", err)
		}

		if len(adrs) == 0 {
			fmt.Println(ui.Warning("No ADRs found. Create one with 'stamp new <title>'"))
			return nil
		}

		rows := make([][]string, len(adrs))
		for i, a := range adrs {
			rows[i] = []string{
				fmt.Sprintf("%04d", a.Number),
				a.Title,
				ui.RenderStatus(a.Status),
				a.Date.Format("2006-01-02"),
			}
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(ui.Gray)).
			Headers("NUM", "TITLE", "STATUS", "DATE").
			Rows(rows...).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == table.HeaderRow {
					return lipgloss.NewStyle().
						Bold(true).
						Foreground(ui.Cyan).
						Padding(0, 1)
				}
				return lipgloss.NewStyle().Padding(0, 1)
			})

		fmt.Println(t)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
