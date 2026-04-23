package sunset

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_Go(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	if result.Language != "Go" {
		t.Errorf("expected language 'Go', got %q", result.Language)
	}
	if result.LanguageID != "go" {
		t.Errorf("expected language ID 'go', got %q", result.LanguageID)
	}
	if result.Path != path {
		t.Errorf("expected path %q, got %q", path, result.Path)
	}
	if len(result.Source) == 0 {
		t.Error("expected non-empty source")
	}
}

func TestParseFile_TypeScript(t *testing.T) {
	path := filepath.Join("testdata", "js-sample", "src", "index.ts")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	if result.Language != "TypeScript" {
		t.Errorf("expected language 'TypeScript', got %q", result.Language)
	}
}

func TestParseFile_Python(t *testing.T) {
	path := filepath.Join("testdata", "python-sample", "main.py")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	if result.Language != "Python" {
		t.Errorf("expected language 'Python', got %q", result.Language)
	}
}

func TestParseFile_NonexistentFile(t *testing.T) {
	_, err := ParseFile("nonexistent.go")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestParseFile_UnsupportedLanguage(t *testing.T) {
	// Create a temporary file with unsupported extension
	tmp := filepath.Join(t.TempDir(), "test.xyz")
	if err := os.WriteFile(tmp, []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := ParseFile(tmp)
	if err == nil {
		t.Error("expected error for unsupported extension")
	}
}

func TestParseFile_WithLanguageOption(t *testing.T) {
	// Force Go parser on the Go sample file
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path, WithLanguage("go"))
	if err != nil {
		t.Fatalf("ParseFile with WithLanguage failed: %v", err)
	}
	defer result.Close()

	if result.LanguageID != "go" {
		t.Errorf("expected language ID 'go', got %q", result.LanguageID)
	}
}

func TestParseFile_TreeNotNil(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	if result.Tree == nil {
		t.Fatal("expected non-nil Tree")
	}
}

func TestNode_Type(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	root := result.Tree.RootNode()
	if root.Type() != "source_file" {
		t.Errorf("expected root type 'source_file', got %q", root.Type())
	}
}

func TestNode_Children(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	root := result.Tree.RootNode()
	children := root.Children()
	if len(children) == 0 {
		t.Error("expected root to have children")
	}
}

func TestNode_Raw(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	root := result.Tree.RootNode()
	raw := root.Raw()
	if raw == nil {
		t.Error("expected non-nil raw tree-sitter node")
	}
}

func TestNode_Position(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	root := result.Tree.RootNode()
	pos := root.Position()
	if pos.Row != 0 || pos.Column != 0 {
		t.Errorf("expected root position (0,0), got (%d,%d)", pos.Row, pos.Column)
	}
}

func TestTreeWrapper_Walk(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	count := 0
	result.Tree.Walk(func(node *Node, depth int) bool {
		count++
		return true
	})

	if count == 0 {
		t.Error("expected Walk to visit at least one node")
	}
}

func TestTreeWrapper_Filter(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	funcs := result.Tree.Filter("function_declaration")
	if len(funcs) != 2 {
		t.Errorf("expected 2 function_declarations in main.go, got %d", len(funcs))
	}
}

func TestFileResult_HasErrors_ValidCode(t *testing.T) {
	path := filepath.Join("testdata", "go-sample", "main.go")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	defer result.Close()

	if result.HasErrors() {
		t.Error("expected no errors for valid Go code")
	}
}

func TestLanguages(t *testing.T) {
	langs := Languages()
	if len(langs) < 3 {
		t.Errorf("expected at least 3 languages, got %d", len(langs))
	}

	// Check Go is present
	found := false
	for _, l := range langs {
		if l.ID == "go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Go in languages list")
	}
}
