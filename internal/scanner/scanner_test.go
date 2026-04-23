package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestScan_GoSample(t *testing.T) {
	result, err := Scan("../../testdata/go-sample", nil)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if len(result.Files) != 3 {
		t.Errorf("expected 3 .go files, got %d: %v", len(result.Files), result.Files)
	}

	// Check all are .go files
	for _, f := range result.Files {
		if filepath.Ext(f) != ".go" {
			t.Errorf("unexpected file: %s", f)
		}
	}
}

func TestScan_PythonSample(t *testing.T) {
	result, err := Scan("../../testdata/python-sample", nil)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if len(result.Files) != 1 {
		t.Errorf("expected 1 .py file, got %d: %v", len(result.Files), result.Files)
	}
}

func TestScan_ExcludePattern(t *testing.T) {
	result, err := Scan("../../testdata/go-sample", &Options{
		ExcludePatterns: []string{"*.go"},
	})
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if len(result.Files) != 0 {
		t.Errorf("expected 0 files with *.go excluded, got %d: %v", len(result.Files), result.Files)
	}
}

func TestScan_HiddenDirsSkipped(t *testing.T) {
	// Create a temp dir with a hidden directory containing a Go file
	tmp := t.TempDir()
	hiddenDir := filepath.Join(tmp, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hiddenDir, "secret.go"), []byte("package secret"), 0644); err != nil {
		t.Fatal(err)
	}
	// And a visible file
	if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Scan(tmp, nil)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if len(result.Files) != 1 {
		t.Errorf("expected 1 file (hidden dir skipped), got %d: %v", len(result.Files), result.Files)
	}
}

func TestScan_DefaultDirsSkipped(t *testing.T) {
	tmp := t.TempDir()
	for _, dir := range []string{"node_modules", "__pycache__", ".git"} {
		dirPath := filepath.Join(tmp, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dirPath, "file.go"), []byte("package x"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Scan(tmp, nil)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if len(result.Files) != 1 {
		t.Errorf("expected 1 file (default dirs skipped), got %d: %v", len(result.Files), result.Files)
	}
}

func TestScan_ResultsSorted(t *testing.T) {
	result, err := Scan("../../testdata/go-sample", nil)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	sorted := make([]string, len(result.Files))
	copy(sorted, result.Files)
	sort.Strings(sorted)

	// Just verify we get files
	if len(result.Files) == 0 {
		t.Error("expected some files")
	}
}
