package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stef16robbe/stamp/internal/adr"
	"github.com/stef16robbe/stamp/internal/config"
)

var graphFormat string

// Link represents a relationship between two ADRs
type Link struct {
	Source   int
	Target   int
	Relation string
}

// linkRegex matches lines like "Supersedes [ADR-0001](0001-title.md)"
var linkRegex = regexp.MustCompile(`^(Supersedes|Superseded by|Amends|Amended by|Clarifies|Clarified by)\s+\[ADR-(\d+)\]`)

// parseLinks extracts links from an ADR's StatusExtra field
func parseLinks(a *adr.ADR) []Link {
	var links []Link
	for _, line := range a.StatusExtra {
		match := linkRegex.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		relation := match[1]
		targetNum, err := strconv.Atoi(match[2])
		if err != nil {
			continue
		}
		links = append(links, Link{
			Source:   a.Number,
			Target:   targetNum,
			Relation: relation,
		})
	}
	return links
}

// relationToArrow maps relations to Mermaid arrow styles
var relationToArrow = map[string]string{
	"Supersedes":    "-->|supersedes|",
	"Superseded by": "-->|superseded by|",
	"Amends":        "-.->|amends|",
	"Amended by":    "-.->|amended by|",
	"Clarifies":     "-.->|clarifies|",
	"Clarified by":  "-.->|clarified by|",
}

// statusToStyle maps statuses to Mermaid node styles
var statusToStyle = map[adr.Status]string{
	adr.StatusDraft:      ":::draft",
	adr.StatusProposed:   ":::proposed",
	adr.StatusAccepted:   ":::accepted",
	adr.StatusDeprecated: ":::deprecated",
	adr.StatusSuperseded: ":::superseded",
	adr.StatusRejected:   ":::rejected",
}

func generateMermaid(adrs []*adr.ADR) string {
	var sb strings.Builder

	sb.WriteString("graph TD\n")

	// Define style classes
	sb.WriteString("    classDef draft fill:#6b7280,stroke:#374151\n")
	sb.WriteString("    classDef proposed fill:#3b82f6,stroke:#1d4ed8\n")
	sb.WriteString("    classDef accepted fill:#22c55e,stroke:#15803d\n")
	sb.WriteString("    classDef deprecated fill:#f59e0b,stroke:#d97706\n")
	sb.WriteString("    classDef superseded fill:#a855f7,stroke:#7e22ce\n")
	sb.WriteString("    classDef rejected fill:#ef4444,stroke:#b91c1c\n")
	sb.WriteString("\n")

	// Create nodes for each ADR
	for _, a := range adrs {
		nodeID := fmt.Sprintf("ADR%d", a.Number)
		// Truncate long titles
		title := a.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		label := fmt.Sprintf("%04d: %s", a.Number, title)
		style := statusToStyle[a.Status]
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]%s\n", nodeID, label, style))
	}

	sb.WriteString("\n")

	// Create edges - only use forward relations to avoid duplicates
	seen := make(map[string]bool)
	for _, a := range adrs {
		links := parseLinks(a)
		for _, link := range links {
			// Only output forward relations (supersedes, amends, clarifies)
			// Skip reverse relations to avoid duplicate edges
			if strings.HasSuffix(link.Relation, " by") {
				continue
			}

			edgeKey := fmt.Sprintf("%d-%d", link.Source, link.Target)
			if seen[edgeKey] {
				continue
			}
			seen[edgeKey] = true

			arrow := relationToArrow[link.Relation]
			if arrow == "" {
				arrow = "-->"
			}
			sb.WriteString(fmt.Sprintf("    ADR%d %s ADR%d\n", link.Source, arrow, link.Target))
		}
	}

	return sb.String()
}

func generateDot(adrs []*adr.ADR) string {
	var sb strings.Builder

	sb.WriteString("digraph ADRs {\n")
	sb.WriteString("    rankdir=TB;\n")
	sb.WriteString("    node [shape=box, style=rounded];\n")
	sb.WriteString("\n")

	// Status colors for DOT
	statusColor := map[adr.Status]string{
		adr.StatusDraft:      "#6b7280",
		adr.StatusProposed:   "#3b82f6",
		adr.StatusAccepted:   "#22c55e",
		adr.StatusDeprecated: "#f59e0b",
		adr.StatusSuperseded: "#a855f7",
		adr.StatusRejected:   "#ef4444",
	}

	// Create nodes
	for _, a := range adrs {
		title := a.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		label := fmt.Sprintf("%04d: %s", a.Number, title)
		color := statusColor[a.Status]
		if color == "" {
			color = "#6b7280"
		}
		sb.WriteString(fmt.Sprintf("    ADR%d [label=\"%s\", fillcolor=\"%s\", style=\"filled,rounded\", fontcolor=\"white\"];\n",
			a.Number, label, color))
	}

	sb.WriteString("\n")

	// Create edges
	seen := make(map[string]bool)
	for _, a := range adrs {
		links := parseLinks(a)
		for _, link := range links {
			if strings.HasSuffix(link.Relation, " by") {
				continue
			}

			edgeKey := fmt.Sprintf("%d-%d", link.Source, link.Target)
			if seen[edgeKey] {
				continue
			}
			seen[edgeKey] = true

			style := "solid"
			if link.Relation == "Amends" || link.Relation == "Clarifies" {
				style = "dashed"
			}
			sb.WriteString(fmt.Sprintf("    ADR%d -> ADR%d [label=\"%s\", style=%s];\n",
				link.Source, link.Target, strings.ToLower(link.Relation), style))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Generate a visual graph of ADR relationships",
	Long: `Generate a graph showing ADR relationships in Mermaid or Graphviz DOT format.

Mermaid graphs can be rendered directly in GitHub markdown or using the Mermaid CLI.
DOT graphs can be rendered using Graphviz (e.g., dot -Tpng graph.dot -o graph.png).

Examples:
  stamp graph                     # Output Mermaid format
  stamp graph --format mermaid    # Output Mermaid format (explicit)
  stamp graph --format dot        # Output Graphviz DOT format
  stamp graph > docs/adr-graph.md # Save to file`,
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
			return fmt.Errorf("no ADRs found")
		}

		var output string
		switch graphFormat {
		case "mermaid":
			output = generateMermaid(adrs)
		case "dot":
			output = generateDot(adrs)
		default:
			return fmt.Errorf("invalid format: %s (valid: mermaid, dot)", graphFormat)
		}

		fmt.Print(output)
		return nil
	},
}

func init() {
	graphCmd.Flags().StringVarP(&graphFormat, "format", "f", "mermaid", "Output format: mermaid or dot")
	rootCmd.AddCommand(graphCmd)
}
