---
title: Language detection
---

```go
func Languages() []LanguageInfo
```

`Languages` returns a copy of the registered language descriptions. Each `LanguageInfo` includes `Name`, `ID`, and `Extensions`.

```go
for _, language := range sunset.Languages() {
    fmt.Printf("%s: %v\n", language.ID, language.Extensions)
}
```

`ParseFile` detects from a file extension by default. If an extension does not identify a registered language, pass `WithLanguage` to specify a supported ID explicitly.
