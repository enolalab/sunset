// Package parser provides a wrapper around go-tree-sitter for parsing source code
// into concrete syntax trees.
package parser

import (
	"fmt"
	"os"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// Tree wraps a tree-sitter parse tree with convenience methods.
type Tree struct {
	inner    *tree_sitter.Tree
	source   []byte
	Language string
}

// Parser wraps a tree-sitter parser instance.
type Parser struct {
	inner *tree_sitter.Parser
}

// NewParser creates a new Parser instance.
// The caller must call Close() when done to free C memory.
func NewParser() *Parser {
	return &Parser{
		inner: tree_sitter.NewParser(),
	}
}

// Close releases the underlying C resources.
// Safe to call multiple times.
func (p *Parser) Close() {
	if p.inner != nil {
		p.inner.Close()
		p.inner = nil
	}
}

// SetLanguage sets the tree-sitter language grammar for this parser.
func (p *Parser) SetLanguage(lang *tree_sitter.Language) error {
	if p.inner == nil {
		return fmt.Errorf("parser is closed")
	}
	return p.inner.SetLanguage(lang)
}

// Parse parses the given source code and returns a Tree.
// The languageName is stored in the returned Tree for reference.
func (p *Parser) Parse(source []byte, languageName string) (*Tree, error) {
	if p.inner == nil {
		return nil, fmt.Errorf("parser is closed")
	}

	tree := p.inner.Parse(source, nil)
	if tree == nil {
		return nil, fmt.Errorf("failed to parse source")
	}

	return &Tree{
		inner:    tree,
		source:   source,
		Language: languageName,
	}, nil
}

// ParseFile reads a file and parses it using the given language grammar.
// The caller must call Close() on the returned Tree when done.
func (p *Parser) ParseFile(path string, lang *tree_sitter.Language, languageName string) (*Tree, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	if err := p.SetLanguage(lang); err != nil {
		return nil, err
	}

	return p.Parse(content, languageName)
}

// RootNode returns the root node of the parse tree.
func (t *Tree) RootNode() *tree_sitter.Node {
	return t.inner.RootNode()
}

// Source returns the original source code.
func (t *Tree) Source() []byte {
	return t.source
}

// Close releases the underlying C resources of the tree.
// Safe to call multiple times.
func (t *Tree) Close() {
	if t.inner != nil {
		t.inner.Close()
		t.inner = nil
	}
}

// HasErrors returns true if the parse tree contains any ERROR or MISSING nodes.
func (t *Tree) HasErrors() bool {
	return hasErrorNodes(t.RootNode())
}

func hasErrorNodes(node *tree_sitter.Node) bool {
	if node.IsError() || node.IsMissing() {
		return true
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(uint(i))
		if child != nil && hasErrorNodes(child) {
			return true
		}
	}
	return false
}

// NodeText returns the source text for a given node.
func (t *Tree) NodeText(node *tree_sitter.Node) string {
	return node.Utf8Text(t.source)
}
