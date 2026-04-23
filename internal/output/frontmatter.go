package output

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// RenderFileFrontmatter generates YAML frontmatter for a single file.
func RenderFileFrontmatter(info *FileInfo) (string, error) {
	// Build compact frontmatter struct (no detail arrays)
	fm := struct {
		File          string   `yaml:"file"`
		Language      string   `yaml:"language"`
		Package       string   `yaml:"package,omitempty"`
		Lines         int      `yaml:"lines"`
		FunctionCount int      `yaml:"function_count"`
		TypeCount     int      `yaml:"type_count"`
		ImportCount   int      `yaml:"import_count"`
		Tags          []string `yaml:"tags,omitempty"`
	}{
		File:          info.File,
		Language:      info.Language,
		Package:       info.Package,
		Lines:         info.Lines,
		FunctionCount: info.FunctionCount,
		TypeCount:     info.TypeCount,
		ImportCount:   info.ImportCount,
		Tags:          info.Tags,
	}

	// Limit tags to 10
	if len(fm.Tags) > 10 {
		fm.Tags = fm.Tags[:10]
	}

	return marshalFrontmatter(fm)
}

// RenderProjectFrontmatter generates YAML frontmatter for index.md.
func RenderProjectFrontmatter(info *ProjectInfo) (string, error) {
	return marshalFrontmatter(info)
}

func marshalFrontmatter(v interface{}) (string, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return "", fmt.Errorf("marshaling frontmatter: %w", err)
	}
	if err := enc.Close(); err != nil {
		return "", fmt.Errorf("closing encoder: %w", err)
	}

	return "---\n" + buf.String() + "---\n", nil
}
