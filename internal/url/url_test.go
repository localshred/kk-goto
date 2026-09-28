// kk-goto/internal/url/url_test.go
package url

import "testing"

func TestBuildURL(t *testing.T) {
	cases := []struct {
		name string
		o    Opts
		want string
	}{
		{"user", Opts{User: "foo"}, "https://github.com/foo"},
		{"repo", Opts{User: "foo", Repo: "bar"}, "https://github.com/foo/bar"},
		{"issues", Opts{User: "foo", Repo: "bar", Issues: true}, "https://github.com/foo/bar/issues"},
		{"issue", Opts{User: "foo", Repo: "bar", Issues: true, Number: 12}, "https://github.com/foo/bar/issues/12"},
		{"pulls", Opts{User: "foo", Repo: "bar", Pulls: true}, "https://github.com/foo/bar/pulls"},
		{"pull", Opts{User: "foo", Repo: "bar", Pulls: true, Number: 34}, "https://github.com/foo/bar/pull/34"},
		{"custom host", Opts{Host: "https://ghe.corp", User: "foo"}, "https://ghe.corp/foo"},
	}
	for _, c := range cases {
		got, ok := BuildURL(c.o)
		if !ok || got != c.want {
			t.Errorf("%s: BuildURL = %q ok=%v, want %q", c.name, got, ok, c.want)
		}
	}
}

func TestBuildURLMissingUser(t *testing.T) {
	if _, ok := BuildURL(Opts{}); ok {
		t.Fatal("no user should not build a URL")
	}
}
