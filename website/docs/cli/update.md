---
title: update
---

```text
sunset update [path] [flags]
```

`update` hashes discovered files and reparses only files whose content changed since the cache entry. It always uses the cache and always writes **summary** output. It does not accept `--detail`, `--no-cache`, or `--max-depth`.

| Flag | Default | Meaning |
| --- | --- | --- |
| `--output` | `<path>/.sunset/output` | Destination directory. |
| `--exclude` | empty | Comma-separated filename or relative-path glob patterns. |
| `--concurrency` | CPU count | Maximum parallel parsers. |
| `--quiet` | `false` | Suppress normal progress output. |

```bash
sunset update . --exclude "*_test.go" --quiet
```

An incremental index contains only newly parsed file metadata, not all cached files. Removed cache entries do not remove old per-file Markdown; see [limitations](../reference/limitations.md).
