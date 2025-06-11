package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestMain sets up and tears down test environment
func TestMain(m *testing.M) {
	// Build the binary for testing
	cmd := exec.Command("go", "build", "-o", "newestfiles_test_binary", ".")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to build test binary: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Clean up
	os.Remove("newestfiles_test_binary")

	os.Exit(code)
}

// createTestFiles creates test files with specific modification times
func createTestFiles(t *testing.T, dir string) {
	files := []struct {
		name    string
		content string
		age     time.Duration // how old the file should be
	}{
		{"newest.go", "package main", 0},                 // newest
		{"middle.txt", "hello world", time.Hour},         // middle
		{"oldest.md", "# README", 2 * time.Hour},         // oldest
		{"other.py", "print('hello')", 30 * time.Minute}, // should not match .go/.txt/.md
		{"another.go", "// comment", 45 * time.Minute},   // second newest .go
	}

	now := time.Now()
	for _, file := range files {
		path := filepath.Join(dir, file.name)
		err := ioutil.WriteFile(path, []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}

		// Set modification time
		modTime := now.Add(-file.age)
		err = os.Chtimes(path, modTime, modTime)
		if err != nil {
			t.Fatalf("Failed to set modification time for %s: %v", path, err)
		}
	}
}

func TestPlainTextOutput(t *testing.T) {
	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "newestfiles_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	createTestFiles(t, tmpDir)

	// Change to test directory
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Run the program without -j flag
	cmd := exec.Command(filepath.Join(oldDir, "newestfiles_test_binary"), ".go", ".txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run program: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	lines := strings.Split(outputStr, "\n")

	// Should have 3 files (.go and .txt files)
	expectedFiles := []string{"newest.go", "another.go", "middle.txt"}
	if len(lines) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d. Output: %s", len(expectedFiles), len(lines), outputStr)
	}

	// Check order (newest first)
	if lines[0] != "newest.go" {
		t.Errorf("Expected newest.go first, got %s", lines[0])
	}
}

func TestJSONOutput(t *testing.T) {
	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "newestfiles_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	createTestFiles(t, tmpDir)

	// Change to test directory
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Run the program with -j flag
	cmd := exec.Command(filepath.Join(oldDir, "newestfiles_test_binary"), "-j", ".go", ".txt")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run program: %v", err)
	}

	// Parse JSON output
	var files []string
	err = json.Unmarshal(output, &files)
	if err != nil {
		t.Fatalf("Failed to parse JSON output: %v. Output: %s", err, string(output))
	}

	// Should have 3 files (.go and .txt files)
	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d", expectedCount, len(files))
	}

	// Check order (newest first)
	if len(files) > 0 && files[0] != "newest.go" {
		t.Errorf("Expected newest.go first, got %s", files[0])
	}
}

func TestNoSuffixesProvided(t *testing.T) {
	oldDir, _ := os.Getwd()

	// Run the program without any arguments
	cmd := exec.Command(filepath.Join(oldDir, "newestfiles_test_binary"))
	output, err := cmd.Output()
	if err != nil {
		// This is expected to fail, but we want to check the output
		if exitError, ok := err.(*exec.ExitError); ok {
			output = exitError.Stderr
		}
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Usage:") {
		t.Errorf("Expected usage message, got: %s", outputStr)
	}
}

func TestNoFilesFound(t *testing.T) {
	// Create temporary directory with no matching files
	tmpDir, err := ioutil.TempDir("", "newestfiles_test_empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file that won't match
	testFile := filepath.Join(tmpDir, "test.xyz")
	err = ioutil.WriteFile(testFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to test directory
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Run the program looking for .go files
	cmd := exec.Command(filepath.Join(oldDir, "newestfiles_test_binary"), ".go")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run program: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if !strings.Contains(outputStr, "No files found") {
		t.Errorf("Expected 'No files found' message, got: %s", outputStr)
	}
}

func TestSuffixNormalization(t *testing.T) {
	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "newestfiles_test_suffix")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "test.go")
	err = ioutil.WriteFile(testFile, []byte("package main"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to test directory
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Test with suffix without dot
	cmd := exec.Command(filepath.Join(oldDir, "newestfiles_test_binary"), "go")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run program: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if !strings.Contains(outputStr, "test.go") {
		t.Errorf("Expected to find test.go, got: %s", outputStr)
	}
}

func TestCaseInsensitiveMatching(t *testing.T) {
	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "newestfiles_test_case")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files with different cases
	files := []string{"test.GO", "another.Go", "third.gO"}
	for _, filename := range files {
		testFile := filepath.Join(tmpDir, filename)
		err = ioutil.WriteFile(testFile, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Change to test directory
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	// Test with lowercase suffix
	cmd := exec.Command(filepath.Join(oldDir, "newestfiles_test_binary"), ".go")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run program: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	lines := strings.Split(outputStr, "\n")

	// Should find all 3 files regardless of case
	if len(lines) != 3 {
		t.Errorf("Expected 3 files, got %d. Output: %s", len(lines), outputStr)
	}
}
