package analyzer

import (
	"fmt"
	"sort"
	"strings"
)

// DepGraph represents a file-level dependency graph.
type DepGraph struct {
	// Edges maps a file to its imported files (adjacency list).
	Edges map[string][]string
	// ReverseEdges maps a file to files that import it.
	ReverseEdges map[string][]string
	// edgeSet tracks existing edges for O(1) dedup.
	edgeSet map[string]map[string]bool
}

// NewDepGraph creates an empty dependency graph.
func NewDepGraph() *DepGraph {
	return &DepGraph{
		Edges:        make(map[string][]string),
		ReverseEdges: make(map[string][]string),
		edgeSet:      make(map[string]map[string]bool),
	}
}

// AddEdge adds a dependency from source to target.
func (g *DepGraph) AddEdge(source, target string) {
	// O(1) dedup check
	if g.edgeSet[source] == nil {
		g.edgeSet[source] = make(map[string]bool)
	}
	if g.edgeSet[source][target] {
		return
	}
	g.edgeSet[source][target] = true
	g.Edges[source] = append(g.Edges[source], target)
	g.ReverseEdges[target] = append(g.ReverseEdges[target], source)
}

// DependenciesOf returns all files that the given file depends on.
func (g *DepGraph) DependenciesOf(file string) []string {
	return g.Edges[file]
}

// DependentsOf returns all files that depend on the given file.
func (g *DepGraph) DependentsOf(file string) []string {
	return g.ReverseEdges[file]
}

// DetectCircular finds circular dependencies in the graph.
// Returns a list of cycles, each cycle is a list of files forming the loop.
func (g *DepGraph) DetectCircular() [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	for node := range g.Edges {
		if !visited[node] {
			var path []string
			g.dfs(node, visited, inStack, path, &cycles)
		}
	}

	return cycles
}

func (g *DepGraph) dfs(node string, visited, inStack map[string]bool, path []string, cycles *[][]string) {
	visited[node] = true
	inStack[node] = true
	path = append(path, node)

	for _, neighbor := range g.Edges[node] {
		if !visited[neighbor] {
			g.dfs(neighbor, visited, inStack, path, cycles)
		} else if inStack[neighbor] {
			// Found a cycle — extract it
			cycleStart := -1
			for i, n := range path {
				if n == neighbor {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 {
				cycle := make([]string, len(path)-cycleStart)
				copy(cycle, path[cycleStart:])
				*cycles = append(*cycles, cycle)
			}
		}
	}

	inStack[node] = false
}

// RenderMarkdown generates a Markdown table of dependencies.
func (g *DepGraph) RenderMarkdown() string {
	if len(g.Edges) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n## Dependency Graph\n\n")
	b.WriteString("| File | Depends On | Depended By |\n")
	b.WriteString("|---|---|---|\n")

	// Collect all files
	allFiles := make(map[string]bool)
	for f := range g.Edges {
		allFiles[f] = true
	}
	for f := range g.ReverseEdges {
		allFiles[f] = true
	}

	// Collect and sort all files for deterministic output
	files := make([]string, 0, len(allFiles))
	for f := range allFiles {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, f := range files {
		deps := g.Edges[f]
		revs := g.ReverseEdges[f]
		depsStr := "-"
		if len(deps) > 0 {
			depsStr = strings.Join(deps, ", ")
		}
		revsStr := "-"
		if len(revs) > 0 {
			revsStr = strings.Join(revs, ", ")
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", f, depsStr, revsStr))
	}

	// Report circular dependencies
	cycles := g.DetectCircular()
	if len(cycles) > 0 {
		b.WriteString("\n### ⚠️ Circular Dependencies\n\n")
		for _, cycle := range cycles {
			b.WriteString(fmt.Sprintf("- %s → %s\n", strings.Join(cycle, " → "), cycle[0]))
		}
	}

	return b.String()
}
