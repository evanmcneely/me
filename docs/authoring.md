# Authoring notes

Posts live in `content/posts`, pages live in `content/pages`, and tooltips live in `content/tooltips`.

## Frontmatter

Each post can include YAML frontmatter:

```yaml
---
title: Example title
description: Short summary for the home page.
author: Evan McNeely
date: 2026-04-15
---
```

Pages work the same way, but they publish at the site root using the markdown filename as the slug. For example, `content/pages/about.md` becomes `/about`.

```yaml
---
title: About
description: What this site is for.
created: 2026-04-01
updated: 2026-04-16
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

For live reload while writing posts or pages, use:

```bash
./scripts/dev.sh
```

Then open `http://localhost:3000`.
