---
title: Known limitations
---

These behaviors reflect the current implementation:

- The dependency graph code exists internally but is not emitted by the indexing engine, so generated output does not contain a dependency graph.
- On a cached or `update` run, `index.md` is built from newly parsed files rather than every cached file. It is not a complete view of unchanged files.
- When a source file is deleted, its cache entry is pruned but its existing generated Markdown in `output/files/` is not removed.
- `parse --detail` selects full CST output only when its value is exactly `full`; invalid values currently fall back to summary output instead of failing validation.

Use a clean output directory and a non-cached parse when a complete, current output set is required.
