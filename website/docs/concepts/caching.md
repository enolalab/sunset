---
title: Caching
---

Sunset caches SHA-256 content hashes at `.sunset/cache/cache.json`. A cache entry records the hash, language ID, and update time for a relative source path.

`parse` uses the cache by default; pass `--no-cache` to parse all discovered files without reading or saving it. `update` always uses the cache and has no `--no-cache` option. `clean` deletes the cache and generated output together.

The cache tracks changes, not a complete output manifest. Read the [incremental limitations](../reference/limitations.md) before treating an update run as a full rebuild.
