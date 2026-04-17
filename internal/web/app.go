package web

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
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

type errorData struct {
	siteData
	StatusCode int
	Heading    string
	Message    string
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
	return a.recoverMiddleware(mux)
}

func (a *App) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.handlePage(w, r)
		return
	}

	posts, err := a.service.ListPosts(r.Context())
	if err != nil {
		a.serverError(w, r, "Unable to load posts right now.")
		return
	}

	pages, err := a.service.ListPages(r.Context())
	if err != nil {
		a.serverError(w, r, "Unable to load pages right now.")
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
		a.notFound(w, r)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	if slug == "" || strings.Contains(slug, "/") {
		a.notFound(w, r)
		return
	}

	post, err := a.service.Post(r.Context(), slug)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			a.notFound(w, r)
			return
		}
		a.serverError(w, r, "Unable to load that post right now.")
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
		a.notFound(w, r)
		return
	}

	page, err := a.service.Page(r.Context(), slug)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			a.notFound(w, r)
			return
		}
		a.serverError(w, r, "Unable to load that page right now.")
		return
	}

	a.render(w, "page.html", pageData{
		siteData: a.site(page.Title, r.URL.Path),
		Page:     page,
	})
}

func (a *App) handleTooltip(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/tooltips/") {
		a.renderTooltipError(w, http.StatusNotFound, "Tooltip not found.")
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/tooltips/")
	if slug == "" || strings.Contains(slug, "/") {
		a.renderTooltipError(w, http.StatusNotFound, "Tooltip not found.")
		return
	}

	tooltip, err := a.service.Tooltip(r.Context(), slug)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			a.renderTooltipError(w, http.StatusNotFound, "Tooltip not found.")
			return
		}
		a.renderTooltipError(w, http.StatusInternalServerError, "Tooltip unavailable right now.")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, `<div class="tooltip-card">%s</div>`, tooltip.HTML)
}

func (a *App) render(w http.ResponseWriter, name string, data any) {
	a.renderStatus(w, http.StatusOK, name, data)
}

func (a *App) renderStatus(w http.ResponseWriter, status int, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := a.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (a *App) notFound(w http.ResponseWriter, r *http.Request) {
	a.renderStatus(w, http.StatusNotFound, "error.html", errorData{
		siteData:   a.site("Not Found", r.URL.Path),
		StatusCode: http.StatusNotFound,
		Heading:    "Not Found",
		Message:    "That page does not exist, or it may have moved.",
	})
}

func (a *App) serverError(w http.ResponseWriter, r *http.Request, message string) {
	a.renderStatus(w, http.StatusInternalServerError, "error.html", errorData{
		siteData:   a.site("Server Error", r.URL.Path),
		StatusCode: http.StatusInternalServerError,
		Heading:    "Oops!",
		Message:    message,
	})
}

func (a *App) renderTooltipError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, `<div class="tooltip-card"><p>%s</p></div>`, template.HTMLEscapeString(message))
}

func (a *App) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic serving %s: %v", r.URL.Path, recovered)
				a.serverError(w, r, "Something went wrong on our side.")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *App) site(pageTitle, path string) siteData {
	return siteData{
		SiteTitle: "Evan McNeely",
		ChromaCSS: a.chromaCSS,
		PageTitle: pageTitle,
		Path:      path,
	}
}
