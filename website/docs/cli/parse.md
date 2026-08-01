---
title: parse
---

```text
sunset parse <path> [flags]
```

Parse scans supported files and writes Markdown. With no path it uses the current directory. Default output is `<path>/.sunset/output`.

| Flag | Default | Meaning |
| --- | --- | --- |
| `--output` | `<path>/.sunset/output` | Destination directory. |
| `--detail` | `summary` | `summary` or `full`; `full` writes a CST dump. |
| `--exclude` | empty | Comma-separated filename or relative-path glob patterns. |
| `--concurrency` | CPU count | Maximum parallel parsers. |
| `--no-cache` | `false` | Do not load or save the cache; parse every discovered file. |
| `--max-depth` | `0` | CST depth limit only when `--detail full`; zero is unlimited. |
| `--quiet` | `false` | Suppress normal progress output. |

Only `--detail`, `--no-cache`, and `--max-depth` belong to `parse`. Values other than `full` for `--detail` currently fall back to summary output; validate user input yourself if that distinction matters.

```bash
sunset parse ./service --detail full --max-depth 4 --exclude "*_test.go,dist/*.js"
```
