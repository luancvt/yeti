package yang

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	gyang "github.com/openconfig/goyang/pkg/yang"
)

// NodeKind represents the YANG schema node type.
type NodeKind string

const (
	KindContainer NodeKind = "container"
	KindList      NodeKind = "list"
	KindLeaf      NodeKind = "leaf"
	KindLeafList  NodeKind = "leaf-list"
	KindModule    NodeKind = "module"
)

// Node represents a YANG tree node with the metadata we care about.
type Node struct {
	Name        string
	Path        string
	Kind        NodeKind
	Description string
	Config      bool // true = config (rw), false = state (ro)
	Mandatory   bool
	Default     string
	Units       string
	Key         string // list key, empty for non-lists
	Type        *TypeInfo
	Extra       map[string][]string // YANG metadata (presence, when, must …)
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
	entries, err := fs.ReadDir(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".yang") {
			yangFiles = append(yangFiles, path+"/"+e.Name())
		}
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
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			continue
		}
	}

	// Process resolves imports, augments, etc. Errors are expected with
	// real-world vendor models (unresolved augments, missing modules).
	// We log them but continue — the tree is still usable for browsing.
	errs := ms.Process()
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
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

// ModuleNames returns the sorted list of module names in this collection.
func (ct *CollectionTree) ModuleNames() []string {
	names := make([]string, 0, len(ct.modules))
	for name := range ct.modules {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ModuleChildren returns the top-level nodes for a specific module.
func (ct *CollectionTree) ModuleChildren(module string) ([]Node, error) {
	entry, ok := ct.modules[module]
	if !ok {
		return nil, fmt.Errorf("module not found: %s", module)
	}
	if entry.Dir == nil {
		return nil, nil
	}
	var nodes []Node
	for _, child := range sortedEntries(entry.Dir) {
		nodes = append(nodes, entryToNode(child, "/"+child.Name))
	}
	return nodes, nil
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

// GetNode returns the node at the given slash-separated path, scoped to a module.
func (ct *CollectionTree) GetNode(module, path string) (*Node, error) {
	entry, err := ct.getEntry(module, path)
	if err != nil {
		return nil, err
	}
	node := entryToNode(entry, path)
	return &node, nil
}

// GetChildren returns the direct children of the node at the given path, scoped to a module.
func (ct *CollectionTree) GetChildren(module, path string) ([]Node, error) {
	entry, err := ct.getEntry(module, path)
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

func (ct *CollectionTree) getEntry(module, path string) (*gyang.Entry, error) {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return nil, errors.New("empty path")
	}

	// Try the specified module first
	if mod, ok := ct.modules[module]; ok && mod.Dir != nil {
		if current, ok := mod.Dir[parts[0]]; ok {
			if entry := walkPath(current, parts[1:]); entry != nil {
				return entry, nil
			}
		}
	}

	// Fall back to searching all modules — handles augmented nodes
	// from other modules that weren't merged during Process().
	for _, mod := range ct.modules {
		if mod.Dir == nil {
			continue
		}
		current, ok := mod.Dir[parts[0]]
		if !ok {
			continue
		}
		if entry := walkPath(current, parts[1:]); entry != nil {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("node not found: %s", path)
}

func walkPath(entry *gyang.Entry, parts []string) *gyang.Entry {
	current := entry
	for _, part := range parts {
		if current.Dir == nil {
			return nil
		}
		child, ok := current.Dir[part]
		if !ok {
			return nil
		}
		current = child
	}
	return current
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
		n.Kind = KindList
		n.Key = e.Key
	case e.IsLeaf():
		n.Kind = KindLeaf
	case e.IsLeafList():
		n.Kind = KindLeafList
	case e.IsContainer():
		n.Kind = KindContainer
	case e.Dir != nil:
		n.Kind = KindContainer
	default:
		n.Kind = KindLeaf
	}

	if e.Type != nil {
		n.Type = extractTypeInfo(e.Type)
	}

	if len(e.Extra) > 0 {
		extra := make(map[string][]string, len(e.Extra))
		for key, vals := range e.Extra {
			strs := make([]string, 0, len(vals))
			for _, v := range vals {
				if yv, ok := v.(*gyang.Value); ok {
					strs = append(strs, yv.Name)
				} else {
					strs = append(strs, fmt.Sprintf("%v", v))
				}
			}
			if len(strs) > 0 {
				extra[key] = strs
			}
		}
		if len(extra) > 0 {
			n.Extra = extra
		}
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
