package adr

import (
	"strings"
	"testing"
	"time"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Status
		wantErr bool
	}{
		{"draft lowercase", "draft", StatusDraft, false},
		{"Draft mixed case", "Draft", StatusDraft, false},
		{"DRAFT uppercase", "DRAFT", StatusDraft, false},
		{"proposed", "proposed", StatusProposed, false},
		{"accepted", "accepted", StatusAccepted, false},
		{"deprecated", "deprecated", StatusDeprecated, false},
		{"superseded", "superseded", StatusSuperseded, false},
		{"rejected", "rejected", StatusRejected, false},
		{"with leading whitespace", "  draft", StatusDraft, false},
		{"with trailing whitespace", "draft  ", StatusDraft, false},
		{"with both whitespace", "  draft  ", StatusDraft, false},
		{"invalid status", "invalid", "", true},
		{"empty string", "", "", true},
		{"partial match", "draf", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple title", "Use PostgreSQL", "use-postgresql"},
		{"with special chars", "Use Go 1.21!", "use-go-121"},
		{"multiple spaces", "Use   multiple   spaces", "use-multiple-spaces"},
		{"leading trailing spaces", "  title  ", "title"},
		{"with hyphens", "already-has-hyphens", "already-has-hyphens"},
		{"multiple hyphens", "has---multiple---hyphens", "has-multiple-hyphens"},
		{"mixed case", "MixedCaseTitle", "mixedcasetitle"},
		{"numbers only", "12345", "12345"},
		{"special chars only", "!@#$%", ""},
		{"unicode chars", "café résumé", "caf-rsum"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatFilename(t *testing.T) {
	tests := []struct {
		name   string
		number int
		title  string
		want   string
	}{
		{"single digit", 1, "Use PostgreSQL", "0001-use-postgresql.md"},
		{"double digit", 12, "Use Redis", "0012-use-redis.md"},
		{"triple digit", 123, "Use Kafka", "0123-use-kafka.md"},
		{"four digit", 1234, "Use RabbitMQ", "1234-use-rabbitmq.md"},
		{"five digit", 12345, "Large Number", "12345-large-number.md"},
		{"zero", 0, "Zero Test", "0000-zero-test.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFilename(tt.number, tt.title)
			if got != tt.want {
				t.Errorf("FormatFilename() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatFilenameSpecialChars(t *testing.T) {
	got := FormatFilename(1, "Record Architecture Decisions!")
	want := "0001-record-architecture-decisions.md"
	if got != want {
		t.Errorf("FormatFilename() = %q, want %q", got, want)
	}
}

func TestADRToMarkdown(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		adr  ADR
		want string
	}{
		{
			name: "basic ADR",
			adr: ADR{
				Number:       1,
				Title:        "Record Architecture Decisions",
				Date:         date,
				Status:       StatusAccepted,
				Context:      "We need to record decisions.",
				Decision:     "We will use ADRs.",
				Consequences: "Easy to track decisions.",
			},
			want: `# 1. Record Architecture Decisions

Date: 2024-01-15

## Status

Accepted

## Context

We need to record decisions.

## Decision

We will use ADRs.

## Consequences

Easy to track decisions.
`,
		},
		{
			name: "ADR with status extra",
			adr: ADR{
				Number:       2,
				Title:        "Use PostgreSQL",
				Date:         date,
				Status:       StatusSuperseded,
				StatusExtra:  []string{"Superseded by [ADR 5](0005-use-mysql.md)"},
				Context:      "Need a database.",
				Decision:     "Use PostgreSQL.",
				Consequences: "Need to learn PostgreSQL.",
			},
			want: `# 2. Use PostgreSQL

Date: 2024-01-15

## Status

Superseded

Superseded by [ADR 5](0005-use-mysql.md)

## Context

Need a database.

## Decision

Use PostgreSQL.

## Consequences

Need to learn PostgreSQL.
`,
		},
		{
			name: "ADR with multiple status extra lines",
			adr: ADR{
				Number:       3,
				Title:        "Use Microservices",
				Date:         date,
				Status:       StatusAccepted,
				StatusExtra:  []string{"Amends [ADR 1](0001-record-architecture-decisions.md)", "Clarifies [ADR 2](0002-use-postgresql.md)"},
				Context:      "Growing complexity.",
				Decision:     "Split into microservices.",
				Consequences: "More operational overhead.",
			},
			want: `# 3. Use Microservices

Date: 2024-01-15

## Status

Accepted

Amends [ADR 1](0001-record-architecture-decisions.md)
Clarifies [ADR 2](0002-use-postgresql.md)

## Context

Growing complexity.

## Decision

Split into microservices.

## Consequences

More operational overhead.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.adr.ToMarkdown()
			if got != tt.want {
				t.Errorf("ToMarkdown() mismatch:\ngot:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantNumber  int
		wantTitle   string
		wantStatus  Status
		wantContext string
	}{
		{
			name: "basic ADR",
			input: `# 1. Record Architecture Decisions

Date: 2024-01-15

## Status

Accepted

## Context

We need to record decisions.

## Decision

We will use ADRs.

## Consequences

Easy to track decisions.
`,
			wantNumber:  1,
			wantTitle:   "Record Architecture Decisions",
			wantStatus:  StatusAccepted,
			wantContext: "We need to record decisions.",
		},
		{
			name: "ADR with status extra",
			input: `# 2. Use PostgreSQL

Date: 2024-01-15

## Status

Superseded

Superseded by [ADR 5](0005-use-mysql.md)

## Context

Need a database.

## Decision

Use PostgreSQL.

## Consequences

Need to learn PostgreSQL.
`,
			wantNumber:  2,
			wantTitle:   "Use PostgreSQL",
			wantStatus:  StatusSuperseded,
			wantContext: "Need a database.",
		},
		{
			name: "title without space after hash",
			input: `#1. No Space Title

Date: 2024-01-15

## Status

Draft

## Context

Test context.

## Decision

Test decision.

## Consequences

Test consequences.
`,
			wantNumber:  1,
			wantTitle:   "No Space Title",
			wantStatus:  StatusDraft,
			wantContext: "Test context.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMarkdown(tt.input)
			if err != nil {
				t.Fatalf("ParseMarkdown() error = %v", err)
			}
			if got.Number != tt.wantNumber {
				t.Errorf("Number = %d, want %d", got.Number, tt.wantNumber)
			}
			if got.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", got.Title, tt.wantTitle)
			}
			if got.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", got.Status, tt.wantStatus)
			}
			if got.Context != tt.wantContext {
				t.Errorf("Context = %q, want %q", got.Context, tt.wantContext)
			}
		})
	}
}

func TestParseMarkdownStatusExtra(t *testing.T) {
	input := `# 2. Use PostgreSQL

Date: 2024-01-15

## Status

Superseded

Superseded by [ADR 5](0005-use-mysql.md)
Amended by [ADR 3](0003-add-read-replicas.md)

## Context

Need a database.

## Decision

Use PostgreSQL.

## Consequences

Need to learn PostgreSQL.
`

	got, err := ParseMarkdown(input)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if len(got.StatusExtra) != 2 {
		t.Fatalf("StatusExtra length = %d, want 2", len(got.StatusExtra))
	}

	if !strings.Contains(got.StatusExtra[0], "Superseded by") {
		t.Errorf("StatusExtra[0] = %q, want to contain 'Superseded by'", got.StatusExtra[0])
	}

	if !strings.Contains(got.StatusExtra[1], "Amended by") {
		t.Errorf("StatusExtra[1] = %q, want to contain 'Amended by'", got.StatusExtra[1])
	}
}

func TestParseMarkdownRoundTrip(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	original := &ADR{
		Number:       1,
		Title:        "Record Architecture Decisions",
		Date:         date,
		Status:       StatusAccepted,
		StatusExtra:  []string{"Amends [ADR 0](0000-initial.md)"},
		Context:      "We need to record architectural decisions.",
		Decision:     "We will use ADRs.",
		Consequences: "Decisions will be documented.",
	}

	markdown := original.ToMarkdown()
	parsed, err := ParseMarkdown(markdown)
	if err != nil {
		t.Fatalf("ParseMarkdown() error = %v", err)
	}

	if parsed.Number != original.Number {
		t.Errorf("Number = %d, want %d", parsed.Number, original.Number)
	}
	if parsed.Title != original.Title {
		t.Errorf("Title = %q, want %q", parsed.Title, original.Title)
	}
	if parsed.Status != original.Status {
		t.Errorf("Status = %q, want %q", parsed.Status, original.Status)
	}
	if parsed.Context != original.Context {
		t.Errorf("Context = %q, want %q", parsed.Context, original.Context)
	}
	if parsed.Decision != original.Decision {
		t.Errorf("Decision = %q, want %q", parsed.Decision, original.Decision)
	}
	if parsed.Consequences != original.Consequences {
		t.Errorf("Consequences = %q, want %q", parsed.Consequences, original.Consequences)
	}
	if len(parsed.StatusExtra) != len(original.StatusExtra) {
		t.Errorf("StatusExtra length = %d, want %d", len(parsed.StatusExtra), len(original.StatusExtra))
	}
	if !parsed.Date.Equal(original.Date) {
		t.Errorf("Date = %v, want %v", parsed.Date, original.Date)
	}
}

func TestValidStatuses(t *testing.T) {
	expected := []Status{
		StatusDraft,
		StatusProposed,
		StatusAccepted,
		StatusDeprecated,
		StatusSuperseded,
		StatusRejected,
	}

	if len(ValidStatuses) != len(expected) {
		t.Errorf("ValidStatuses length = %d, want %d", len(ValidStatuses), len(expected))
	}

	for i, status := range expected {
		if ValidStatuses[i] != status {
			t.Errorf("ValidStatuses[%d] = %q, want %q", i, ValidStatuses[i], status)
		}
	}
}
