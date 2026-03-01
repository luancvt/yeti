package config

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Collection struct {
	Name    string `yaml:"name"`
	Display string `yaml:"display"`
	Path    string `yaml:"path"`
}

type Config struct {
	Collections []Collection `yaml:"collections"`
}

func Load(r io.Reader) (*Config, error) {
	var cfg Config
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if len(cfg.Collections) == 0 {
		return nil, fmt.Errorf("config must have at least one collection")
	}

	for i, c := range cfg.Collections {
		if c.Name == "" {
			return nil, fmt.Errorf("collection %d: missing name", i)
		}
		if c.Display == "" {
			return nil, fmt.Errorf("collection %d (%s): missing display", i, c.Name)
		}
		if c.Path == "" {
			return nil, fmt.Errorf("collection %d (%s): missing path", i, c.Name)
		}
	}

	return &cfg, nil
}
