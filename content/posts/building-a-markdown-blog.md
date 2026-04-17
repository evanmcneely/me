---
title: Building a markdown blog in Go
description: Notes on a tiny blog stack with markdown posts, in-memory render caching, and inline hover notes.
author: Evan McNeely
date: 2026-04-15
tags:
  - go
  - markdown
---

The basic loop is intentionally small: write markdown, render it in Go, and keep the rendered HTML in memory so the server does not have to re-parse every post on every request.

The one extra flourish is {{curled underlines|curled-underlines}}. They let a paragraph carry extra context without sending the reader on a side quest.

## Why this shape

- Markdown files are easy to diff and easy to move.
- An in-memory cache makes rendered output cheap to reuse.
- A tiny custom syntax gives inline definitions without a heavyweight CMS.

## Code still matters

```go
func render(slug string) error {
	post, err := service.Post(context.Background(), slug)
	if err != nil {
		return err
	}

	fmt.Println(post.Title)
	return nil
}
```

## Authoring a tooltip

Tooltip notes live in `content/tooltips/*.md`, so a phrase like {{render cache|render-cache}} can be reused across multiple posts.

That keeps the main essay readable while still leaving room for definitions, tangents, and implementation details.
