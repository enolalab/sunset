// Package docstring extracts documentation comments from source code
// for various programming languages.
package docstring

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// Extract returns the documentation comment associated with a node.
// It looks for comment nodes immediately preceding the given node.
// Returns empty string if no doc comment is found.
func Extract(node *tree_sitter.Node, source []byte, languageID string) string {
	switch languageID {
	case "go":
		return extractGo(node, source)
	case "javascript", "typescript":
		return extractJSDoc(node, source)
	case "python":
		return extractPythonDocstring(node, source)
	default:
		return extractGo(node, source) // fallback to line comments
	}
}

// extractGo extracts Go-style // line comments above a node.
func extractGo(node *tree_sitter.Node, source []byte) string {
	parent := node.Parent()
	if parent == nil {
		return ""
	}

	// Find this node's index among siblings
	idx := -1
	for i := 0; i < int(parent.ChildCount()); i++ {
		sibling := parent.Child(uint(i))
		if sibling != nil && sibling.StartByte() == node.StartByte() && sibling.EndByte() == node.EndByte() {
			idx = i
			break
		}
	}
	if idx <= 0 {
		return ""
	}

	// Collect consecutive comment nodes before this node
	var comments []string
	for i := idx - 1; i >= 0; i-- {
		prev := parent.Child(uint(i))
		if prev == nil {
			break
		}
		if prev.Kind() == "comment" {
			text := prev.Utf8Text(source)
			// Strip leading "// " or "//"
			if len(text) > 3 && text[:3] == "// " {
				text = text[3:]
			} else if len(text) > 2 && text[:2] == "//" {
				text = text[2:]
			}
			comments = append([]string{text}, comments...)
		} else {
			break
		}
	}

	if len(comments) == 0 {
		return ""
	}

	result := ""
	for i, c := range comments {
		if i > 0 {
			result += " "
		}
		result += c
	}
	return result
}

// extractJSDoc extracts /** JSDoc */ comments above a node.
func extractJSDoc(node *tree_sitter.Node, source []byte) string {
	parent := node.Parent()
	if parent == nil {
		return ""
	}

	idx := -1
	for i := 0; i < int(parent.ChildCount()); i++ {
		sibling := parent.Child(uint(i))
		if sibling != nil && sibling.StartByte() == node.StartByte() && sibling.EndByte() == node.EndByte() {
			idx = i
			break
		}
	}
	if idx <= 0 {
		return ""
	}

	prev := parent.Child(uint(idx - 1))
	if prev == nil {
		return ""
	}

	if prev.Kind() == "comment" {
		text := prev.Utf8Text(source)
		return cleanJSDoc(text)
	}

	return ""
}

// cleanJSDoc strips /** */ delimiters and * prefixes from JSDoc.
func cleanJSDoc(text string) string {
	if len(text) < 5 {
		return text
	}

	// Remove /** and */
	if text[:3] == "/**" {
		text = text[3:]
	}
	if len(text) > 2 && text[len(text)-2:] == "*/" {
		text = text[:len(text)-2]
	}

	// Clean up line prefixes
	result := ""
	lines := splitLines(text)
	for _, line := range lines {
		line = trimLeft(line, " \t")
		if len(line) > 0 && line[0] == '*' {
			line = line[1:]
			line = trimLeft(line, " ")
		}
		line = trimRight(line, " \t")
		if line != "" {
			if result != "" {
				result += " "
			}
			result += line
		}
	}
	return result
}

// extractPythonDocstring extracts """docstring""" from inside a function/class body.
func extractPythonDocstring(node *tree_sitter.Node, source []byte) string {
	// Python docstrings are the first expression in a function/class body
	var body *tree_sitter.Node

	// Look for "body" or "block" child
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(uint(i))
		if child != nil && (child.Kind() == "block" || child.Kind() == "body") {
			body = child
			break
		}
	}

	if body == nil || body.ChildCount() == 0 {
		return ""
	}

	// First statement in body
	first := body.NamedChild(0)
	if first == nil {
		return ""
	}

	// Check if it's an expression_statement containing a string
	if first.Kind() == "expression_statement" && first.NamedChildCount() > 0 {
		expr := first.NamedChild(0)
		if expr != nil && expr.Kind() == "string" {
			text := expr.Utf8Text(source)
			return cleanPythonDocstring(text)
		}
	}

	return ""
}

// cleanPythonDocstring strips triple quotes from docstrings.
func cleanPythonDocstring(text string) string {
	// Remove """ or '''
	for _, delim := range []string{`"""`, `'''`} {
		if len(text) >= 6 && text[:3] == delim && text[len(text)-3:] == delim {
			text = text[3 : len(text)-3]
			break
		}
	}
	text = trimLeft(text, " \t\n")
	text = trimRight(text, " \t\n")

	// Collapse to single line for frontmatter
	result := ""
	for _, line := range splitLines(text) {
		line = trimLeft(line, " \t")
		line = trimRight(line, " \t")
		if line != "" {
			if result != "" {
				result += " "
			}
			result += line
		}
	}
	return result
}

// Helper functions to avoid strings package dependency
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimLeft(s string, cutset string) string {
	for len(s) > 0 && containsByte(cutset, s[0]) {
		s = s[1:]
	}
	return s
}

func trimRight(s string, cutset string) string {
	for len(s) > 0 && containsByte(cutset, s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return s
}

func containsByte(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}
