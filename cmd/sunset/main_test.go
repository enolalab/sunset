package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	binary := filepath.Join(tmp, "sunset")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/sunset/")
	cmd.Dir = findProjectRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return binary
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	// Walk up to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}
}

func runSunset(t *testing.T, binary string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = findProjectRoot(t)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return string(out), exitCode
}

func TestCLI_NoArgs_ShowsHelp(t *testing.T) {
	bin := buildBinary(t)
	out, code := runSunset(t, bin)

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(out, "Sunset") {
		t.Error("output should mention Sunset")
	}
	if !strings.Contains(out, "parse") {
		t.Error("output should list parse command")
	}
}

func TestCLI_Version(t *testing.T) {
	bin := buildBinary(t)
	out, code := runSunset(t, bin, "version")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(out, "sunset") {
		t.Error("version output should contain 'sunset'")
	}
}

func TestCLI_Languages(t *testing.T) {
	bin := buildBinary(t)
	out, code := runSunset(t, bin, "languages")

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	for _, lang := range []string{"go", "javascript", "typescript", "python"} {
		if !strings.Contains(out, lang) {
			t.Errorf("output should contain %s", lang)
		}
	}
}

func TestCLI_ParseHelp(t *testing.T) {
	bin := buildBinary(t)
	out, _ := runSunset(t, bin, "parse", "--help")

	flags := []string{"--output", "--detail", "--exclude", "--concurrency", "--no-cache", "--quiet"}
	for _, flag := range flags {
		if !strings.Contains(out, strings.TrimPrefix(flag, "--")) {
			t.Errorf("parse --help should mention %s", flag)
		}
	}
}

func TestCLI_ParseTestdata(t *testing.T) {
	bin := buildBinary(t)
	tmp := t.TempDir()
	outDir := filepath.Join(t.TempDir(), "out")

	// Create source files
	if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "utils.go"), []byte("package main\n\nfunc helper() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	out, code := runSunset(t, bin, "parse", "--output", outDir, "--no-cache", tmp)
	t.Logf("parse output: %s", out)
	if code != 0 {
		t.Fatalf("parse failed (code %d): %s", code, out)
	}

	// Check index.md
	if _, err := os.Stat(filepath.Join(outDir, "index.md")); err != nil {
		t.Error("index.md should be created")
	}

	// Check files directory
	entries, err := os.ReadDir(filepath.Join(outDir, "files"))
	if err != nil {
		t.Fatalf("reading files dir: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 output files, got %d", len(entries))
	}
}

func TestCLI_InvalidPath(t *testing.T) {
	bin := buildBinary(t)
	out, code := runSunset(t, bin, "parse", "/nonexistent/path")

	if code == 0 {
		t.Error("expected non-zero exit code for invalid path")
	}
	if !strings.Contains(out, "Error") {
		t.Errorf("expected error message, got: %s", out)
	}
}

func TestCLI_UnknownCommand(t *testing.T) {
	bin := buildBinary(t)
	_, code := runSunset(t, bin, "foobar")

	if code == 0 {
		t.Error("expected non-zero exit code for unknown command")
	}
}

func TestCLI_Clean(t *testing.T) {
	bin := buildBinary(t)
	tmp := t.TempDir()

	// Create .sunset dir
	if err := os.MkdirAll(filepath.Join(tmp, ".sunset", "cache"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, ".sunset", "output"), 0755); err != nil {
		t.Fatal(err)
	}

	out, code := runSunset(t, bin, "clean", tmp)
	if code != 0 {
		t.Errorf("clean failed (code %d): %s", code, out)
	}

	if _, err := os.Stat(filepath.Join(tmp, ".sunset")); !os.IsNotExist(err) {
		t.Error(".sunset should be cleaned")
	}
}
