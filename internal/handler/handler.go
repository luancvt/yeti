package handler

import (
	"net/http"
	"sort"

	"github.com/terjelafton/yeti/internal/view"
	"github.com/terjelafton/yeti/internal/yang"
)

type Handler struct {
	trees        map[string]*yang.CollectionTree
	displayNames map[string]string
}

func New(trees map[string]*yang.CollectionTree, displayNames map[string]string) *Handler {
	return &Handler{trees: trees, displayNames: displayNames}
}

// Collections returns sorted collection info with display names.
func (h *Handler) Collections() []view.CollectionInfo {
	infos := make([]view.CollectionInfo, 0, len(h.trees))
	for name := range h.trees {
		display := h.displayNames[name]
		if display == "" {
			display = name
		}
		infos = append(infos, view.CollectionInfo{Name: name, Display: display})
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})
	return infos
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	view.Index(h.Collections(), "", "").Render(r.Context(), w)
}

func (h *Handler) Browse(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	module := r.PathValue("module")

	tree, ok := h.trees[collection]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Verify the module exists
	if _, err := tree.ModuleChildren(module); err != nil {
		http.NotFound(w, r)
		return
	}

	view.Index(h.Collections(), collection, module).Render(r.Context(), w)
}

func (h *Handler) Models(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	if collection == "" {
		collection = r.FormValue("collection")
	}
	tree, ok := h.trees[collection]
	if !ok {
		http.NotFound(w, r)
		return
	}
	view.ModelPicker(tree.ModuleNames(), collection).Render(r.Context(), w)
}

func (h *Handler) Tree(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	module := r.PathValue("module")
	tree, ok := h.trees[collection]
	if !ok {
		http.NotFound(w, r)
		return
	}

	pathStr := r.PathValue("path")

	var nodes []yang.Node
	var err error
	if pathStr == "" {
		// No path — return module's top-level children
		nodes, err = tree.ModuleChildren(module)
	} else {
		nodes, err = tree.GetChildren(module, "/"+pathStr)
	}

	if err != nil {
		http.NotFound(w, r)
		return
	}

	view.TreeNodeList(nodes, collection, module).Render(r.Context(), w)
}

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	module := r.PathValue("module")
	tree, ok := h.trees[collection]
	if !ok {
		http.NotFound(w, r)
		return
	}

	path := "/" + r.PathValue("path")
	node, err := tree.GetNode(module, path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	view.Detail(node).Render(r.Context(), w)
}

func (h *Handler) EmptyTree(w http.ResponseWriter, r *http.Request) {
	view.EmptyTree().Render(r.Context(), w)
}

func (h *Handler) EmptyDetail(w http.ResponseWriter, r *http.Request) {
	view.EmptyDetail().Render(r.Context(), w)
}
