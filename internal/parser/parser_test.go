package parser

import (
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

func goLang() *tree_sitter.Language {
	return tree_sitter.NewLanguage(tree_sitter_go.Language())
}

func pythonLang() *tree_sitter.Language {
	return tree_sitter.NewLanguage(tree_sitter_python.Language())
}

func TestParser_ParseGoCode(t *testing.T) {
	p := NewParser()
	defer p.Close()

	if err := p.SetLanguage(goLang()); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}

	source := []byte("package main\n\nfunc main() {}\n")
	tree, err := p.Parse(source, "go")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer tree.Close()

	root := tree.RootNode()
	if root == nil {
		t.Fatal("expected non-nil root node")
	}
	if root.Kind() != "source_file" {
		t.Errorf("expected root type 'source_file', got %q", root.Kind())
	}
}

func TestParser_ParsePythonCode(t *testing.T) {
	p := NewParser()
	defer p.Close()

	if err := p.SetLanguage(pythonLang()); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}

	source := []byte("def hello():\n    print('hi')\n")
	tree, err := p.Parse(source, "python")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	defer tree.Close()

	root := tree.RootNode()
	if root.Kind() != "module" {
		t.Errorf("expected root type 'module', got %q", root.Kind())
	}
}

func TestParser_ParseInvalidContent(t *testing.T) {
	p := NewParser()
	defer p.Close()

	if err := p.SetLanguage(goLang()); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}

	source := []byte("this is not valid go @@@ code {{{")
	tree, err := p.Parse(source, "go")
	if err != nil {
		t.Fatalf("expected tree-sitter to return tree with errors, not parse error: %v", err)
	}
	defer tree.Close()

	if !tree.HasErrors() {
		t.Error("expected tree to have errors for invalid input")
	}
}

func TestParser_ParseEmptyFile(t *testing.T) {
	p := NewParser()
	defer p.Close()

	if err := p.SetLanguage(goLang()); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}

	tree, err := p.Parse([]byte(""), "go")
	if err != nil {
		t.Fatalf("expected no error for empty file, got: %v", err)
	}
	defer tree.Close()

	root := tree.RootNode()
	if root == nil {
		t.Fatal("expected non-nil root node for empty file")
	}
}

func TestParser_CloseMultipleTimes(t *testing.T) {
	p := NewParser()
	p.Close()
	p.Close() // should not panic
}

func TestTree_CloseMultipleTimes(t *testing.T) {
	p := NewParser()
	defer p.Close()
	if err := p.SetLanguage(goLang()); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}

	tree, _ := p.Parse([]byte("package main"), "go")
	tree.Close()
	tree.Close() // should not panic
}

func TestParser_ParseAfterClose(t *testing.T) {
	p := NewParser()
	p.Close()

	_, err := p.Parse([]byte("package main"), "go")
	if err == nil {
		t.Error("expected error when parsing after Close()")
	}
}

func TestTree_NodeText(t *testing.T) {
	p := NewParser()
	defer p.Close()
	if err := p.SetLanguage(goLang()); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}

	source := []byte("package main\n\nfunc hello() {}\n")
	tree, _ := p.Parse(source, "go")
	defer tree.Close()

	root := tree.RootNode()
	text := tree.NodeText(root)
	if text != string(source) {
		t.Errorf("expected root text to equal source, got %q", text)
	}
}

func TestTree_HasErrors_ValidCode(t *testing.T) {
	p := NewParser()
	defer p.Close()
	if err := p.SetLanguage(goLang()); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}

	tree, _ := p.Parse([]byte("package main\n\nfunc main() {}\n"), "go")
	defer tree.Close()

	if tree.HasErrors() {
		t.Error("expected no errors for valid Go code")
	}
}
