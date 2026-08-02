---
title: Installation
---

## Releases

The current public release is **v1.0.0**. Download stable binaries and
`checksums.txt` from the [GitHub Releases page](https://github.com/enolalabs/sunset/releases),
which is the authoritative record of public release availability.

The release process creates and verifies a **draft** before a maintainer
publishes it. Draft assets, including planned `v1.0.1` assets, are not public
downloads and must not be treated as an installation source.

## Development installation from source

The canonical Go module is `github.com/enolalabs/sunset`. To use the current
main branch for development, clone the public repository and install from that
checkout. This is not a stable, version-pinned installation method.

```bash
git clone https://github.com/enolalabs/sunset.git
cd sunset
go install ./cmd/sunset
sunset version
```

## Local build

```bash
git clone https://github.com/enolalabs/sunset.git
cd sunset
make build
./bin/sunset version
```
