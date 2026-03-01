package handler

import (
	"net/http"

	"github.com/terjelafton/yeti/internal/view"
	"github.com/terjelafton/yeti/internal/yang"
)

type Handler struct {
	trees map[string]*yang.CollectionTree
}

func New(trees map[string]*yang.CollectionTree) *Handler {
	return &Handler{trees: trees}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	tree := h.trees["test"]
	nodes := tree.Children()
	view.Index(nodes, "test").Render(r.Context(), w)
}

func (h *Handler) TreeChildren(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	tree, ok := h.trees[collection]
	if !ok {
		http.NotFound(w, r)
		return
	}

	path := "/" + r.PathValue("path")
	children, err := tree.GetChildren(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	view.TreeNodeList(children, collection).Render(r.Context(), w)
}

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	tree, ok := h.trees[collection]
	if !ok {
		http.NotFound(w, r)
		return
	}

	path := "/" + r.PathValue("path")
	node, err := tree.GetNode(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	view.Detail(node).Render(r.Context(), w)
}
