---
title: Scanning and exclusions
---

Directory scans recurse from the supplied root. Sunset always skips `.git`, `.idea`, `.vscode`, `node_modules`, `__pycache__`, `vendor`, and every other hidden directory whose name starts with `.`.

If a root `.gitignore` exists, Sunset applies its patterns, including negation, directory-only, and anchored patterns. It then applies `--exclude` patterns. An exclude is matched against both the basename and the relative path; simple suffix patterns such as `*.test.ts` are also supported.

Unreadable entries are skipped. Only files with a supported extension are returned by the scanner.
