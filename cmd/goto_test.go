package cmd

import (
	"bytes"
	"testing"

	"github.com/localshred/kk-goto/internal/alias"
	"github.com/localshred/kk-goto/internal/url"
	"github.com/localshred/koolkit/pkg/run"
)

func TestDefaultTargetNumberImpliesIssues(t *testing.T) {
	got := defaultTarget(url.Opts{User: "foo", Repo: "bar", Number: 5})
	if !got.Issues || got.Pulls {
		t.Fatalf("number without -i/-p should default to issues, got %+v", got)
	}
	u, ok := url.BuildURL(got)
	if !ok || u != "https://github.com/foo/bar/issues/5" {
		t.Fatalf("url = %q ok=%v", u, ok)
	}
}

func TestDefaultTargetRespectsExplicitPulls(t *testing.T) {
	got := defaultTarget(url.Opts{User: "foo", Repo: "bar", Number: 5, Pulls: true})
	if got.Issues {
		t.Fatalf("explicit -p must not be overridden by issues default, got %+v", got)
	}
}

func TestResolveAppliesAliases(t *testing.T) {
	a := alias.Aliases{"l": {User: "localshred", Repos: map[string]string{"d": "dotfiles"}}}
	got := resolve(a, url.Opts{User: "l", Repo: "d"}, "l")
	if got.User != "localshred" || got.Repo != "dotfiles" {
		t.Fatalf("resolve = %+v", got)
	}
}

func TestOpenDryPrintsURL(t *testing.T) {
	var out bytes.Buffer
	f := &run.Fake{}
	if code := openURL(f, "https://github.com/foo", true, &out); code != 0 {
		t.Fatalf("dry exit = %d", code)
	}
	if len(f.Calls) != 0 {
		t.Fatal("dry run must not shell out")
	}
	if !bytes.Contains(out.Bytes(), []byte("https://github.com/foo")) {
		t.Fatalf("dry output = %q", out.String())
	}
}

func TestOpenShellsOpen(t *testing.T) {
	var out bytes.Buffer
	f := &run.Fake{Present: map[string]bool{"open": true}}
	openURL(f, "https://github.com/foo", false, &out)
	if len(f.Calls) != 1 || f.Calls[0][0] != "open" || f.Calls[0][1] != "https://github.com/foo" {
		t.Fatalf("calls = %v", f.Calls)
	}
}
