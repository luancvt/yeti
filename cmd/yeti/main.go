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

	collections := []struct{ name, path string }{
		{"test", "test"},
		{"xr-7112", "xr-7112"},
	}

	trees := make(map[string]*yang.CollectionTree)
	for _, c := range collections {
		log.Printf("Parsing collection %s...", c.name)
		tree, err := yang.ParseCollection(modelsFS, c.path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing %s: %v\n", c.name, err)
			os.Exit(1)
		}
		trees[c.name] = tree
		log.Printf("  %d top-level nodes", len(tree.Children()))
	}

	h := handler.New(trees)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("GET /tree/{collection}/{path...}", h.TreeChildren)
	mux.HandleFunc("GET /detail/{collection}/{path...}", h.Detail)

	log.Println("Yeti running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
