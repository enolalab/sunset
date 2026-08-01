---
title: CLI overview
---

```text
sunset <command> [flags]
```

| Command | Purpose |
| --- | --- |
| `parse <path>` | Scan and document a file or directory. |
| `update [path]` | Reparse changed files incrementally. |
| `languages` | List supported language IDs and extensions. |
| `version` | Print the Sunset version. |
| `clean [path]` | Remove the project's `.sunset` cache and output. |
| `help` | Print root usage. |

`--help` and `--version` are also accepted at the root. Commands exit nonzero for unknown commands and parse runs that collect file errors.
