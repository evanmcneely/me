---
title: Render cache
---

The cache stores rendered HTML plus post metadata keyed by slug and file freshness data.

If the source file has not changed, the server can reuse the cached render instead of reparsing markdown.
