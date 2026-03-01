package config_test

import (
	"strings"
	"testing"

	"github.com/terjelafton/yeti/internal/config"
)

func TestLoadConfig(t *testing.T) {
	input := `
collections:
  - name: test-collection
    display: "Test Collection"
    path: test
`
	cfg, err := config.Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Collections) != 1 {
		t.Fatalf("expected 1 collection, got %d", len(cfg.Collections))
	}

	c := cfg.Collections[0]
	if c.Name != "test-collection" {
		t.Errorf("expected name 'test-collection', got %q", c.Name)
	}
	if c.Display != "Test Collection" {
		t.Errorf("expected display 'Test Collection', got %q", c.Display)
	}
	if c.Path != "test" {
		t.Errorf("expected path 'test', got %q", c.Path)
	}
}

func TestLoadConfigValidation(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing name", `collections: [{ display: "X", path: "x" }]`},
		{"missing display", `collections: [{ name: "x", path: "x" }]`},
		{"missing path", `collections: [{ name: "x", display: "X" }]`},
		{"empty collections", `collections: []`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := config.Load(strings.NewReader(tt.input))
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
