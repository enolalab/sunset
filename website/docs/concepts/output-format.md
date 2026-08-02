---
title: Output format
---

Default output is `.sunset/output/` beneath the parsed root. `index.md` contains project frontmatter and a file table. Per-file output is written to `files/<sanitized-relative-path>.md`; path separators become underscores.

Every per-file document starts with YAML frontmatter:

```yaml
---
file: main.go
language: go
package: main
lines: 28
function_count: 2
type_count: 0
import_count: 3
tags:
  - has-functions
---
```

Summary mode renders extracted functions, types, imports, and documentation where supported. `parse --detail full` instead renders a concrete syntax tree. `update` always writes summary output.
