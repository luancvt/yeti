package yang

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	gyang "github.com/openconfig/goyang/pkg/yang"
)

// Node represents a YANG tree node with the metadata we care about.
type Node struct {
	Name        string
	Path        string
	Kind        string // "container", "list", "leaf", "leaf-list", "module"
	Description string
	Config      bool   // true = config (rw), false = state (ro)
	Mandatory   bool
	Default     string
	Units       string
	Key         string // list key, empty for non-lists
	Type        *TypeInfo
	Children    []*Node
}

// TypeInfo holds resolved type information for leaf nodes.
type TypeInfo struct {
	Name           string
	Kind           string // goyang TypeKind as string
	Range          string
	Length         string
	Pattern        []string
	EnumValues     []string
	Default        string
	FractionDigits int
}

// CollectionTree holds the parsed YANG modules for one collection.
type CollectionTree struct {
	modules map[string]*gyang.Entry
}

// ParseCollection parses all .yang files in the given path within the filesystem
// and returns the merged tree as a queryable CollectionTree.
func ParseCollection(fsys fs.FS, path string) (*CollectionTree, error) {
	ms := gyang.NewModules()

	var yangFiles []string
	err := fs.WalkDir(fsys, path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(p, ".yang") {
			yangFiles = append(yangFiles, p)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", path, err)
	}

	if len(yangFiles) == 0 {
		return nil, fmt.Errorf("no .yang files found in %s", path)
	}

	for _, f := range yangFiles {
		data, err := fs.ReadFile(fsys, f)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", f, err)
		}
		if err := ms.Parse(string(data), f); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, err)
		}
	}

	errs := ms.Process()
	if len(errs) > 0 {
		return nil, fmt.Errorf("processing modules: %v", errs[0])
	}

	tree := &CollectionTree{
		modules: make(map[string]*gyang.Entry),
	}

	for _, mod := range ms.Modules {
		if mod == nil {
			continue
		}
		entry := gyang.ToEntry(mod)
		if entry == nil {
			continue
		}
		tree.modules[entry.Name] = entry
	}

	return tree, nil
}

// Children returns the top-level nodes across all modules in the collection.
func (ct *CollectionTree) Children() []Node {
	var nodes []Node
	for _, entry := range ct.modules {
		if entry.Dir == nil {
			continue
		}
		for _, child := range sortedEntries(entry.Dir) {
			nodes = append(nodes, entryToNode(child, "/"+child.Name))
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})
	return nodes
}

// GetNode returns the node at the given slash-separated path.
func (ct *CollectionTree) GetNode(path string) (*Node, error) {
	entry, err := ct.getEntry(path)
	if err != nil {
		return nil, err
	}
	node := entryToNode(entry, path)
	return &node, nil
}

// GetChildren returns the direct children of the node at the given path.
func (ct *CollectionTree) GetChildren(path string) ([]Node, error) {
	entry, err := ct.getEntry(path)
	if err != nil {
		return nil, err
	}

	if entry.Dir == nil {
		return nil, nil
	}

	var children []Node
	for _, child := range sortedEntries(entry.Dir) {
		childPath := path + "/" + child.Name
		children = append(children, entryToNode(child, childPath))
	}
	return children, nil
}

func (ct *CollectionTree) getEntry(path string) (*gyang.Entry, error) {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return nil, fmt.Errorf("empty path")
	}

	var current *gyang.Entry
	for _, mod := range ct.modules {
		if mod.Dir == nil {
			continue
		}
		if child, ok := mod.Dir[parts[0]]; ok {
			current = child
			break
		}
	}
	if current == nil {
		return nil, fmt.Errorf("node not found: %s", parts[0])
	}

	for _, part := range parts[1:] {
		if current.Dir == nil {
			return nil, fmt.Errorf("node %q has no children", current.Name)
		}
		child, ok := current.Dir[part]
		if !ok {
			return nil, fmt.Errorf("node not found: %s", part)
		}
		current = child
	}

	return current, nil
}

func entryToNode(e *gyang.Entry, path string) Node {
	n := Node{
		Name:        e.Name,
		Path:        path,
		Description: e.Description,
		Config:      e.Config != gyang.TSFalse,
		Mandatory:   e.Mandatory == gyang.TSTrue,
		Units:       e.Units,
	}

	if len(e.Default) > 0 {
		n.Default = e.Default[0]
	}

	switch {
	case e.IsList():
		n.Kind = "list"
		n.Key = e.Key
	case e.IsLeaf():
		n.Kind = "leaf"
	case e.IsLeafList():
		n.Kind = "leaf-list"
	case e.IsContainer():
		n.Kind = "container"
	case e.Dir != nil:
		n.Kind = "container"
	default:
		n.Kind = "leaf"
	}

	if e.Type != nil {
		n.Type = extractTypeInfo(e.Type)
	}

	return n
}

func extractTypeInfo(t *gyang.YangType) *TypeInfo {
	ti := &TypeInfo{
		Name:           t.Name,
		Kind:           t.Kind.String(),
		FractionDigits: t.FractionDigits,
	}

	if t.Range != nil {
		ti.Range = t.Range.String()
	}
	if t.Length != nil {
		ti.Length = t.Length.String()
	}
	ti.Pattern = append(ti.Pattern, t.Pattern...)
	if t.Enum != nil {
		ti.EnumValues = t.Enum.Names()
	}

	return ti
}

func sortedEntries(dir map[string]*gyang.Entry) []*gyang.Entry {
	names := make([]string, 0, len(dir))
	for k := range dir {
		names = append(names, k)
	}
	sort.Strings(names)

	entries := make([]*gyang.Entry, 0, len(names))
	for _, name := range names {
		entries = append(entries, dir[name])
	}
	return entries
}
