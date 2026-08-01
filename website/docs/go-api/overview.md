---
title: Go API overview
---

Import the public package with its canonical module path:

```go
import "github.com/enolalabs/sunset/pkg/sunset"
```

The public API parses one source file at a time. Its main entry point is `ParseFile`; use `Languages` to inspect the registry. A successful `FileResult` owns tree-sitter resources, so callers must call `Close` when finished.

```go
result, err := sunset.ParseFile("main.go")
if err != nil {
    return err
}
defer result.Close()
```

`FileResult` exposes `Path`, `Language`, `LanguageID`, `Source`, and `Tree`, plus `HasErrors` and `Close`. `Tree` is a `*TreeWrapper`, not a function. There is no public package-level `Walk`, and nodes use `Type`, not `Kind`.
