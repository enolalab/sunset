---
title: Installation
---

## From source

The canonical Go module is `github.com/enolalabs/sunset`.

```bash
go install github.com/enolalabs/sunset/cmd/sunset@latest
sunset version
```

## Releases

Published release binaries are available from the [GitHub Releases page](https://github.com/enolalabs/sunset/releases). A public release contains platform archives and `checksums.txt`; use version-explicit download URLs when verifying an archive.

The release process creates and verifies a **draft** before a maintainer publishes it. Draft assets, including the planned `v1.0.1` assets, are not public downloads and must not be treated as an installation source.

## Local build

```bash
git clone https://github.com/enolalabs/sunset.git
cd sunset
make build
./bin/sunset version
```
