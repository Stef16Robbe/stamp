package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Directory != "docs/adr" {
		t.Errorf("Directory = %q, want %q", cfg.Directory, "docs/adr")
	}
}

func TestConfigFileName(t *testing.T) {
	if ConfigFileName != ".stamp.yaml" {
		t.Errorf("ConfigFileName = %q, want %q", ConfigFileName, ".stamp.yaml")
	}
}

func TestConfigSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &Config{
		Directory: "custom/adr/path",
	}

	if err := cfg.Save(tmpDir); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Verify file was created
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Save() did not create config file")
	}

	// Verify content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Config file is empty")
	}

	// Should contain the directory setting
	contentStr := string(content)
	if !contains(contentStr, "custom/adr/path") {
		t.Errorf("Config file doesn't contain directory path, got: %s", contentStr)
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create temp directory structure
	tmpDir, err := os.MkdirTemp("", "stamp-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save config
	originalCfg := &Config{
		Directory: "my/custom/adrs",
	}

	if err := originalCfg.Save(tmpDir); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Change to the temp directory to test Load()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Load config
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loadedCfg.Directory != originalCfg.Directory {
		t.Errorf("Loaded Directory = %q, want %q", loadedCfg.Directory, originalCfg.Directory)
	}
}

func TestFindConfigFileInCurrentDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Resolve symlinks (macOS uses /var -> /private/var)
	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks: %v", err)
	}

	// Create config file
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte("directory: docs/adr\n"), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Change to temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Find config
	found, err := FindConfigFile()
	if err != nil {
		t.Fatalf("FindConfigFile() error: %v", err)
	}

	if found != configPath {
		t.Errorf("FindConfigFile() = %q, want %q", found, configPath)
	}
}

func TestFindConfigFileInParentDir(t *testing.T) {
	// Create temp directory structure: parent/child
	parentDir, err := os.MkdirTemp("", "stamp-config-test-parent-*")
	if err != nil {
		t.Fatalf("Failed to create parent dir: %v", err)
	}
	defer os.RemoveAll(parentDir)

	// Resolve symlinks (macOS uses /var -> /private/var)
	parentDir, err = filepath.EvalSymlinks(parentDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks: %v", err)
	}

	childDir := filepath.Join(parentDir, "child")
	if err := os.Mkdir(childDir, 0755); err != nil {
		t.Fatalf("Failed to create child dir: %v", err)
	}

	// Create config in parent directory
	configPath := filepath.Join(parentDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte("directory: docs/adr\n"), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Change to child directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(childDir); err != nil {
		t.Fatalf("Failed to change to child dir: %v", err)
	}

	// Find config should find it in parent
	found, err := FindConfigFile()
	if err != nil {
		t.Fatalf("FindConfigFile() error: %v", err)
	}

	if found != configPath {
		t.Errorf("FindConfigFile() = %q, want %q", found, configPath)
	}
}

func TestFindConfigFileNotFound(t *testing.T) {
	// Create temp directory without config
	tmpDir, err := os.MkdirTemp("", "stamp-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Find config should fail
	_, err = FindConfigFile()
	if err == nil {
		t.Error("FindConfigFile() expected error when no config exists")
	}
}

func TestLoadNoConfigFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory (no config file)
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Load should fail
	_, err = Load()
	if err == nil {
		t.Error("Load() expected error when no config file exists")
	}
}

func TestADRDirectory(t *testing.T) {
	// Create temp directory structure
	tmpDir, err := os.MkdirTemp("", "stamp-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Resolve symlinks (macOS uses /var -> /private/var)
	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks: %v", err)
	}

	// Create config file
	cfg := &Config{Directory: "docs/adr"}
	if err := cfg.Save(tmpDir); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Change to temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Get ADR directory
	adrDir, err := cfg.ADRDirectory()
	if err != nil {
		t.Fatalf("ADRDirectory() error: %v", err)
	}

	expected := filepath.Join(tmpDir, "docs/adr")
	if adrDir != expected {
		t.Errorf("ADRDirectory() = %q, want %q", adrDir, expected)
	}
}

func TestADRDirectoryFromSubdirectory(t *testing.T) {
	// Create temp directory structure: parent/child
	parentDir, err := os.MkdirTemp("", "stamp-config-test-parent-*")
	if err != nil {
		t.Fatalf("Failed to create parent dir: %v", err)
	}
	defer os.RemoveAll(parentDir)

	// Resolve symlinks (macOS uses /var -> /private/var)
	parentDir, err = filepath.EvalSymlinks(parentDir)
	if err != nil {
		t.Fatalf("Failed to resolve symlinks: %v", err)
	}

	childDir := filepath.Join(parentDir, "src", "components")
	if err := os.MkdirAll(childDir, 0755); err != nil {
		t.Fatalf("Failed to create child dir: %v", err)
	}

	// Create config in parent directory
	cfg := &Config{Directory: "docs/decisions"}
	if err := cfg.Save(parentDir); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Change to nested child directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(childDir); err != nil {
		t.Fatalf("Failed to change to child dir: %v", err)
	}

	// Get ADR directory - should resolve relative to config file location
	adrDir, err := cfg.ADRDirectory()
	if err != nil {
		t.Fatalf("ADRDirectory() error: %v", err)
	}

	expected := filepath.Join(parentDir, "docs/decisions")
	if adrDir != expected {
		t.Errorf("ADRDirectory() = %q, want %q", adrDir, expected)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stamp-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create invalid YAML config
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Change to temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Load should fail with invalid YAML
	_, err = Load()
	if err == nil {
		t.Error("Load() expected error for invalid YAML")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
