package parser

import (
	"strings"
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

// generateGoSource creates a Go source file with n functions.
func generateGoSource(n int) []byte {
	var b strings.Builder
	b.WriteString("package main\n\nimport \"fmt\"\n\n")
	for i := range n {
		b.WriteString("// FunctionDoc is the docstring.\n")
		b.WriteString("func function")
		b.WriteString(strings.Repeat("x", 1)) // keep name short
		b.WriteString(string(rune('A' + i%26)))
		b.WriteString("() {\n")
		b.WriteString("\tfmt.Println(\"hello\")\n")
		b.WriteString("\tx := 1 + 2\n")
		b.WriteString("\t_ = x\n")
		b.WriteString("}\n\n")
	}
	return []byte(b.String())
}

func setupBenchParser(b *testing.B) *Parser {
	b.Helper()
	p := NewParser()
	if err := p.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language())); err != nil {
		b.Fatalf("SetLanguage: %v", err)
	}
	return p
}

func BenchmarkWalk_SmallFile(b *testing.B) {
	p := setupBenchParser(b)
	defer p.Close()

	source := []byte("package main\n\nfunc main() {\n\tx := 1\n\tfmt.Println(x)\n}\n")
	tree, _ := p.Parse(source, "go")
	defer tree.Close()

	b.ResetTimer()
	for range b.N {
		Walk(tree.RootNode(), func(n *tree_sitter.Node, depth int) bool {
			return true
		})
	}
}

func BenchmarkWalk_LargeFile(b *testing.B) {
	p := setupBenchParser(b)
	defer p.Close()

	// ~150 functions ≈ ~1000 lines
	source := generateGoSource(150)
	tree, _ := p.Parse(source, "go")
	defer tree.Close()

	nodeCount := CountNodes(tree.RootNode())
	b.Logf("Source: %d bytes, %d nodes", len(source), nodeCount)

	b.ResetTimer()
	for range b.N {
		Walk(tree.RootNode(), func(n *tree_sitter.Node, depth int) bool {
			return true
		})
	}
}

func BenchmarkFilter_LargeFile(b *testing.B) {
	p := setupBenchParser(b)
	defer p.Close()

	source := generateGoSource(150)
	tree, _ := p.Parse(source, "go")
	defer tree.Close()

	b.ResetTimer()
	for range b.N {
		Filter(tree.RootNode(), "function_declaration")
	}
}

func BenchmarkCountNodes_LargeFile(b *testing.B) {
	p := setupBenchParser(b)
	defer p.Close()

	source := generateGoSource(150)
	tree, _ := p.Parse(source, "go")
	defer tree.Close()

	b.ResetTimer()
	for range b.N {
		CountNodes(tree.RootNode())
	}
}
