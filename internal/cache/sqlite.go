package cache

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type CachedPost struct {
	Slug        string
	SourcePath  string
	ModUnix     int64
	Size        int64
	Title       string
	Description string
	Author      string
	PublishedAt time.Time
	Tags        []string
	HTML        string
	Excerpt     string
	ReadTime    int
}

type CachedTooltip struct {
	Slug       string
	SourcePath string
	ModUnix    int64
	Size       int64
	Title      string
	HTML       string
}

type CachedPage struct {
	Slug        string
	SourcePath  string
	ModUnix     int64
	Size        int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	HTML        string
	Excerpt     string
	ReadTime    int
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SQLiteStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SQLiteStore) GetPost(slug string, modUnix, size int64) (*CachedPost, error) {
	const query = `
		SELECT slug, source_path, mod_unix, size_bytes, title, description, author, published_at, tags_json, html, excerpt, read_time
		FROM rendered_posts
		WHERE slug = ? AND mod_unix = ? AND size_bytes = ?`

	row := s.db.QueryRow(query, slug, modUnix, size)
	var post CachedPost
	var tagsJSON string
	var publishedAt string
	if err := row.Scan(
		&post.Slug,
		&post.SourcePath,
		&post.ModUnix,
		&post.Size,
		&post.Title,
		&post.Description,
		&post.Author,
		&publishedAt,
		&tagsJSON,
		&post.HTML,
		&post.Excerpt,
		&post.ReadTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if publishedAt != "" {
		parsedTime, parseErr := time.Parse(time.RFC3339, publishedAt)
		if parseErr != nil {
			return nil, fmt.Errorf("parse cached post date: %w", parseErr)
		}
		post.PublishedAt = parsedTime
	}
	if err := json.Unmarshal([]byte(tagsJSON), &post.Tags); err != nil {
		return nil, fmt.Errorf("decode cached tags: %w", err)
	}

	return &post, nil
}

func (s *SQLiteStore) UpsertPost(post CachedPost) error {
	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO rendered_posts (
			slug, source_path, mod_unix, size_bytes, title, description, author, published_at, tags_json, html, excerpt, read_time
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			source_path = excluded.source_path,
			mod_unix = excluded.mod_unix,
			size_bytes = excluded.size_bytes,
			title = excluded.title,
			description = excluded.description,
			author = excluded.author,
			published_at = excluded.published_at,
			tags_json = excluded.tags_json,
			html = excluded.html,
			excerpt = excluded.excerpt,
			read_time = excluded.read_time`

	_, err = s.db.Exec(
		query,
		post.Slug,
		post.SourcePath,
		post.ModUnix,
		post.Size,
		post.Title,
		post.Description,
		post.Author,
		post.PublishedAt.Format(time.RFC3339),
		string(tagsJSON),
		post.HTML,
		post.Excerpt,
		post.ReadTime,
	)
	return err
}

func (s *SQLiteStore) GetTooltip(slug string, modUnix, size int64) (*CachedTooltip, error) {
	const query = `
		SELECT slug, source_path, mod_unix, size_bytes, title, html
		FROM rendered_tooltips
		WHERE slug = ? AND mod_unix = ? AND size_bytes = ?`

	row := s.db.QueryRow(query, slug, modUnix, size)
	var tooltip CachedTooltip
	if err := row.Scan(&tooltip.Slug, &tooltip.SourcePath, &tooltip.ModUnix, &tooltip.Size, &tooltip.Title, &tooltip.HTML); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &tooltip, nil
}

func (s *SQLiteStore) GetPage(slug string, modUnix, size int64) (*CachedPage, error) {
	const query = `
		SELECT slug, source_path, mod_unix, size_bytes, title, description, created_at, updated_at, html, excerpt, read_time
		FROM rendered_pages
		WHERE slug = ? AND mod_unix = ? AND size_bytes = ?`

	row := s.db.QueryRow(query, slug, modUnix, size)
	var page CachedPage
	var createdAt string
	var updatedAt string
	if err := row.Scan(
		&page.Slug,
		&page.SourcePath,
		&page.ModUnix,
		&page.Size,
		&page.Title,
		&page.Description,
		&createdAt,
		&updatedAt,
		&page.HTML,
		&page.Excerpt,
		&page.ReadTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if createdAt != "" {
		parsedTime, parseErr := time.Parse(time.RFC3339, createdAt)
		if parseErr != nil {
			return nil, fmt.Errorf("parse cached page created date: %w", parseErr)
		}
		page.CreatedAt = parsedTime
	}
	if updatedAt != "" {
		parsedTime, parseErr := time.Parse(time.RFC3339, updatedAt)
		if parseErr != nil {
			return nil, fmt.Errorf("parse cached page updated date: %w", parseErr)
		}
		page.UpdatedAt = parsedTime
	}

	return &page, nil
}

func (s *SQLiteStore) UpsertPage(page CachedPage) error {
	const query = `
		INSERT INTO rendered_pages (
			slug, source_path, mod_unix, size_bytes, title, description, created_at, updated_at, html, excerpt, read_time
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			source_path = excluded.source_path,
			mod_unix = excluded.mod_unix,
			size_bytes = excluded.size_bytes,
			title = excluded.title,
			description = excluded.description,
			created_at = excluded.created_at,
			updated_at = excluded.updated_at,
			html = excluded.html,
			excerpt = excluded.excerpt,
			read_time = excluded.read_time`

	_, err := s.db.Exec(
		query,
		page.Slug,
		page.SourcePath,
		page.ModUnix,
		page.Size,
		page.Title,
		page.Description,
		page.CreatedAt.Format(time.RFC3339),
		page.UpdatedAt.Format(time.RFC3339),
		page.HTML,
		page.Excerpt,
		page.ReadTime,
	)
	return err
}

func (s *SQLiteStore) UpsertTooltip(tooltip CachedTooltip) error {
	const query = `
		INSERT INTO rendered_tooltips (slug, source_path, mod_unix, size_bytes, title, html)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			source_path = excluded.source_path,
			mod_unix = excluded.mod_unix,
			size_bytes = excluded.size_bytes,
			title = excluded.title,
			html = excluded.html`

	_, err := s.db.Exec(query, tooltip.Slug, tooltip.SourcePath, tooltip.ModUnix, tooltip.Size, tooltip.Title, tooltip.HTML)
	return err
}

func (s *SQLiteStore) migrate() error {
	const schema = `
		CREATE TABLE IF NOT EXISTS rendered_posts (
			slug TEXT PRIMARY KEY,
			source_path TEXT NOT NULL,
			mod_unix INTEGER NOT NULL,
			size_bytes INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			author TEXT NOT NULL,
			published_at TEXT NOT NULL,
			tags_json TEXT NOT NULL,
			html TEXT NOT NULL,
			excerpt TEXT NOT NULL,
			read_time INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS rendered_tooltips (
			slug TEXT PRIMARY KEY,
			source_path TEXT NOT NULL,
			mod_unix INTEGER NOT NULL,
			size_bytes INTEGER NOT NULL,
			title TEXT NOT NULL,
			html TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS rendered_pages (
			slug TEXT PRIMARY KEY,
			source_path TEXT NOT NULL,
			mod_unix INTEGER NOT NULL,
			size_bytes INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			html TEXT NOT NULL,
			excerpt TEXT NOT NULL,
			read_time INTEGER NOT NULL
		);`

	_, err := s.db.Exec(schema)
	return err
}
