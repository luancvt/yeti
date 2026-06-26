package handler

import (
	"cmp"
	"log"
	"net/http"
	"slices"

	"github.com/a-h/templ"

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
	slices.SortFunc(infos, func(a, b view.CollectionInfo) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return infos
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	render(w, r, view.Index(h.Collections(), "", ""))
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

	render(w, r, view.Index(h.Collections(), collection, module))
}

func (h *Handler) Models(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	collection := r.PathValue("collection")
	if collection == "" {
		collection = r.FormValue("collection")
	}
	tree, ok := h.trees[collection]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("reset") == "true" {
		render(w, r, view.ModelPickerWithReset(tree.ModuleNames(), collection))
	} else {
		render(w, r, view.ModelPicker(tree.ModuleNames(), collection))
	}
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

	if pathStr == "" {
		render(w, r, view.TreeNodeListWithReset(nodes, collection, module))
	} else {
		render(w, r, view.TreeNodeList(nodes, collection, module))
	}
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

	render(w, r, view.Detail(node))
}

func (*Handler) EmptyTree(w http.ResponseWriter, r *http.Request) {
	render(w, r, view.EmptyTree())
}

func (*Handler) EmptyDetail(w http.ResponseWriter, r *http.Request) {
	render(w, r, view.EmptyDetail())
}

func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("render: %v", err)
	}
}
