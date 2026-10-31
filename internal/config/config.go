package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// defaultTemplate is written to the config path on first run.
// editor and data_dir are commented out on purpose, so $EDITOR and
// $XDG_DATA_HOME keep working until the user decides to override them.
const defaultTemplate = `# loom configuration
# Remove the leading "#" to override a setting.

# Editor for loom. If unset: $EDITOR, then nano.
# editor = "nano"

# Where scratch and notes are stored. If unset: $XDG_DATA_HOME/loom
# or ~/.local/share/loom.
# data_dir = "~/.local/share/loom"

[scratch]
# Date format for scratch session headers.
# Supported: %Y %y %m %d %H %M %S %b %a %%
date_format = "%Y-%m-%d %H:%M"

[note]
# Tags added to every new note, e.g. ["inbox"].
default_tags = []
`

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
// On first run, a commented default config file is created.
func Load() (Config, error) {
	cfg := Default()

	p, err := Path()
	if err != nil {
		return cfg, err
	}
	if _, err := toml.DecodeFile(p, &cfg); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return cfg, err // file exists but is invalid: report it
		}
		_ = writeDefault(p) // first run: best effort, never block loom
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

// EditorCommand resolves the editor: config/LOOM_EDITOR, then $EDITOR, then nano.
// The result is split into fields so values like "code -w" work.
func (c Config) EditorCommand() []string {
	for _, e := range []string{c.Editor, os.Getenv("EDITOR")} {
		if parts := strings.Fields(e); len(parts) > 0 {
			return parts
		}
	}
	return []string{"nano"}
}

// writeDefault creates the config file with defaultTemplate.
// O_EXCL guarantees an existing config is never overwritten.
func writeDefault(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(defaultTemplate); err != nil {
		f.Close()
		return err
	}
	return f.Close()
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
