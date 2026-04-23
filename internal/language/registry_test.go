package language

import (
	"errors"
	"testing"
)

func TestDetect_Go(t *testing.T) {
	lang, err := Detect("main.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lang.ID != "go" {
		t.Errorf("expected 'go', got %q", lang.ID)
	}
	if lang.Name != "Go" {
		t.Errorf("expected 'Go', got %q", lang.Name)
	}
}

func TestDetect_JavaScript(t *testing.T) {
	cases := []struct {
		filename string
		id       string
	}{
		{"app.js", "javascript"},
		{"component.jsx", "javascript"},
	}
	for _, tc := range cases {
		t.Run(tc.filename, func(t *testing.T) {
			lang, err := Detect(tc.filename)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if lang.ID != tc.id {
				t.Errorf("expected %q, got %q", tc.id, lang.ID)
			}
		})
	}
}

func TestDetect_TypeScript(t *testing.T) {
	cases := []struct {
		filename string
		id       string
	}{
		{"app.ts", "typescript"},
		{"component.tsx", "typescript"},
	}
	for _, tc := range cases {
		t.Run(tc.filename, func(t *testing.T) {
			lang, err := Detect(tc.filename)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if lang.ID != tc.id {
				t.Errorf("expected %q, got %q", tc.id, lang.ID)
			}
		})
	}
}

func TestDetect_Python(t *testing.T) {
	lang, err := Detect("script.py")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lang.ID != "python" {
		t.Errorf("expected 'python', got %q", lang.ID)
	}
}

func TestDetect_UnsupportedExtension(t *testing.T) {
	_, err := Detect("unknown.xyz")
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}
	if !errors.Is(err, ErrUnsupportedLanguage) {
		t.Errorf("expected ErrUnsupportedLanguage, got: %v", err)
	}
}

func TestDetect_NoExtension(t *testing.T) {
	_, err := Detect("Makefile")
	if err == nil {
		t.Fatal("expected error for file without extension")
	}
	if !errors.Is(err, ErrUnsupportedLanguage) {
		t.Errorf("expected ErrUnsupportedLanguage, got: %v", err)
	}
}

func TestDetect_CaseInsensitive(t *testing.T) {
	lang, err := Detect("MAIN.GO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lang.ID != "go" {
		t.Errorf("expected 'go', got %q", lang.ID)
	}
}

func TestGet_ValidID(t *testing.T) {
	lang, err := Get("go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lang.Name != "Go" {
		t.Errorf("expected 'Go', got %q", lang.Name)
	}
}

func TestGet_InvalidID(t *testing.T) {
	_, err := Get("cobol")
	if err == nil {
		t.Fatal("expected error for unknown language ID")
	}
	if !errors.Is(err, ErrUnsupportedLanguage) {
		t.Errorf("expected ErrUnsupportedLanguage, got: %v", err)
	}
}

func TestAll_ReturnsRegistered(t *testing.T) {
	all := All()
	if len(all) < 3 {
		t.Errorf("expected at least 3 languages (go, js, ts, python), got %d", len(all))
	}
}

func TestIsSupported(t *testing.T) {
	if !IsSupported("main.go") {
		t.Error("expected main.go to be supported")
	}
	if IsSupported("main.rb") {
		t.Error("expected main.rb to not be supported")
	}
}

func TestSupportedExtensions(t *testing.T) {
	exts := SupportedExtensions()
	if len(exts) < 5 {
		t.Errorf("expected at least 5 extensions (.go, .js, .jsx, .ts, .tsx, .py), got %d", len(exts))
	}

	// Check .go is in the list
	found := false
	for _, ext := range exts {
		if ext == ".go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected .go in supported extensions")
	}
}

func TestGrammar_NotNil(t *testing.T) {
	for _, lang := range All() {
		if lang.Grammar == nil {
			t.Errorf("grammar for %s is nil", lang.Name)
		}
	}
}
