package chromium

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFolderName(t *testing.T) {
	cases := []struct {
		s    string
		want bool
	}{
		{"Default", true},
		{"Profile 1", true},
		{"Profile 12", true},
		{"Work", false},
		{"Profile", false},
		{"Profile x", false},
	}
	for _, tc := range cases {
		if got := isFolderName(tc.s); got != tc.want {
			t.Errorf("isFolderName(%q) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

func TestFolderForDisplayName(t *testing.T) {
	raw := []byte(`{
  "profile": {
    "info_cache": {
      "Default": { "name": "Personal" },
      "Profile 1": { "name": "Work" }
    }
  }
}`)
	folder, err := folderForDisplayName(raw, "Work")
	if err != nil {
		t.Fatal(err)
	}
	if folder != "Profile 1" {
		t.Fatalf("got %q, want Profile 1", folder)
	}
}

func TestResolveProfileDir_FolderNameSkipsFile(t *testing.T) {
	dir := t.TempDir()
	got, err := ResolveProfileDir(dir, "Profile 2")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Profile 2" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveProfileDir_FromFile(t *testing.T) {
	dir := t.TempDir()
	ls := filepath.Join(dir, "Local State")
	raw := []byte(`{
  "profile": {
    "info_cache": {
      "Default": { "name": "Personal" },
      "Profile 1": { "name": "Work" }
    }
  }
}`)
	if err := os.WriteFile(ls, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ResolveProfileDir(dir, "Work")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Profile 1" {
		t.Fatalf("got %q", got)
	}
}
