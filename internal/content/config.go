package content

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Options struct {
	ConfigPath     string
	ContentDir     string
	StaticDir      string
	TemplateDir    string
	OutputDir      string
	DataDir        string
	IncludeDrafts  bool
	CompressOutput bool
}

type Config struct {
	BaseURL          string           `toml:"base_url"`
	Title            string           `toml:"title"`
	Description      string           `toml:"description"`
	DefaultLanguage  string           `toml:"default_language"`
	OutputDir        string           `toml:"output_dir"`
	Author           string           `toml:"author"`
	BuildSearchIndex bool             `toml:"build_search_index"`
	GenerateFeeds    bool             `toml:"generate_feeds"`
	Taxonomies       []TaxonomyConfig `toml:"taxonomies"`
	Extra            map[string]any   `toml:"extra"`
	Search           SearchConfig     `toml:"search"`
}

type TaxonomyConfig struct {
	Name   string `toml:"name"`
	Feed   bool   `toml:"feed"`
	Render *bool  `toml:"render"`
}

type SearchConfig struct {
	IndexFormat string `toml:"index_format"`
}

func (c TaxonomyConfig) ShouldRender() bool {
	if c.Render == nil {
		return true
	}
	return *c.Render
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return cfg, err
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080"
	}
	if cfg.DefaultLanguage == "" {
		cfg.DefaultLanguage = "en"
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = filepath.ToSlash("public")
	}
	if cfg.Extra == nil {
		cfg.Extra = map[string]any{}
	}
	return cfg, nil
}
