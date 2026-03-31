package chromium

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolveProfileDir maps a profile label or folder name to a Chromium profile directory name (e.g. Default, Profile 1).
func ResolveProfileDir(userDataDir, spec string) (string, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return "", fmt.Errorf("empty profile")
	}
	if isFolderName(spec) {
		return spec, nil
	}
	path := filepath.Join(userDataDir, "Local State")
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read Local State: %w", err)
	}
	folder, err := folderForDisplayName(raw, spec)
	if err != nil {
		return "", err
	}
	return folder, nil
}

func isFolderName(s string) bool {
	if s == "Default" || s == "System Profile" {
		return true
	}
	if strings.HasPrefix(s, "Profile ") {
		rest := strings.TrimPrefix(s, "Profile ")
		for _, c := range rest {
			if c < '0' || c > '9' {
				return false
			}
		}
		return len(rest) > 0
	}
	return false
}

type localState struct {
	Profile struct {
		InfoCache map[string]profileEntry `json:"info_cache"`
	} `json:"profile"`
}

type profileEntry struct {
	Name string `json:"name"`
}

func folderForDisplayName(raw []byte, displayName string) (string, error) {
	var ls localState
	if err := json.Unmarshal(raw, &ls); err != nil {
		return "", fmt.Errorf("parse Local State: %w", err)
	}
	if ls.Profile.InfoCache == nil {
		return "", fmt.Errorf("no profiles in Local State")
	}
	var matches []string
	for folder, ent := range ls.Profile.InfoCache {
		if ent.Name == displayName {
			matches = append(matches, folder)
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no profile named %q in Local State", displayName)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("ambiguous profile name %q (folders: %v)", displayName, matches)
	}
}
