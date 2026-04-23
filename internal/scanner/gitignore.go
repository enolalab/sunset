package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GitignoreRules holds parsed .gitignore patterns from the project root.
type GitignoreRules struct {
	patterns []gitignorePattern
}

type gitignorePattern struct {
	pattern  string
	negated  bool
	dirOnly  bool
	anchored bool // pattern contains /
}

// LoadGitignore reads and parses .gitignore from the given root directory.
// Returns nil if no .gitignore is found.
func LoadGitignore(root string) *GitignoreRules {
	path := filepath.Join(root, ".gitignore")
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	rules := &GitignoreRules{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || line[0] == '#' {
			continue
		}
		rules.patterns = append(rules.patterns, parsePattern(line))
	}

	return rules
}

func parsePattern(raw string) gitignorePattern {
	p := gitignorePattern{}

	// Check negation
	if raw[0] == '!' {
		p.negated = true
		raw = raw[1:]
	}

	// Check dir-only
	if strings.HasSuffix(raw, "/") {
		p.dirOnly = true
		raw = strings.TrimSuffix(raw, "/")
	}

	// Check if anchored (contains / not at end)
	p.anchored = strings.Contains(raw, "/")

	// Remove leading /
	raw = strings.TrimPrefix(raw, "/")

	p.pattern = raw
	return p
}

// IsIgnored returns true if the given relative path should be ignored.
func (g *GitignoreRules) IsIgnored(relPath string) bool {
	if g == nil {
		return false
	}

	// Normalize to forward slashes for pattern matching
	relPath = filepath.ToSlash(relPath)
	ignored := false

	for _, p := range g.patterns {
		if matchesPattern(relPath, p) {
			ignored = !p.negated
		}
	}

	return ignored
}

// matchesPattern checks if a path matches a single gitignore pattern.
func matchesPattern(relPath string, p gitignorePattern) bool {
	pattern := p.pattern

	if p.anchored {
		// Match from root
		matched, _ := filepath.Match(pattern, relPath)
		if matched {
			return true
		}
		// Also try as prefix for directory patterns
		if p.dirOnly && strings.HasPrefix(relPath, pattern+"/") {
			return true
		}
		return false
	}

	// Non-anchored: match against any path component or the basename
	name := filepath.Base(relPath)

	// Match against basename
	if matched, _ := filepath.Match(pattern, name); matched {
		return true
	}

	// For dir-only patterns, check if any directory component matches
	if p.dirOnly {
		parts := strings.Split(relPath, "/")
		for _, part := range parts[:len(parts)-1] { // exclude the filename itself
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}
		}
		// Also check if the path starts with the pattern as a dir prefix
		if strings.HasPrefix(relPath, pattern+"/") {
			return true
		}
	}

	// Match against each suffix of the path
	parts := strings.Split(relPath, "/")
	for i := range parts {
		subpath := strings.Join(parts[i:], "/")
		if matched, _ := filepath.Match(pattern, subpath); matched {
			return true
		}
	}

	return false
}
