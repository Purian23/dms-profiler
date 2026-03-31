package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const DefaultPath = "~/.config/dms-profiler/config.toml"

// File is the TOML config root.
type File struct {
	Browser Browser `toml:"browser"`
	Rules   []Rule  `toml:"rules"`
	Default Default `toml:"default"`
}

type Browser struct {
	Command     string `toml:"command"`
	UserDataDir string `toml:"user_data_dir"`
}

// Rule is one [[rules]] entry: use match or matches (non-empty matches ignores match).
type Rule struct {
	Match   string   `toml:"match"`
	Matches []string `toml:"matches"`
	Profile string   `toml:"profile"`
}

type FlatRule struct {
	Match   string
	Profile string
}

// FlattenRules turns each [[rules]] block into ordered prefix rows (first match wins).
func FlattenRules(rules []Rule) []FlatRule {
	out := make([]FlatRule, 0)
	for _, r := range rules {
		if len(r.Matches) > 0 {
			for _, m := range r.Matches {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				out = append(out, FlatRule{Match: m, Profile: r.Profile})
			}
			continue
		}
		if strings.TrimSpace(r.Match) != "" {
			out = append(out, FlatRule{Match: strings.TrimSpace(r.Match), Profile: r.Profile})
		}
	}
	return out
}

type Default struct {
	Profile string `toml:"profile"`
}

// Load reads config from path (~ and env vars expanded).
func Load(path string) (*File, error) {
	expanded, err := ExpandPath(path)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(expanded)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var f File
	if err := toml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if f.Browser.Command == "" {
		f.Browser.Command = "google-chrome-stable"
	}
	if f.Browser.UserDataDir == "" {
		f.Browser.UserDataDir = "~/.config/google-chrome"
	}
	ud, err := ExpandPath(f.Browser.UserDataDir)
	if err != nil {
		return nil, err
	}
	f.Browser.UserDataDir = ud
	if f.Default.Profile == "" {
		f.Default.Profile = "Default"
	}
	return &f, nil
}

// ExpandPath expands ~ and $VAR in a path.
func ExpandPath(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	s = os.ExpandEnv(s)
	if strings.HasPrefix(s, "~"+string(os.PathSeparator)) || s == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if s == "~" {
			return home, nil
		}
		return filepath.Join(home, s[2:]), nil
	}
	if strings.HasPrefix(s, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, s[1:]), nil
	}
	return filepath.Clean(s), nil
}
