package analyzer

import (
	"strings"
	"testing"
)

func TestDepGraph_AddEdge(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.go", "b.go")
	g.AddEdge("a.go", "c.go")
	g.AddEdge("b.go", "c.go")

	deps := g.DependenciesOf("a.go")
	if len(deps) != 2 {
		t.Errorf("expected 2 dependencies for a.go, got %d", len(deps))
	}

	revs := g.DependentsOf("c.go")
	if len(revs) != 2 {
		t.Errorf("expected 2 dependents for c.go, got %d", len(revs))
	}
}

func TestDepGraph_NoDuplicateEdges(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.go", "b.go")
	g.AddEdge("a.go", "b.go") // duplicate

	if len(g.DependenciesOf("a.go")) != 1 {
		t.Errorf("expected 1 dependency (no dups), got %d", len(g.DependenciesOf("a.go")))
	}
}

func TestDepGraph_DetectCircular(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.go", "b.go")
	g.AddEdge("b.go", "a.go") // circular!

	cycles := g.DetectCircular()
	if len(cycles) == 0 {
		t.Error("expected circular dependency to be detected")
	}
}

func TestDepGraph_NoCircular(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.go", "b.go")
	g.AddEdge("b.go", "c.go")

	cycles := g.DetectCircular()
	if len(cycles) != 0 {
		t.Errorf("expected no cycles, got %d", len(cycles))
	}
}

func TestDepGraph_RenderMarkdown(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("main.go", "handler.go")
	g.AddEdge("handler.go", "model.go")

	md := g.RenderMarkdown()
	if !strings.Contains(md, "## Dependency Graph") {
		t.Error("expected markdown to contain header")
	}
	if !strings.Contains(md, "main.go") {
		t.Error("expected markdown to contain main.go")
	}
	if !strings.Contains(md, "handler.go") {
		t.Error("expected markdown to contain handler.go")
	}
}

func TestDepGraph_RenderMarkdown_Empty(t *testing.T) {
	g := NewDepGraph()
	md := g.RenderMarkdown()
	if md != "" {
		t.Errorf("expected empty string for empty graph, got %q", md)
	}
}

func TestDepGraph_RenderMarkdown_Circular(t *testing.T) {
	g := NewDepGraph()
	g.AddEdge("a.go", "b.go")
	g.AddEdge("b.go", "a.go")

	md := g.RenderMarkdown()
	if !strings.Contains(md, "Circular") {
		t.Error("expected markdown to mention circular dependencies")
	}
}
