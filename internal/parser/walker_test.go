package parser

import (
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

func parseGoSource(t *testing.T, source string) *Tree {
	t.Helper()
	p := NewParser()
	t.Cleanup(p.Close)
	if err := p.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language())); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}
	tree, err := p.Parse([]byte(source), "go")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	t.Cleanup(tree.Close)
	return tree
}

func TestWalk_VisitsAllNodes(t *testing.T) {
	tree := parseGoSource(t, "package main\n\nfunc main() {}\n")

	count := 0
	Walk(tree.RootNode(), func(n *tree_sitter.Node, depth int) bool {
		count++
		return true
	})

	if count == 0 {
		t.Error("expected Walk to visit at least one node")
	}

	// Verify count matches CountNodes
	expected := CountNodes(tree.RootNode())
	if count != expected {
		t.Errorf("Walk visited %d nodes, CountNodes returned %d", count, expected)
	}
}

func TestWalk_SkipChildren(t *testing.T) {
	tree := parseGoSource(t, "package main\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n")

	visited := 0
	Walk(tree.RootNode(), func(n *tree_sitter.Node, depth int) bool {
		visited++
		// Skip children after depth 1
		return depth < 1
	})

	total := CountNodes(tree.RootNode())
	if visited >= total {
		t.Errorf("expected fewer nodes when skipping children, visited=%d total=%d", visited, total)
	}
}

func TestWalk_DepthCorrect(t *testing.T) {
	tree := parseGoSource(t, "package main\n\nfunc main() {}\n")

	Walk(tree.RootNode(), func(n *tree_sitter.Node, depth int) bool {
		if n == tree.RootNode() && depth != 0 {
			t.Errorf("root node should have depth 0, got %d", depth)
		}
		return true
	})
}

func TestFilter_FunctionDeclaration(t *testing.T) {
	tree := parseGoSource(t, "package main\n\nfunc foo() {}\n\nfunc bar() {}\n")

	funcs := Filter(tree.RootNode(), "function_declaration")
	if len(funcs) != 2 {
		t.Errorf("expected 2 function_declarations, got %d", len(funcs))
	}

	for _, f := range funcs {
		if f.Kind() != "function_declaration" {
			t.Errorf("expected kind 'function_declaration', got %q", f.Kind())
		}
	}
}

func TestFilter_NoMatch(t *testing.T) {
	tree := parseGoSource(t, "package main\n")

	results := Filter(tree.RootNode(), "class_declaration")
	if len(results) != 0 {
		t.Errorf("expected 0 matches for class_declaration in Go, got %d", len(results))
	}
}

func TestDepth_RootIsZero(t *testing.T) {
	tree := parseGoSource(t, "package main\n")

	d := Depth(tree.RootNode())
	if d != 0 {
		t.Errorf("expected root depth 0, got %d", d)
	}
}

func TestDepth_ChildIsOne(t *testing.T) {
	tree := parseGoSource(t, "package main\n")

	root := tree.RootNode()
	if root.ChildCount() == 0 {
		t.Skip("no children to test")
	}

	child := root.Child(0)
	d := Depth(child)
	if d != 1 {
		t.Errorf("expected child depth 1, got %d", d)
	}
}

func TestCountNodes(t *testing.T) {
	tree := parseGoSource(t, "package main\n\nfunc main() {}\n")

	count := CountNodes(tree.RootNode())
	if count < 3 {
		t.Errorf("expected at least 3 nodes (source_file, package_clause, func_decl), got %d", count)
	}
}

func TestMaxDepth(t *testing.T) {
	tree := parseGoSource(t, "package main\n\nfunc main() {\n\tx := 1\n}\n")

	maxD := MaxDepth(tree.RootNode())
	if maxD < 2 {
		t.Errorf("expected max depth >= 2 for nested code, got %d", maxD)
	}
}
