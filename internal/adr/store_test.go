package adr

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	store := NewStore("/tmp/test-adrs")
	if store.Directory != "/tmp/test-adrs" {
		t.Errorf("Directory = %q, want %q", store.Directory, "/tmp/test-adrs")
	}
}

func TestStoreList(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "stamp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Test empty directory
	adrs, err := store.List()
	if err != nil {
		t.Fatalf("List() error on empty dir: %v", err)
	}
	if len(adrs) != 0 {
		t.Errorf("List() on empty dir returned %d ADRs, want 0", len(adrs))
	}

	// Create some ADR files
	adr1 := &ADR{
		Number:       1,
		Title:        "First Decision",
		Date:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Status:       StatusAccepted,
		Context:      "Context 1",
		Decision:     "Decision 1",
		Consequences: "Consequences 1",
	}
	adr2 := &ADR{
		Number:       2,
		Title:        "Second Decision",
		Date:         time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
		Status:       StatusDraft,
		Context:      "Context 2",
		Decision:     "Decision 2",
		Consequences: "Consequences 2",
	}

	if err := store.Save(adr1); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	if err := store.Save(adr2); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Create a non-ADR file that should be ignored
	nonADRPath := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(nonADRPath, []byte("# README"), 0644); err != nil {
		t.Fatalf("Failed to create non-ADR file: %v", err)
	}

	// Create a directory that should be ignored
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Test listing
	adrs, err = store.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(adrs) != 2 {
		t.Errorf("List() returned %d ADRs, want 2", len(adrs))
	}

	// Verify sorting by number
	if adrs[0].Number != 1 {
		t.Errorf("First ADR number = %d, want 1", adrs[0].Number)
	}
	if adrs[1].Number != 2 {
		t.Errorf("Second ADR number = %d, want 2", adrs[1].Number)
	}
}

func TestStoreListNonExistentDirectory(t *testing.T) {
	store := NewStore("/nonexistent/path/that/does/not/exist")
	_, err := store.List()
	if err == nil {
		t.Error("List() expected error for non-existent directory")
	}
}

func TestStoreLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Create an ADR file
	content := `# 1. Test Decision

Date: 2024-01-15

## Status

Accepted

## Context

Test context.

## Decision

Test decision.

## Consequences

Test consequences.
`
	filename := "0001-test-decision.md"
	filePath := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test loading
	adr, err := store.Load(filename)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if adr.Number != 1 {
		t.Errorf("Number = %d, want 1", adr.Number)
	}
	if adr.Title != "Test Decision" {
		t.Errorf("Title = %q, want %q", adr.Title, "Test Decision")
	}
	if adr.Status != StatusAccepted {
		t.Errorf("Status = %q, want %q", adr.Status, StatusAccepted)
	}
	if adr.Filename != filename {
		t.Errorf("Filename = %q, want %q", adr.Filename, filename)
	}
}

func TestStoreLoadNonExistent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)
	_, err = store.Load("nonexistent.md")
	if err == nil {
		t.Error("Load() expected error for non-existent file")
	}
}

func TestStoreSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	adr := &ADR{
		Number:       1,
		Title:        "Test Decision",
		Date:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Status:       StatusAccepted,
		Context:      "Test context.",
		Decision:     "Test decision.",
		Consequences: "Test consequences.",
	}

	if err := store.Save(adr); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Verify filename was set
	expectedFilename := "0001-test-decision.md"
	if adr.Filename != expectedFilename {
		t.Errorf("Filename = %q, want %q", adr.Filename, expectedFilename)
	}

	// Verify file exists
	filePath := filepath.Join(tmpDir, expectedFilename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Save() did not create file")
	}

	// Verify content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Saved file is empty")
	}
}

func TestStoreSaveWithExistingFilename(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	adr := &ADR{
		Number:       1,
		Title:        "Test Decision",
		Date:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Status:       StatusAccepted,
		Filename:     "custom-filename.md",
		Context:      "Test context.",
		Decision:     "Test decision.",
		Consequences: "Test consequences.",
	}

	if err := store.Save(adr); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Verify custom filename was preserved
	if adr.Filename != "custom-filename.md" {
		t.Errorf("Filename = %q, want %q", adr.Filename, "custom-filename.md")
	}

	// Verify file exists with custom name
	filePath := filepath.Join(tmpDir, "custom-filename.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Save() did not create file with custom filename")
	}
}

func TestStoreNextNumber(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Test empty directory
	num, err := store.NextNumber()
	if err != nil {
		t.Fatalf("NextNumber() error on empty dir: %v", err)
	}
	if num != 1 {
		t.Errorf("NextNumber() on empty dir = %d, want 1", num)
	}

	// Create some files
	files := []string{
		"0001-first.md",
		"0002-second.md",
		"0005-fifth.md", // Gap in sequence
	}
	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		content := `# 1. Test

Date: 2024-01-15

## Status

Draft

## Context

Test

## Decision

Test

## Consequences

Test
`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", f, err)
		}
	}

	// Test with existing files
	num, err = store.NextNumber()
	if err != nil {
		t.Fatalf("NextNumber() error: %v", err)
	}
	if num != 6 {
		t.Errorf("NextNumber() = %d, want 6", num)
	}
}

func TestStoreNextNumberNonExistentDirectory(t *testing.T) {
	store := NewStore("/nonexistent/path")
	num, err := store.NextNumber()
	if err != nil {
		t.Fatalf("NextNumber() should not error for non-existent directory: %v", err)
	}
	if num != 1 {
		t.Errorf("NextNumber() for non-existent dir = %d, want 1", num)
	}
}

func TestStoreNextNumberIgnoresNonADRFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Create non-ADR files
	nonADRFiles := []string{
		"README.md",
		"notes.txt",
		"1-invalid-format.md", // Missing leading zeros
	}
	for _, f := range nonADRFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", f, err)
		}
	}

	num, err := store.NextNumber()
	if err != nil {
		t.Fatalf("NextNumber() error: %v", err)
	}
	if num != 1 {
		t.Errorf("NextNumber() = %d, want 1 (should ignore non-ADR files)", num)
	}
}

func TestStoreFindByNumber(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Create ADRs
	adr1 := &ADR{
		Number:       1,
		Title:        "First",
		Date:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Status:       StatusAccepted,
		Context:      "Context",
		Decision:     "Decision",
		Consequences: "Consequences",
	}
	adr2 := &ADR{
		Number:       2,
		Title:        "Second",
		Date:         time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
		Status:       StatusDraft,
		Context:      "Context",
		Decision:     "Decision",
		Consequences: "Consequences",
	}

	store.Save(adr1)
	store.Save(adr2)

	// Test finding existing ADR
	found, err := store.FindByNumber(1)
	if err != nil {
		t.Fatalf("FindByNumber(1) error: %v", err)
	}
	if found.Number != 1 {
		t.Errorf("Found ADR number = %d, want 1", found.Number)
	}
	if found.Title != "First" {
		t.Errorf("Found ADR title = %q, want %q", found.Title, "First")
	}

	// Test finding second ADR
	found, err = store.FindByNumber(2)
	if err != nil {
		t.Fatalf("FindByNumber(2) error: %v", err)
	}
	if found.Number != 2 {
		t.Errorf("Found ADR number = %d, want 2", found.Number)
	}

	// Test finding non-existent ADR
	_, err = store.FindByNumber(99)
	if err == nil {
		t.Error("FindByNumber(99) expected error for non-existent ADR")
	}
	if !os.IsNotExist(err) {
		t.Errorf("FindByNumber(99) error = %v, want os.ErrNotExist", err)
	}
}

func TestFilenameRegex(t *testing.T) {
	tests := []struct {
		filename string
		match    bool
	}{
		{"0001-test.md", true},
		{"0012-some-title.md", true},
		{"1234-another-one.md", true},
		{"0001-a.md", true},
		{"README.md", false},
		{"1-no-padding.md", false},
		{"0001-test.txt", false},
		{"0001.md", false},
		{"test-0001.md", false},
		{"0001-test.md.bak", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := filenameRegex.MatchString(tt.filename)
			if got != tt.match {
				t.Errorf("filenameRegex.MatchString(%q) = %v, want %v", tt.filename, got, tt.match)
			}
		})
	}
}
