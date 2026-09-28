package alias

import (
	"os"
	"path/filepath"
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
