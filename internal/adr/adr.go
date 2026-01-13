package adr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Status string

const (
	StatusDraft      Status = "Draft"
	StatusProposed   Status = "Proposed"
	StatusAccepted   Status = "Accepted"
	StatusDeprecated Status = "Deprecated"
	StatusSuperseded Status = "Superseded"
	StatusRejected   Status = "Rejected"
)

var ValidStatuses = []Status{
	StatusDraft,
	StatusProposed,
	StatusAccepted,
	StatusDeprecated,
	StatusSuperseded,
	StatusRejected,
}

func ParseStatus(s string) (Status, error) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	for _, status := range ValidStatuses {
		if strings.ToLower(string(status)) == normalized {
			return status, nil
		}
	}
	return "", fmt.Errorf("invalid status: %s (valid: draft, proposed, accepted, deprecated, superseded, rejected)", s)
}

type ADR struct {
	Number       int
	Title        string
	Date         time.Time
	Status       Status
	StatusExtra  []string // Additional lines in status section (links, etc.)
	Context      string
	Decision     string
	Consequences string
	Filename     string
}

func (a *ADR) ToMarkdown() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %d. %s\n\n", a.Number, a.Title))
	sb.WriteString(fmt.Sprintf("Date: %s\n\n", a.Date.Format("2006-01-02")))
	sb.WriteString("## Status\n\n")
	sb.WriteString(string(a.Status))
	sb.WriteString("\n")
	for _, extra := range a.StatusExtra {
		sb.WriteString("\n")
		sb.WriteString(extra)
	}
	if len(a.StatusExtra) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString("## Context\n\n")
	sb.WriteString(a.Context)
	sb.WriteString("\n\n")
	sb.WriteString("## Decision\n\n")
	sb.WriteString(a.Decision)
	sb.WriteString("\n\n")
	sb.WriteString("## Consequences\n\n")
	sb.WriteString(a.Consequences)
	sb.WriteString("\n")

	return sb.String()
}

var (
	titleRegex  = regexp.MustCompile(`^#\s*(\d+)\.\s*(.+)$`)
	dateRegex   = regexp.MustCompile(`^Date:\s*(.+)$`)
	headerRegex = regexp.MustCompile(`^##\s*(.+)$`)
)

func ParseMarkdown(content string) (*ADR, error) {
	adr := &ADR{}
	lines := strings.Split(content, "\n")

	var currentSection string
	var sectionContent []string

	flushSection := func() {
		if currentSection == "" {
			return
		}
		text := strings.TrimSpace(strings.Join(sectionContent, "\n"))
		switch currentSection {
		case "Status":
			statusLines := strings.Split(text, "\n")
			if len(statusLines) > 0 {
				status, err := ParseStatus(statusLines[0])
				if err == nil {
					adr.Status = status
				} else {
					adr.Status = Status(statusLines[0])
				}
				if len(statusLines) > 1 {
					for _, line := range statusLines[1:] {
						if strings.TrimSpace(line) != "" {
							adr.StatusExtra = append(adr.StatusExtra, line)
						}
					}
				}
			}
		case "Context":
			adr.Context = text
		case "Decision":
			adr.Decision = text
		case "Consequences":
			adr.Consequences = text
		}
		sectionContent = nil
	}

	for _, line := range lines {
		if match := titleRegex.FindStringSubmatch(line); match != nil {
			num, _ := strconv.Atoi(match[1])
			adr.Number = num
			adr.Title = strings.TrimSpace(match[2])
			continue
		}

		if match := dateRegex.FindStringSubmatch(line); match != nil {
			dateStr := strings.TrimSpace(match[1])
			if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				adr.Date = t
			}
			continue
		}

		if match := headerRegex.FindStringSubmatch(line); match != nil {
			flushSection()
			currentSection = strings.TrimSpace(match[1])
			continue
		}

		if currentSection != "" {
			sectionContent = append(sectionContent, line)
		}
	}

	flushSection()

	return adr, nil
}

func Slugify(title string) string {
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func FormatFilename(number int, title string) string {
	return fmt.Sprintf("%04d-%s.md", number, Slugify(title))
}
