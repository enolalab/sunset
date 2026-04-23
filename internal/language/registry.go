// Package language provides a registry for mapping file extensions to
// tree-sitter language grammars.
package language

import (
	"fmt"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// ErrUnsupportedLanguage is returned when a file's language cannot be detected.
var ErrUnsupportedLanguage = fmt.Errorf("unsupported language")

// Language represents a supported programming language with its tree-sitter grammar.
type Language struct {
	// Name is the human-readable language name (e.g., "Go", "Python").
	Name string

	// ID is a short lowercase identifier (e.g., "go", "python").
	ID string

	// Extensions lists the file extensions associated with this language (e.g., ".go").
	Extensions []string

	// Grammar returns the tree-sitter Language for parsing.
	Grammar *tree_sitter.Language
}

// registry holds all registered languages.
var registry []*Language

// Register adds a language to the global registry.
func Register(lang *Language) {
	registry = append(registry, lang)
}

// Detect identifies the language of a file based on its extension.
// Returns ErrUnsupportedLanguage if no match is found.
func Detect(filename string) (*Language, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return nil, fmt.Errorf("%w: no extension for %q", ErrUnsupportedLanguage, filename)
	}

	for _, lang := range registry {
		for _, langExt := range lang.Extensions {
			if langExt == ext {
				return lang, nil
			}
		}
	}

	return nil, fmt.Errorf("%w: unknown extension %q", ErrUnsupportedLanguage, ext)
}

// Get returns a language by its ID (e.g., "go", "python").
// Returns ErrUnsupportedLanguage if not found.
func Get(id string) (*Language, error) {
	id = strings.ToLower(id)
	for _, lang := range registry {
		if lang.ID == id {
			return lang, nil
		}
	}
	return nil, fmt.Errorf("%w: unknown language %q", ErrUnsupportedLanguage, id)
}

// All returns all registered languages.
func All() []*Language {
	result := make([]*Language, len(registry))
	copy(result, registry)
	return result
}

// SupportedExtensions returns all supported file extensions.
func SupportedExtensions() []string {
	var exts []string
	for _, lang := range registry {
		exts = append(exts, lang.Extensions...)
	}
	return exts
}

// IsSupported returns true if the file extension is supported.
func IsSupported(filename string) bool {
	_, err := Detect(filename)
	return err == nil
}

// Supported returns all registered languages (alias for All).
func Supported() []*Language {
	return All()
}
