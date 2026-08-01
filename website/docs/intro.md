---
sidebar_position: 1
title: Sunset
slug: /
---

Sunset scans supported source files and writes Markdown documentation backed by tree-sitter concrete syntax trees. It is useful for codebase exploration, documentation inputs, and AI-oriented indexes.

Start with [installation](getting-started/installation.md), then run [`sunset parse`](cli/parse.md) against a project. The generated documentation is local to the project under `.sunset/` unless an output path is supplied.

Sunset supports Go, JavaScript, TypeScript, and Python. Its command-line indexer and public Go parsing API are separate surfaces; the Go API parses one file at a time.
