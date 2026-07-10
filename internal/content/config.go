package content

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Options struct {
	ConfigPath     string
	ContentDir     string
	StaticDir      string
	TemplateDir    string
	OutputDir      string
	IncludeDrafts  bool
	CompressOutput bool
}

type Config struct {
	BaseURL          string           `toml:"base_url"`
	Title            string           `toml:"title"`
	Description      string           `toml:"description"`
	DefaultLanguage  string           `toml:"default_language"`
	Author           string           `toml:"author"`
	BuildSearchIndex bool             `toml:"build_search_index"`
	GenerateFeeds    bool             `toml:"generate_feeds"`
	Taxonomies       []TaxonomyConfig `toml:"taxonomies"`
	Extra            map[string]any   `toml:"extra"`
}

type TaxonomyConfig struct {
	Name   string `toml:"name"`
	Feed   bool   `toml:"feed"`
	Render *bool  `toml:"render"`
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
	if cfg.Extra == nil {
		cfg.Extra = map[string]any{}
	}
	return cfg, nil
}
