package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"evanmcneely/internal/blog"
	"evanmcneely/internal/cache"
	"evanmcneely/internal/render"
	"evanmcneely/internal/web"
)

func main() {
	addr := getenv("ADDR", ":8080")
	contentDir := getenv("CONTENT_DIR", "content")

	store := cache.NewStore()
	defer store.Close()

	renderer := render.NewMarkdownRenderer()
	service := blog.NewService(contentDir, store, renderer)
	app, err := web.NewApp(service, renderer.HighlightCSS())
	if err != nil {
		log.Fatalf("build app: %v", err)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           app.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("blog listening on http://localhost%s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
