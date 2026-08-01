---
title: Contributing
---

```bash
git clone https://github.com/enolalabs/sunset.git
cd sunset
make build
make test
make lint
```

Keep the public surface in `pkg/sunset` minimal, add tests for behavior changes, and close objects that wrap C resources. The project uses tree-sitter grammars, so language additions need a registry binding, a grammar dependency, representative test data, and any applicable docstring extraction support.

Before changing release instructions, run:

```bash
bash .github/scripts/release/check-doc-snippets.sh
```

The repository's `CONTRIBUTING.md` is the source for project contribution conventions.
