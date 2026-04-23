// Package cache provides file-level hash caching for incremental updates.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// CacheDir is the directory name for cache storage.
	CacheDir = ".sunset"
	// CacheSubDir is the cache subdirectory within .sunset.
	CacheSubDir = "cache"
	// CacheFile is the filename within the cache directory.
	CacheFile = "cache.json"
)

// Entry stores hash and metadata for a single file.
type Entry struct {
	Hash      string    `json:"hash"`
	Language  string    `json:"language,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Cache manages file hash entries for incremental builds.
type Cache struct {
	Entries map[string]Entry `json:"entries"`
	dir     string
}

// Load reads the cache from disk. Returns an empty cache if file doesn't exist.
func Load(rootDir string) (*Cache, error) {
	c := &Cache{
		Entries: make(map[string]Entry),
		dir:     filepath.Join(rootDir, CacheDir, CacheSubDir),
	}

	data, err := os.ReadFile(filepath.Join(c.dir, CacheFile))
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil // fresh cache
		}
		return nil, fmt.Errorf("reading cache: %w", err)
	}

	if err := json.Unmarshal(data, &c.Entries); err != nil {
		// Corrupt cache — start fresh
		c.Entries = make(map[string]Entry)
		return c, nil
	}

	return c, nil
}

// Save writes the cache to disk.
func (c *Cache) Save() error {
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return fmt.Errorf("creating cache dir: %w", err)
	}

	data, err := json.MarshalIndent(c.Entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling cache: %w", err)
	}

	return os.WriteFile(filepath.Join(c.dir, CacheFile), data, 0644)
}

// IsChanged returns true if the file content differs from the cached hash.
func (c *Cache) IsChanged(relPath string, content []byte) bool {
	hash := HashContent(content)
	entry, exists := c.Entries[relPath]
	if !exists {
		return true
	}
	return entry.Hash != hash
}

// Update stores the hash for a file.
func (c *Cache) Update(relPath string, content []byte, langID string) {
	c.Entries[relPath] = Entry{
		Hash:      HashContent(content),
		Language:  langID,
		UpdatedAt: time.Now(),
	}
}

// Remove deletes a file entry from the cache.
func (c *Cache) Remove(relPath string) {
	delete(c.Entries, relPath)
}

// Prune removes cache entries for files that no longer exist in the file list.
// Returns the list of removed paths.
func (c *Cache) Prune(currentFiles []string) []string {
	current := make(map[string]bool, len(currentFiles))
	for _, f := range currentFiles {
		current[f] = true
	}

	var removed []string
	for path := range c.Entries {
		if !current[path] {
			removed = append(removed, path)
			delete(c.Entries, path)
		}
	}
	return removed
}

// Clean removes the entire cache directory.
func Clean(rootDir string) error {
	return os.RemoveAll(filepath.Join(rootDir, CacheDir))
}

// HashContent computes SHA256 hex digest of content.
func HashContent(content []byte) string {
	h := sha256.Sum256(content)
	return hex.EncodeToString(h[:])
}
