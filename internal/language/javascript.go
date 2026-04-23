package language

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

func init() {
	Register(&Language{
		Name:       "JavaScript",
		ID:         "javascript",
		Extensions: []string{".js", ".jsx"},
		Grammar:    tree_sitter.NewLanguage(tree_sitter_javascript.Language()),
	})

	Register(&Language{
		Name:       "TypeScript",
		ID:         "typescript",
		Extensions: []string{".ts", ".tsx"},
		Grammar:    tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript()),
	})
}
