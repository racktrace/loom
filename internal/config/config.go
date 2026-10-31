package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds every user-adjustable setting.
type Config struct {
	Editor  string `toml:"editor"`
	DataDir string `toml:"data_dir"`

	Scratch struct {
		DateFormat string `toml:"date_format"`
	} `toml:"scratch"`

	Note struct {
		DefaultTags []string `toml:"default_tags"`
	} `toml:"note"`
}

// Default returns the built-in values used when the config file or a field is missing.
func Default() Config {
	var c Config
	c.DataDir = defaultDataDir()
	c.Scratch.DateFormat = "%Y-%m-%d %H:%M"
	c.Note.DefaultTags = []string{}
	return c
}

// Path returns the config file location, e.g. ~/.config/loom/config.toml.
// os.UserConfigDir already respects $XDG_CONFIG_HOME on Linux.
func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "loom", "config.toml"), nil
}

// Load applies defaults, then config.toml, then env vars. Later sources win.
func Load() (Config, error) {
	cfg := Default()

	p, err := Path()
	if err != nil {
		return cfg, err
	}
	if _, err := toml.DecodeFile(p, &cfg); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return cfg, err // file exists but is invalid: report it
	}

	if v := os.Getenv("LOOM_DIR"); v != "" {
		cfg.DataDir = v
	}
	if v := os.Getenv("LOOM_EDITOR"); v != "" {
		cfg.Editor = v
	}

	cfg.DataDir = expandHome(cfg.DataDir)
	return cfg, nil
}

func (c Config) EditorCommand() []string {
	for _, e := range []string{c.Editor, os.Getenv("EDITOR")} {
		if parts := strings.Fields(e); len(parts) > 0 {
			return parts
		}
	}
	return []string{"nano"}
}

// defaultDataDir follows XDG: $XDG_DATA_HOME/loom, else ~/.local/share/loom.
func defaultDataDir() string {
	if x := os.Getenv("XDG_DATA_HOME"); filepath.IsAbs(x) {
		return filepath.Join(x, "loom")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "loom")
}

// expandHome turns "~/notes" into "/home/user/notes".
func expandHome(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	return p
}
