package adr

import "time"

func NewADR(number int, title string) *ADR {
	return &ADR{
		Number:       number,
		Title:        title,
		Date:         time.Now(),
		Status:       StatusDraft,
		Context:      "[Why is this decision needed?]",
		Decision:     "[What was decided?]",
		Consequences: "[What are the implications?]",
	}
}
