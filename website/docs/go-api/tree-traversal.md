---
title: Tree traversal
---

`TreeWrapper` provides `RootNode`, `Walk`, and `Filter`. `Walk` is depth-first; return `false` from its callback to skip a node's children.

```go
result.Tree.Walk(func(node *sunset.Node, depth int) bool {
    fmt.Printf("%s at %d:%d\n", node.Type(), node.Position().Row, node.Position().Column)
    return true
})
```

`Node` provides `Raw`, `Type`, `Text`, `Position`, `EndPosition`, `Children`, `NamedChildren`, `ChildByFieldName`, `Parent`, `IsNamed`, `IsError`, `IsMissing`, and `ChildCount`. `Raw` exposes the underlying tree-sitter node for advanced use.

`Position.Row` and `Position.Column` are **zero-based**, matching tree-sitter positions. `Filter(nodeType)` returns every node whose grammar type matches `nodeType`.
