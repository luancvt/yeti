package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/terjelafton/yeti/internal/config"
	"github.com/terjelafton/yeti/internal/handler"
	"github.com/terjelafton/yeti/internal/yang"
	"github.com/terjelafton/yeti/static"
)

func main() {
	configPath := envOr("YETI_CONFIG", "config.yaml")
	modelsDir := envOr("YETI_MODELS_DIR", "models")

	f, err := os.Open(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening config %s: %v\n", configPath, err)
		os.Exit(1)
	}
	defer f.Close()

	cfg, err := config.Load(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
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

	staticHandler := http.StripPrefix("/static/", http.FileServerFS(static.FS))
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			staticHandler.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	log.Println("Yeti running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", httpHandler))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
