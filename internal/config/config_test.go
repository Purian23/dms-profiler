package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFlattenRules(t *testing.T) {
	rules := []Rule{
		{Matches: []string{"https://a/", "https://b/"}, Profile: "Work"},
		{Match: "https://c/", Profile: "Personal"},
		{Matches: []string{}, Match: "https://solo/", Profile: "Default"},
	}
	got := FlattenRules(rules)
	want := []FlatRule{
		{"https://a/", "Work"},
		{"https://b/", "Work"},
		{"https://c/", "Personal"},
		{"https://solo/", "Default"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FlattenRules() = %#v, want %#v", got, want)
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir:", err)
	}
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"~", home},
		{"~/foo/bar", filepath.Join(home, "foo/bar")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}
	for _, tc := range cases {
		got, err := ExpandPath(tc.in)
		if err != nil {
			t.Errorf("ExpandPath(%q) error: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ExpandPath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	cfg := `
[browser]
command = "chromium"
user_data_dir = "/tmp/chrome"

[[rules]]
match = "https://work.example.com/"
profile = "Work"

[default]
profile = "Personal"
`
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if f.Browser.Command != "chromium" {
		t.Errorf("command = %q, want chromium", f.Browser.Command)
	}
	if f.Browser.UserDataDir != "/tmp/chrome" {
		t.Errorf("user_data_dir = %q, want /tmp/chrome", f.Browser.UserDataDir)
	}
	if len(f.Rules) != 1 || f.Rules[0].Match != "https://work.example.com/" {
		t.Errorf("unexpected rules: %v", f.Rules)
	}
	if f.Default.Profile != "Personal" {
		t.Errorf("default profile = %q, want Personal", f.Default.Profile)
	}
}

func TestLoad_Defaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte("[browser]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if f.Browser.Command != "google-chrome-stable" {
		t.Errorf("default command = %q", f.Browser.Command)
	}
	if f.Default.Profile != "Default" {
		t.Errorf("default profile = %q", f.Default.Profile)
	}
}
