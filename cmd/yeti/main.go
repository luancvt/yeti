package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/luancvt/yeti/internal/config"
	"github.com/luancvt/yeti/internal/handler"
	"github.com/luancvt/yeti/internal/yang"
	"github.com/luancvt/yeti/static"
)

func main() {
	configPath := envOr("YETI_CONFIG", "config.yaml")
	modelsDir := envOr("YETI_MODELS_DIR", "models")

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	modelsFS := os.DirFS(modelsDir)
	trees := make(map[string]*yang.CollectionTree)
	displayNames := make(map[string]string)

	for _, c := range cfg.Collections {
		log.Printf("Parsing collection %s...", c.Name)
		tree, err := yang.ParseCollection(modelsFS, c.Path)
		if err != nil {
			log.Printf("  skipping %s: %v", c.Name, err)
			continue
		}
		trees[c.Name] = tree
		displayNames[c.Name] = c.Display
		log.Printf("  %d modules", len(tree.ModuleNames()))
	}

	h := handler.New(trees, displayNames)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("GET /{collection}/{module}", h.Browse)
	mux.HandleFunc("GET /models", h.Models)
	mux.HandleFunc("GET /models/{collection}", h.Models)
	mux.HandleFunc("GET /tree/{collection}/{module}", h.Tree)
	mux.HandleFunc("GET /tree/{collection}/{module}/{path...}", h.Tree)
	mux.HandleFunc("GET /detail/{collection}/{module}/{path...}", h.Detail)
	mux.HandleFunc("GET /empty/tree", h.EmptyTree)
	mux.HandleFunc("GET /empty/detail", h.EmptyDetail)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	staticHandler := http.StripPrefix("/static/", http.FileServerFS(static.FS))
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			staticHandler.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      httpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Yeti running on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}

func loadConfig(path string) (*config.Config, error) {
	f, err := os.Open(path) //nolint:gosec // G304: config path from env/default
	if err != nil {
		return nil, fmt.Errorf("opening config %s: %w", path, err)
	}

	cfg, err := config.Load(f)
	closeErr := f.Close()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("closing config: %w", closeErr)
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
