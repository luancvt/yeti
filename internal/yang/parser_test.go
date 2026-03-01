package yang_test

import (
	"os"
	"testing"

	yeti "github.com/terjelafton/yeti/internal/yang"
)

func TestParseCollection(t *testing.T) {
	modelsFS := os.DirFS("../../models")

	tree, err := yeti.ParseCollection(modelsFS, "test")
	if err != nil {
		t.Fatalf("ParseCollection failed: %v", err)
	}

	if tree == nil {
		t.Fatal("expected non-nil tree")
	}

	// The test module should have a top-level "interfaces" container
	children := tree.Children()
	found := false
	for _, child := range children {
		if child.Name == "interfaces" {
			found = true
			break
		}
	}
	if !found {
		names := make([]string, len(children))
		for i, c := range children {
			names[i] = c.Name
		}
		t.Errorf("expected 'interfaces' in top-level children, got %v", names)
	}
}
