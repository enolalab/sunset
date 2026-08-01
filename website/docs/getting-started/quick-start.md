---
title: Quick start
---

Run a summary parse from the project root:

```bash
sunset parse .
```

Sunset writes an index to `.sunset/output/index.md` and one Markdown file per parsed source file under `.sunset/output/files/`.

```bash
# Inspect a complete concrete syntax tree instead of a summary.
sunset parse . --detail full

# Reprocess only changed files in an existing index.
sunset update .

# Remove both generated output and the cache.
sunset clean .
```

See [parse](../cli/parse.md) for its flags and [limitations](../reference/limitations.md) before relying on incremental output as a complete project index.
