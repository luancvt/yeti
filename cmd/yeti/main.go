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

	tree, err := yang.ParseCollection(modelsFS, "test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing models: %v\n", err)
		os.Exit(1)
	}

	trees := map[string]*yang.CollectionTree{"test": tree}
	h := handler.New(trees)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("GET /tree/{collection}/{path...}", h.TreeChildren)
	mux.HandleFunc("GET /detail/{collection}/{path...}", h.Detail)

	log.Println("Yeti running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
