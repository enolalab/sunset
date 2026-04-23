package language

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

func init() {
	Register(&Language{
		Name:       "Go",
		ID:         "go",
		Extensions: []string{".go"},
		Grammar:    tree_sitter.NewLanguage(tree_sitter_go.Language()),
	})
}
