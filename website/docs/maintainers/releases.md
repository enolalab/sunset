---
title: Releases
---

Releases are an explicit maintainer-controlled process. A candidate commit must be reachable from `main`, pass release validation, have canonical module identity `github.com/enolalabs/sunset`, and have matching release notes and install snippets.

Pushing an exact `vX.Y.Z` tag runs the release workflow. It builds and validates five native archives, creates or refreshes a draft release, verifies its checksums and an external Go consumer, then waits for authorization in the protected `release` environment before publishing.

The expected archive targets are Linux amd64 and arm64, macOS amd64 and arm64, and Windows amd64, plus `checksums.txt`. Draft releases are private staging artifacts. They are not public downloads and must not be presented as published assets.

After promotion, published assets and their version-explicit GitHub URLs are immutable. Recovery requires a new version, not editing or replacing the published release. See `docs/releasing.md` in the repository for the operator runbook.
