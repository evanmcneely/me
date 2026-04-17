package cache

import (
	"sync"
	"time"
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

type Store struct {
	mu       sync.RWMutex
	posts    map[string]CachedPost
	pages    map[string]CachedPage
	tooltips map[string]CachedTooltip
}

func NewStore() *Store {
	return &Store{
		posts:    make(map[string]CachedPost),
		pages:    make(map[string]CachedPage),
		tooltips: make(map[string]CachedTooltip),
	}
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) GetPost(slug string, modUnix, size int64) (*CachedPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[slug]
	if !ok || post.ModUnix != modUnix || post.Size != size {
		return nil, nil
	}

	copy := post
	return &copy, nil
}

func (s *Store) UpsertPost(post CachedPost) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.posts[post.Slug] = post
	return nil
}

func (s *Store) GetTooltip(slug string, modUnix, size int64) (*CachedTooltip, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tooltip, ok := s.tooltips[slug]
	if !ok || tooltip.ModUnix != modUnix || tooltip.Size != size {
		return nil, nil
	}

	copy := tooltip
	return &copy, nil
}

func (s *Store) UpsertTooltip(tooltip CachedTooltip) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tooltips[tooltip.Slug] = tooltip
	return nil
}

func (s *Store) GetPage(slug string, modUnix, size int64) (*CachedPage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	page, ok := s.pages[slug]
	if !ok || page.ModUnix != modUnix || page.Size != size {
		return nil, nil
	}

	copy := page
	return &copy, nil
}

func (s *Store) UpsertPage(page CachedPage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pages[page.Slug] = page
	return nil
}
