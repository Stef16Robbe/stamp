package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/adr"
	"github.com/stef16robbe/stamp/internal/config"
	"github.com/stef16robbe/stamp/internal/ui"
)

var validRelations = map[string]string{
	"supersedes":    "Supersedes",
	"superseded-by": "Superseded by",
	"amends":        "Amends",
	"amended-by":    "Amended by",
	"clarifies":     "Clarifies",
	"clarified-by":  "Clarified by",
}

var linkCmd = &cobra.Command{
	Use:   "link <source> <target> <relation>",
	Short: "Link two ADRs",
	Long: `Creates a link between two ADRs by adding a relationship to the source ADR's status section.

Valid relations: supersedes, superseded-by, amends, amended-by, clarifies, clarified-by

Example:
  stamp link 2 1 supersedes
  # Adds "Supersedes [ADR-0001](0001-...md)" to ADR 0002`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceNum, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid source ADR number: %s", args[0])
		}

		targetNum, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid target ADR number: %s", args[1])
		}

		relation := strings.ToLower(args[2])
		relationDisplay, ok := validRelations[relation]
		if !ok {
			validList := make([]string, 0, len(validRelations))
			for k := range validRelations {
				validList = append(validList, k)
			}
			return fmt.Errorf("invalid relation: %s (valid: %s)", args[2], strings.Join(validList, ", "))
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

		source, err := store.FindByNumber(sourceNum)
		if err != nil {
			return fmt.Errorf("source ADR %04d not found", sourceNum)
		}

		target, err := store.FindByNumber(targetNum)
		if err != nil {
			return fmt.Errorf("target ADR %04d not found", targetNum)
		}

		linkLine := fmt.Sprintf("%s [ADR-%04d](%s)", relationDisplay, target.Number, target.Filename)
		source.StatusExtra = append(source.StatusExtra, linkLine)

		if err := store.Save(source); err != nil {
			return fmt.Errorf("failed to save ADR: %w", err)
		}

		arrow := lipgloss.NewStyle().Foreground(ui.Magenta).Render(" → ")
		adrStyle := lipgloss.NewStyle().Foreground(ui.Cyan).Bold(true)
		fmt.Println(ui.Success("Linked " + adrStyle.Render(fmt.Sprintf("ADR-%04d", sourceNum)) + arrow + adrStyle.Render(fmt.Sprintf("ADR-%04d", targetNum)) + ui.Muted(" ("+relationDisplay+")")))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
