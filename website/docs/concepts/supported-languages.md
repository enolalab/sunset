---
title: Supported languages
---

Sunset detects languages from a filename extension, case-insensitively.

| Language | ID | Extensions |
| --- | --- | --- |
| Go | `go` | `.go` |
| JavaScript | `javascript` | `.js`, `.jsx` |
| TypeScript | `typescript` | `.ts`, `.tsx` |
| Python | `python` | `.py` |

Unsupported extensions are not included in directory scans. Parsing a single file with an unsupported extension returns an unsupported-language error unless the public Go API is given `WithLanguage`.
