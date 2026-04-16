# Authoring notes

Posts live in `content/posts` and tooltips live in `content/tooltips`.

## Frontmatter

Each post can include YAML frontmatter:

```yaml
---
title: Example title
description: Short summary for the home page.
author: Evan McNeely
date: 2026-04-15
tags:
  - go
  - writing
---
```

## Tooltip syntax

Write inline tooltip references like this:

```md
The parser uses {{render caching|render-cache}} to avoid unnecessary work.
```

That looks for `content/tooltips/render-cache.md` and fetches it on hover or keyboard focus.

## Running locally

```bash
go run ./cmd/blog
```

Then open `http://localhost:8080`.
