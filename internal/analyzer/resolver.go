package analyzer

import (
	"os"
	"path/filepath"
	"strings"
)

// resolveImport attempts to resolve an import path to a project-relative file path.
func resolveImport(importPath string, langID string, currentFile string, projectFiles []string, rootDir string) string {
	switch langID {
	case "go":
		return resolveGoImport(importPath, projectFiles, rootDir)
	case "javascript", "typescript":
		return resolveJSImport(importPath, currentFile, projectFiles)
	case "python":
		return resolvePythonImport(importPath, currentFile, projectFiles)
	default:
		return ""
	}
}

// resolveGoImport resolves Go import paths.
// It reads go.mod to determine module path and maps internal imports to files.
func resolveGoImport(importPath string, projectFiles []string, rootDir string) string {
	modulePath := readGoModulePath(rootDir)
	if modulePath == "" {
		return ""
	}

	// Check if this import is within our module
	if !strings.HasPrefix(importPath, modulePath) {
		return "" // external
	}

	// Get the relative package path
	relPkg := strings.TrimPrefix(importPath, modulePath)
	relPkg = strings.TrimPrefix(relPkg, "/")

	// Find files in this package directory
	for _, f := range projectFiles {
		dir := filepath.Dir(f)
		if dir == relPkg || filepath.ToSlash(dir) == relPkg {
			return f // return first file in the package
		}
	}

	return ""
}

// readGoModulePath extracts the module path from go.mod.
func readGoModulePath(rootDir string) string {
	data, err := os.ReadFile(filepath.Join(rootDir, "go.mod"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

// resolveJSImport resolves JS/TS relative imports.
func resolveJSImport(importPath string, currentFile string, projectFiles []string) string {
	// Only handle relative imports
	if !strings.HasPrefix(importPath, ".") {
		return "" // external (npm package)
	}

	dir := filepath.Dir(currentFile)
	candidate := filepath.Join(dir, importPath)
	candidate = filepath.Clean(candidate)
	candidate = filepath.ToSlash(candidate)

	// Try exact match and common extensions
	extensions := []string{"", ".ts", ".tsx", ".js", ".jsx", "/index.ts", "/index.js"}
	for _, ext := range extensions {
		tryPath := candidate + ext
		for _, f := range projectFiles {
			if filepath.ToSlash(f) == tryPath {
				return f
			}
		}
	}

	return ""
}

// resolvePythonImport resolves Python relative/absolute imports.
func resolvePythonImport(importPath string, currentFile string, projectFiles []string) string {
	// Handle "from X import Y" — importPath is the module part
	parts := strings.Split(importPath, ".")

	// Relative import (starts with .)
	if strings.HasPrefix(importPath, ".") {
		dir := filepath.Dir(currentFile)
		// Count leading dots
		dots := 0
		for _, c := range importPath {
			if c == '.' {
				dots++
			} else {
				break
			}
		}
		// Go up directories
		for i := 1; i < dots; i++ {
			dir = filepath.Dir(dir)
		}
		rest := strings.TrimLeft(importPath, ".")
		if rest != "" {
			parts = strings.Split(rest, ".")
			candidate := filepath.Join(dir, filepath.Join(parts...))
			return findPythonFile(candidate, projectFiles)
		}
		return ""
	}

	// Absolute import — try as path from root
	candidate := filepath.Join(parts...)
	return findPythonFile(candidate, projectFiles)
}

func findPythonFile(candidate string, projectFiles []string) string {
	candidate = filepath.ToSlash(candidate)
	// Try as .py file
	tryPath := candidate + ".py"
	for _, f := range projectFiles {
		if filepath.ToSlash(f) == tryPath {
			return f
		}
	}
	// Try as package (__init__.py)
	tryPath = candidate + "/__init__.py"
	for _, f := range projectFiles {
		if filepath.ToSlash(f) == tryPath {
			return f
		}
	}
	return ""
}

// isExternalImport determines if an import is from an external package.
func isExternalImport(importPath string, langID string) bool {
	switch langID {
	case "go":
		// External if it contains a domain-like prefix
		return strings.Contains(importPath, ".") || isGoStdlib(importPath)
	case "javascript", "typescript":
		// External if not starting with . or /
		return !strings.HasPrefix(importPath, ".") && !strings.HasPrefix(importPath, "/")
	case "python":
		// Can't easily tell without inspecting the filesystem
		return false
	default:
		return false
	}
}

// isGoStdlib checks if a Go import path is from the standard library.
func isGoStdlib(path string) bool {
	stdPkgs := map[string]bool{
		"fmt": true, "os": true, "io": true, "net": true, "net/http": true,
		"encoding": true, "encoding/json": true, "encoding/xml": true,
		"strings": true, "strconv": true, "bytes": true, "bufio": true,
		"path": true, "path/filepath": true, "sync": true, "context": true,
		"testing": true, "log": true, "errors": true, "time": true,
		"math": true, "sort": true, "regexp": true, "reflect": true,
	}
	return stdPkgs[path]
}
