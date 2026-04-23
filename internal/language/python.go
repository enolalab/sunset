package language

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

func init() {
	Register(&Language{
		Name:       "Python",
		ID:         "python",
		Extensions: []string{".py"},
		Grammar:    tree_sitter.NewLanguage(tree_sitter_python.Language()),
	})
}
