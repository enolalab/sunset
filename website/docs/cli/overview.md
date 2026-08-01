---
title: CLI overview
---

```text
sunset <command> [flags]
```

| Command | Purpose |
| --- | --- |
| `parse [directory]` | Scan and document a directory root; defaults to the current directory. |
| `update [path]` | Reparse changed files incrementally. |
| `languages` | List supported language IDs and extensions. |
| `version` | Print the Sunset version. |
| `clean [path]` | Remove the project's `.sunset` cache and output. |
| `help` | Print root usage. |

`--help` and `--version` are also accepted at the root. Commands exit nonzero for unknown commands and parse runs that collect file errors.

## Current limitation

The CLI parse workflow supports directory roots only. Single-file parsing is a
current limitation of the CLI; use the public Go API's `ParseFile` for a single
source file.
