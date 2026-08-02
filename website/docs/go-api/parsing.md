---
title: Parsing files
---

```go
func ParseFile(path string, opts ...Option) (*FileResult, error)
func WithLanguage(id string) Option
```

`ParseFile` reads `path`, detects its language from the basename extension, parses it, and returns a `FileResult`. It returns an error when the file cannot be read, the language is unsupported, or tree-sitter setup or parsing fails.

Use `WithLanguage` to override extension detection with a registered language ID such as `go`, `javascript`, `typescript`, or `python`.

```go
result, err := sunset.ParseFile("script.source", sunset.WithLanguage("python"))
if err != nil {
    return err
}
defer result.Close()
```

Call `Close` exactly when the result is no longer needed. `HasErrors` reports whether the parsed tree contains syntax errors.
