package blog

import (
	"html/template"
	"time"
)

type Post struct {
	Slug        string
	Title       string
	Description string
	Author      string
	PublishedAt time.Time
	Tags        []string
	HTML        template.HTML
	Excerpt     string
	ReadTime    int
}

type Tooltip struct {
	Slug  string
	Title string
	HTML  template.HTML
}
