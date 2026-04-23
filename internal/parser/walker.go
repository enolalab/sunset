package parser

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// WalkFunc is the callback type for tree traversal.
// Return true to continue walking children, false to skip.
type WalkFunc func(node *tree_sitter.Node, depth int) bool

// Walk performs a depth-first traversal of the tree starting from the given node.
// The callback receives each node and its depth (root = 0).
// If the callback returns false, children of that node are skipped.
func Walk(node *tree_sitter.Node, fn WalkFunc) {
	walkRecursive(node, 0, fn)
}

func walkRecursive(node *tree_sitter.Node, depth int, fn WalkFunc) {
	if node == nil {
		return
	}

	if !fn(node, depth) {
		return
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(uint(i))
		if child != nil {
			walkRecursive(child, depth+1, fn)
		}
	}
}

// Filter returns all descendant nodes matching the given node type.
func Filter(node *tree_sitter.Node, nodeType string) []*tree_sitter.Node {
	var result []*tree_sitter.Node
	Walk(node, func(n *tree_sitter.Node, depth int) bool {
		if n.Kind() == nodeType {
			result = append(result, n)
		}
		return true
	})
	return result
}

// Depth calculates the depth of a node relative to the tree root.
// Root node has depth 0.
func Depth(node *tree_sitter.Node) int {
	depth := 0
	current := node.Parent()
	for current != nil {
		depth++
		current = current.Parent()
	}
	return depth
}

// CountNodes returns the total number of nodes in the subtree rooted at node.
func CountNodes(node *tree_sitter.Node) int {
	count := 0
	Walk(node, func(n *tree_sitter.Node, depth int) bool {
		count++
		return true
	})
	return count
}

// MaxDepth returns the maximum depth of the subtree rooted at node.
func MaxDepth(node *tree_sitter.Node) int {
	maxD := 0
	Walk(node, func(n *tree_sitter.Node, depth int) bool {
		if depth > maxD {
			maxD = depth
		}
		return true
	})
	return maxD
}
