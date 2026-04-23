package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitignore_MatchPattern(t *testing.T) {
	tmp := t.TempDir()

	// Create .gitignore
	if err := os.WriteFile(filepath.Join(tmp, ".gitignore"), []byte("*.log\nbuild/\n"), 0644); err != nil {
		t.Fatal(err)
	}

	rules := LoadGitignore(tmp)
	if rules == nil {
		t.Fatal("expected rules to be loaded")
	}

	if !rules.IsIgnored("app.log") {
		t.Error("expected app.log to be ignored")
	}
	if !rules.IsIgnored("debug.log") {
		t.Error("expected debug.log to be ignored")
	}
	if rules.IsIgnored("main.go") {
		t.Error("main.go should not be ignored")
	}
}

func TestGitignore_NoFile(t *testing.T) {
	rules := LoadGitignore(t.TempDir())
	if rules != nil {
		t.Error("expected nil rules when no .gitignore")
	}
}

func TestGitignore_NilRules(t *testing.T) {
	var rules *GitignoreRules
	if rules.IsIgnored("anything") {
		t.Error("nil rules should not ignore anything")
	}
}

func TestGitignore_DirPattern(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, ".gitignore"), []byte("build/\n"), 0644); err != nil {
		t.Fatal(err)
	}

	rules := LoadGitignore(tmp)
	if !rules.IsIgnored("build/output.go") {
		t.Error("expected build/output.go to be ignored by build/ pattern")
	}
}

func TestGitignore_Negation(t *testing.T) {
	tmp := t.TempDir()
	content := "*.log\n!important.log\n"
	if err := os.WriteFile(filepath.Join(tmp, ".gitignore"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rules := LoadGitignore(tmp)
	if !rules.IsIgnored("debug.log") {
		t.Error("expected debug.log to be ignored")
	}
	if rules.IsIgnored("important.log") {
		t.Error("important.log should NOT be ignored (negated)")
	}
}

func TestGitignore_Comments(t *testing.T) {
	tmp := t.TempDir()
	content := "# This is a comment\n*.tmp\n"
	if err := os.WriteFile(filepath.Join(tmp, ".gitignore"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rules := LoadGitignore(tmp)
	if !rules.IsIgnored("cache.tmp") {
		t.Error("expected cache.tmp to be ignored")
	}
}

func TestGitignore_IntegrationWithScan(t *testing.T) {
	tmp := t.TempDir()

	// Create files
	if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "generated.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, ".gitignore"), []byte("generated.go\n"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Scan(tmp, nil)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if len(result.Files) != 1 {
		t.Errorf("expected 1 file (generated.go excluded), got %d: %v", len(result.Files), result.Files)
	}
}
