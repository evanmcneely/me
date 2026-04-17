package blog

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"evanmcneely/internal/cache"
	"evanmcneely/internal/render"

	"gopkg.in/yaml.v3"
)

var validSlug = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

type Service struct {
	postsDir    string
	pagesDir    string
	tooltipsDir string
	cache       *cache.Store
	renderer    *render.MarkdownRenderer
}

type postFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Date        string `yaml:"date"`
}

type pageFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Created     string `yaml:"created"`
	Updated     string `yaml:"updated"`
}

type tooltipFrontmatter struct {
	Title string `yaml:"title"`
}

func NewService(contentDir string, store *cache.Store, renderer *render.MarkdownRenderer) *Service {
	return &Service{
		postsDir:    filepath.Join(contentDir, "posts"),
		pagesDir:    filepath.Join(contentDir, "pages"),
		tooltipsDir: filepath.Join(contentDir, "tooltips"),
		cache:       store,
		renderer:    renderer,
	}
}

func (s *Service) ListPosts(ctx context.Context) ([]Post, error) {
	entries, err := os.ReadDir(s.postsDir)
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0, len(entries))
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		post, err := s.Post(ctx, slug)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PublishedAt.After(posts[j].PublishedAt)
	})

	return posts, nil
}

func (s *Service) ListPages(ctx context.Context) ([]Page, error) {
	entries, err := os.ReadDir(s.pagesDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	pages := make([]Page, 0, len(entries))
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		page, err := s.Page(ctx, slug)
		if err != nil {
			continue
		}
		pages = append(pages, page)
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Title < pages[j].Title
	})

	return pages, nil
}

func (s *Service) Post(_ context.Context, slug string) (Post, error) {
	if !validSlug.MatchString(slug) {
		return Post{}, fs.ErrNotExist
	}

	path := filepath.Join(s.postsDir, slug+".md")
	info, err := os.Stat(path)
	if err != nil {
		return Post{}, err
	}

	if cached, err := s.cache.GetPost(slug, info.ModTime().Unix(), info.Size()); err == nil && cached != nil {
		return hydratePost(*cached), nil
	} else if err != nil {
		return Post{}, err
	}

	rendered, err := s.renderPost(path, slug, info)
	if err != nil {
		return Post{}, err
	}
	if err := s.cache.UpsertPost(rendered); err != nil {
		return Post{}, err
	}

	return hydratePost(rendered), nil
}

func (s *Service) Page(_ context.Context, slug string) (Page, error) {
	if !validSlug.MatchString(slug) {
		return Page{}, fs.ErrNotExist
	}

	path := filepath.Join(s.pagesDir, slug+".md")
	info, err := os.Stat(path)
	if err != nil {
		return Page{}, err
	}

	if cached, err := s.cache.GetPage(slug, info.ModTime().Unix(), info.Size()); err == nil && cached != nil {
		return hydratePage(*cached), nil
	} else if err != nil {
		return Page{}, err
	}

	rendered, err := s.renderPage(path, slug, info)
	if err != nil {
		return Page{}, err
	}
	if err := s.cache.UpsertPage(rendered); err != nil {
		return Page{}, err
	}

	return hydratePage(rendered), nil
}

func (s *Service) Tooltip(_ context.Context, slug string) (Tooltip, error) {
	if !validSlug.MatchString(slug) {
		return Tooltip{}, fs.ErrNotExist
	}

	path := filepath.Join(s.tooltipsDir, slug+".md")
	info, err := os.Stat(path)
	if err != nil {
		return Tooltip{}, err
	}

	if cached, err := s.cache.GetTooltip(slug, info.ModTime().Unix(), info.Size()); err == nil && cached != nil {
		return Tooltip{Slug: cached.Slug, Title: cached.Title, HTML: template.HTML(cached.HTML)}, nil
	} else if err != nil {
		return Tooltip{}, err
	}

	rendered, err := s.renderTooltip(path, slug, info)
	if err != nil {
		return Tooltip{}, err
	}
	if err := s.cache.UpsertTooltip(rendered); err != nil {
		return Tooltip{}, err
	}

	return Tooltip{Slug: rendered.Slug, Title: rendered.Title, HTML: template.HTML(rendered.HTML)}, nil
}

func (s *Service) renderPost(path, slug string, info os.FileInfo) (cache.CachedPost, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return cache.CachedPost{}, err
	}

	metaBlock, body, err := splitFrontmatter(input)
	if err != nil {
		return cache.CachedPost{}, err
	}

	var meta postFrontmatter
	if err := yaml.Unmarshal(metaBlock, &meta); err != nil {
		return cache.CachedPost{}, err
	}

	html, err := s.renderer.Render(string(body))
	if err != nil {
		return cache.CachedPost{}, err
	}

	plain := render.PlainText(string(body))
	publishedAt, err := parseDate(meta.Date)
	if err != nil {
		return cache.CachedPost{}, err
	}
	description := strings.TrimSpace(meta.Description)
	if description == "" {
		description = excerpt(plain, 170)
	}

	return cache.CachedPost{
		Slug:        slug,
		SourcePath:  path,
		ModUnix:     info.ModTime().Unix(),
		Size:        info.Size(),
		Title:       strings.TrimSpace(meta.Title),
		Description: description,
		Author:      strings.TrimSpace(meta.Author),
		PublishedAt: publishedAt,
		HTML:        html,
		Excerpt:     excerpt(plain, 220),
		ReadTime:    readTime(plain),
	}, nil
}

func (s *Service) renderTooltip(path, slug string, info os.FileInfo) (cache.CachedTooltip, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return cache.CachedTooltip{}, err
	}

	metaBlock, body, err := splitFrontmatter(input)
	if err != nil {
		return cache.CachedTooltip{}, err
	}

	var meta tooltipFrontmatter
	if len(metaBlock) > 0 {
		if err := yaml.Unmarshal(metaBlock, &meta); err != nil {
			return cache.CachedTooltip{}, err
		}
	}

	html, err := s.renderer.Render(string(body))
	if err != nil {
		return cache.CachedTooltip{}, err
	}

	return cache.CachedTooltip{
		Slug:       slug,
		SourcePath: path,
		ModUnix:    info.ModTime().Unix(),
		Size:       info.Size(),
		Title:      strings.TrimSpace(meta.Title),
		HTML:       html,
	}, nil
}

func (s *Service) renderPage(path, slug string, info os.FileInfo) (cache.CachedPage, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return cache.CachedPage{}, err
	}

	metaBlock, body, err := splitFrontmatter(input)
	if err != nil {
		return cache.CachedPage{}, err
	}

	var meta pageFrontmatter
	if err := yaml.Unmarshal(metaBlock, &meta); err != nil {
		return cache.CachedPage{}, err
	}

	html, err := s.renderer.Render(string(body))
	if err != nil {
		return cache.CachedPage{}, err
	}

	plain := render.PlainText(string(body))
	createdAt, err := parseDate(meta.Created)
	if err != nil {
		return cache.CachedPage{}, err
	}
	updatedAt, err := parseDate(meta.Updated)
	if err != nil {
		return cache.CachedPage{}, err
	}
	if createdAt.IsZero() {
		createdAt = info.ModTime()
	}
	if updatedAt.IsZero() {
		updatedAt = info.ModTime()
	}
	description := strings.TrimSpace(meta.Description)
	if description == "" {
		description = excerpt(plain, 170)
	}
	title := strings.TrimSpace(meta.Title)
	if title == "" {
		title = slugToTitle(slug)
	}

	return cache.CachedPage{
		Slug:        slug,
		SourcePath:  path,
		ModUnix:     info.ModTime().Unix(),
		Size:        info.Size(),
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		HTML:        html,
		Excerpt:     excerpt(plain, 220),
		ReadTime:    readTime(plain),
	}, nil
}

func hydratePost(cached cache.CachedPost) Post {
	return Post{
		Slug:        cached.Slug,
		Title:       cached.Title,
		Description: cached.Description,
		Author:      cached.Author,
		PublishedAt: cached.PublishedAt,
		HTML:        template.HTML(cached.HTML),
		Excerpt:     cached.Excerpt,
		ReadTime:    cached.ReadTime,
	}
}

func hydratePage(cached cache.CachedPage) Page {
	return Page{
		Slug:        cached.Slug,
		Title:       cached.Title,
		Description: cached.Description,
		CreatedAt:   cached.CreatedAt,
		UpdatedAt:   cached.UpdatedAt,
		HTML:        template.HTML(cached.HTML),
		Excerpt:     cached.Excerpt,
		ReadTime:    cached.ReadTime,
	}
}

func splitFrontmatter(input []byte) ([]byte, []byte, error) {
	text := string(input)
	if !strings.HasPrefix(text, "---\n") {
		return nil, input, nil
	}

	parts := strings.SplitN(text[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return nil, nil, errors.New("invalid frontmatter")
	}

	return []byte(parts[0]), []byte(parts[1]), nil
}

func parseDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse date %q: %w", value, err)
	}
	return date, nil
}

func excerpt(text string, limit int) string {
	if len(text) <= limit {
		return text
	}
	cut := strings.LastIndex(text[:limit], " ")
	if cut < 1 {
		cut = limit
	}
	return strings.TrimSpace(text[:cut]) + "..."
}

func readTime(text string) int {
	words := len(strings.Fields(text))
	minutes := words / 200
	if words%200 != 0 {
		minutes++
	}
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}

func slugToTitle(slug string) string {
	parts := strings.Split(slug, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
