package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadLines(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	content := `# Comment line
line1
  line2  

# Another comment
line3
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	lines, err := ReadLines(filePath)
	if err != nil {
		t.Fatalf("ReadLines returned error: %v", err)
	}

	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(lines))
	}
	for i, v := range expected {
		if lines[i] != v {
			t.Errorf("Line %d expected %q, got %q", i, v, lines[i])
		}
	}
}

func TestLoadEnv(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	// Set a pre-existing env var to test non-overwrite behavior
	_ = os.Setenv("TEST_EXISTING_VAR", "original_value")
	defer os.Unsetenv("TEST_EXISTING_VAR")
	defer os.Unsetenv("TEST_NEW_VAR")
	defer os.Unsetenv("TEST_QUOTED_VAR")
	defer os.Unsetenv("TEST_SINGLE_QUOTED_VAR")

	content := `# Test .env file
TEST_NEW_VAR=hello_world
TEST_EXISTING_VAR=should_not_overwrite
TEST_QUOTED_VAR="quoted value with spaces"
TEST_SINGLE_QUOTED_VAR='single quoted value'
INVALID_LINE_NO_EQUALS
# Comment line
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write .env file: %v", err)
	}

	LoadEnv(envPath)

	if val := os.Getenv("TEST_NEW_VAR"); val != "hello_world" {
		t.Errorf("Expected TEST_NEW_VAR='hello_world', got %q", val)
	}
	if val := os.Getenv("TEST_EXISTING_VAR"); val != "original_value" {
		t.Errorf("Expected TEST_EXISTING_VAR='original_value', got %q", val)
	}
	if val := os.Getenv("TEST_QUOTED_VAR"); val != "quoted value with spaces" {
		t.Errorf("Expected TEST_QUOTED_VAR='quoted value with spaces', got %q", val)
	}
	if val := os.Getenv("TEST_SINGLE_QUOTED_VAR"); val != "single quoted value" {
		t.Errorf("Expected TEST_SINGLE_QUOTED_VAR='single quoted value', got %q", val)
	}

	// Test missing file does not panic or error
	LoadEnv(filepath.Join(tmpDir, "non_existent.env"))
}
