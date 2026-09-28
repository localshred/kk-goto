// Package url builds SCM (GitHub) navigation URLs. Ported from mx goto.
package url

import "fmt"

type Opts struct {
	Host   string
	User   string
	Repo   string
	Number int
	Issues bool
	Pulls  bool
}

func (o Opts) host() string {
	if o.Host != "" {
		return o.Host
	}
	return "https://github.com"
}

// BuildURL returns the navigation URL for o. ok is false when required
// fields for any target are missing (no user).
func BuildURL(o Opts) (string, bool) {
	if o.User == "" {
		return "", false
	}
	switch {
	case o.Issues && o.Number > 0:
		return fmt.Sprintf("%s/%s/%s/issues/%d", o.host(), o.User, o.Repo, o.Number), true
	case o.Pulls && o.Number > 0:
		return fmt.Sprintf("%s/%s/%s/pull/%d", o.host(), o.User, o.Repo, o.Number), true
	case o.Issues:
		return fmt.Sprintf("%s/%s/%s/issues", o.host(), o.User, o.Repo), true
	case o.Pulls:
		return fmt.Sprintf("%s/%s/%s/pulls", o.host(), o.User, o.Repo), true
	case o.Repo != "":
		return fmt.Sprintf("%s/%s/%s", o.host(), o.User, o.Repo), true
	default:
		return fmt.Sprintf("%s/%s", o.host(), o.User), true
	}
}
