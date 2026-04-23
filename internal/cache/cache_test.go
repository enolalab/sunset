package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCache_LoadEmpty(t *testing.T) {
	tmp := t.TempDir()
	c, err := Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(c.Entries) != 0 {
		t.Errorf("expected empty cache, got %d entries", len(c.Entries))
	}
}

func TestCache_SaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	c, err := Load(tmp)
	if err != nil {
		t.Fatal(err)
	}

	c.Update("main.go", []byte("package main"), "go")
	c.Update("utils.go", []byte("package utils"), "go")

	if err := c.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Verify file exists
	cachePath := filepath.Join(tmp, CacheDir, CacheSubDir, CacheFile)
	if _, err := os.Stat(cachePath); err != nil {
		t.Fatalf("cache file not created: %v", err)
	}

	// Reload
	c2, err := Load(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(c2.Entries) != 2 {
		t.Errorf("expected 2 entries after reload, got %d", len(c2.Entries))
	}
}

func TestCache_IsChanged(t *testing.T) {
	tmp := t.TempDir()
	c, _ := Load(tmp)

	content := []byte("package main")

	// New file → changed
	if !c.IsChanged("main.go", content) {
		t.Error("new file should be marked as changed")
	}

	// Update cache
	c.Update("main.go", content, "go")

	// Same content → not changed
	if c.IsChanged("main.go", content) {
		t.Error("same content should NOT be changed")
	}

	// Different content → changed
	if !c.IsChanged("main.go", []byte("package main\nfunc main() {}")) {
		t.Error("modified content should be changed")
	}
}

func TestCache_Remove(t *testing.T) {
	tmp := t.TempDir()
	c, _ := Load(tmp)

	c.Update("main.go", []byte("pkg"), "go")
	c.Remove("main.go")

	if _, exists := c.Entries["main.go"]; exists {
		t.Error("entry should be removed")
	}
}

func TestCache_Prune(t *testing.T) {
	tmp := t.TempDir()
	c, _ := Load(tmp)

	c.Update("main.go", []byte("pkg"), "go")
	c.Update("deleted.go", []byte("pkg"), "go")
	c.Update("utils.go", []byte("pkg"), "go")

	// Only main.go and utils.go exist now
	removed := c.Prune([]string{"main.go", "utils.go"})

	if len(removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(removed))
	}
	if len(removed) > 0 && removed[0] != "deleted.go" {
		t.Errorf("expected 'deleted.go' removed, got %q", removed[0])
	}
	if _, exists := c.Entries["deleted.go"]; exists {
		t.Error("deleted.go should be pruned from cache")
	}
}

func TestCache_Clean(t *testing.T) {
	tmp := t.TempDir()
	c, _ := Load(tmp)
	c.Update("main.go", []byte("pkg"), "go")
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}

	if err := Clean(tmp); err != nil {
		t.Fatalf("Clean: %v", err)
	}

	// Directory should be gone
	if _, err := os.Stat(filepath.Join(tmp, CacheDir)); !os.IsNotExist(err) {
		t.Error("cache dir should be removed after Clean")
	}
}

func TestCache_CorruptFile(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, CacheDir, CacheSubDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	// Write invalid JSON
	if err := os.WriteFile(filepath.Join(dir, CacheFile), []byte("{invalid"), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := Load(tmp)
	if err != nil {
		t.Fatalf("Load should not fail on corrupt cache: %v", err)
	}
	if len(c.Entries) != 0 {
		t.Error("corrupt cache should result in empty entries")
	}
}

func TestHashContent_Deterministic(t *testing.T) {
	content := []byte("hello world")
	h1 := HashContent(content)
	h2 := HashContent(content)
	if h1 != h2 {
		t.Error("hash should be deterministic")
	}
	if len(h1) != 64 { // SHA256 hex = 64 chars
		t.Errorf("expected 64 char hash, got %d", len(h1))
	}
}

func TestHashContent_DifferentInputs(t *testing.T) {
	h1 := HashContent([]byte("hello"))
	h2 := HashContent([]byte("world"))
	if h1 == h2 {
		t.Error("different inputs should produce different hashes")
	}
}
