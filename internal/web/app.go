package web

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"evanmcneely/internal/blog"
)

//go:embed templates/*.html static/*
var assets embed.FS

type App struct {
	service     *blog.Service
	templates   *template.Template
	chromaCSS   template.CSS
	staticFiles http.Handler
}

type siteData struct {
	SiteTitle string
	ChromaCSS template.CSS
	PageTitle string
	Path      string
}

type homeData struct {
	siteData
	Posts []blog.Post
	Pages []blog.Page
}

type postData struct {
	siteData
	Post blog.Post
}

type pageData struct {
	siteData
	Page blog.Page
}

func NewApp(service *blog.Service, chromaCSS string) (*App, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return "Draft"
			}
			return t.Format("02 Jan 2006")
		},
		"isoDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02")
		},
	}).ParseFS(assets, "templates/*.html")
	if err != nil {
		return nil, err
	}

	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		return nil, err
	}

	return &App{
		service:     service,
		templates:   tmpl,
		chromaCSS:   template.CSS(chromaCSS),
		staticFiles: http.FileServer(http.FS(staticFS)),
	}, nil
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", a.staticFiles))
	mux.HandleFunc("/", a.handleHome)
	mux.HandleFunc("/posts/", a.handlePost)
	mux.HandleFunc("/tooltips/", a.handleTooltip)
	return mux
}

func (a *App) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.handlePage(w, r)
		return
	}

	posts, err := a.service.ListPosts(r.Context())
	if err != nil {
		http.Error(w, "unable to load posts", http.StatusInternalServerError)
		return
	}

	pages, err := a.service.ListPages(r.Context())
	if err != nil {
		http.Error(w, "unable to load pages", http.StatusInternalServerError)
		return
	}

	a.render(w, "home.html", homeData{
		siteData: a.site("Home", r.URL.Path),
		Posts:    posts,
		Pages:    pages,
	})
}

func (a *App) handlePost(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/posts/") {
		http.NotFound(w, r)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(w, r)
		return
	}

	post, err := a.service.Post(r.Context(), slug)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "unable to load post", http.StatusInternalServerError)
		return
	}

	a.render(w, "post.html", postData{
		siteData: a.site(post.Title, r.URL.Path),
		Post:     post,
	})
}

func (a *App) handlePage(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(w, r)
		return
	}

	page, err := a.service.Page(r.Context(), slug)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "unable to load page", http.StatusInternalServerError)
		return
	}

	a.render(w, "page.html", pageData{
		siteData: a.site(page.Title, r.URL.Path),
		Page:     page,
	})
}

func (a *App) handleTooltip(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/tooltips/") {
		http.NotFound(w, r)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/tooltips/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(w, r)
		return
	}

	tooltip, err := a.service.Tooltip(r.Context(), slug)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "unable to load tooltip", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<div class="tooltip-card">%s</div>`, tooltip.HTML)
}

func (a *App) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (a *App) site(pageTitle, path string) siteData {
	return siteData{
		SiteTitle: "Evan McNeely",
		ChromaCSS: a.chromaCSS,
		PageTitle: pageTitle,
		Path:      path,
	}
}
