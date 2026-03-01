package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/terjelafton/yeti/internal/handler"
	"github.com/terjelafton/yeti/internal/yang"
)

func main() {
	modelsFS := os.DirFS("models")

	// Auto-discover collections from models/ directory
	entries, err := os.ReadDir("models")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading models directory: %v\n", err)
		os.Exit(1)
	}

	trees := make(map[string]*yang.CollectionTree)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		log.Printf("Parsing collection %s...", name)
		tree, err := yang.ParseCollection(modelsFS, name)
		if err != nil {
			log.Printf("  skipping %s: %v", name, err)
			continue
		}
		trees[name] = tree
		log.Printf("  %d modules", len(tree.ModuleNames()))
	}

	h := handler.New(trees)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("GET /{collection}/{module}", h.Browse)
	mux.HandleFunc("GET /models/{collection}", h.Models)
	mux.HandleFunc("GET /tree/{collection}/{module}", h.Tree)
	mux.HandleFunc("GET /tree/{collection}/{module}/{path...}", h.Tree)
	mux.HandleFunc("GET /detail/{collection}/{module}/{path...}", h.Detail)

	log.Println("Yeti running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
