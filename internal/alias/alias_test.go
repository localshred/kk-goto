package alias

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeJSON(t *testing.T, dir, body string) string {
	t.Helper()
	p := filepath.Join(dir, "goto.json")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func writeJSONNamed(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestResolveUserAndRepo(t *testing.T) {
	dir := t.TempDir()
	p := writeJSON(t, dir, `{
	  "l": {"user": "localshred", "scm": "github", "repos": {"d": "dotfiles"}}
	}`)
	a := Load([]string{p})
	user, host := a.ResolveUser("l")
	if user != "localshred" || host != "" {
		t.Fatalf("ResolveUser(l) = %q,%q", user, host)
	}
	if got := a.ResolveRepo("l", "d"); got != "dotfiles" {
		t.Fatalf("ResolveRepo(l,d) = %q", got)
	}
}

func TestResolveUnknownPassesThrough(t *testing.T) {
	a := Load(nil) // no config at all
	user, _ := a.ResolveUser("foo")
	if user != "foo" {
		t.Fatalf("unknown user should pass through, got %q", user)
	}
	if got := a.ResolveRepo("foo", "bar"); got != "bar" {
		t.Fatalf("unknown repo should pass through, got %q", got)
	}
}

func TestLoadSkipsMalformed(t *testing.T) {
	dir := t.TempDir()
	p := writeJSON(t, dir, `{ not json `)
	if a := Load([]string{p}); len(a) != 0 {
		t.Fatalf("malformed config should load empty, got %v", a)
	}
}

func TestLoadEarlierFilesWin(t *testing.T) {
	dir := t.TempDir()
	p1 := writeJSONNamed(t, dir, "first.json", `{"l": {"user": "first", "scm": "github"}}`)
	p2 := writeJSONNamed(t, dir, "second.json", `{"l": {"user": "second", "scm": "gitlab"}}`)
	a := Load([]string{p1, p2})
	user, _ := a.ResolveUser("l")
	if user != "first" {
		t.Fatalf("earlier file should win, got %q (expected first)", user)
	}
}

func TestDefaultPaths(t *testing.T) {
	// Test with KK_GOTO_CONFIG containing empty segments
	t.Run("empty-segments-filtered", func(t *testing.T) {
		t.Setenv("KK_GOTO_CONFIG", ":x:y::z:")
		t.Setenv("XDG_CONFIG_HOME", "")
		paths := DefaultPaths()
		if len(paths) != 3 {
			t.Fatalf("expected 3 paths, got %d: %v", len(paths), paths)
		}
		if paths[0] != "x" || paths[1] != "y" || paths[2] != "z" {
			t.Fatalf("expected [x y z], got %v", paths)
		}
	})

	// Test default path expansion
	t.Run("default-path-expands-home", func(t *testing.T) {
		t.Setenv("KK_GOTO_CONFIG", "")
		t.Setenv("XDG_CONFIG_HOME", "")
		paths := DefaultPaths()
		if len(paths) != 1 {
			t.Fatalf("expected 1 default path, got %d: %v", len(paths), paths)
		}
		p := paths[0]
		if !filepath.IsAbs(p) {
			t.Fatalf("expected absolute path, got %q", p)
		}
		if !strings.Contains(p, ".config") || !strings.Contains(p, "koolkit") {
			t.Fatalf("path missing .config/koolkit, got %q", p)
		}
		if strings.Contains(p, "~") {
			t.Fatalf("path contains literal ~, got %q", p)
		}
	})
}
