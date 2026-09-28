// Package alias loads user/repo aliases for kk-goto from JSON config files.
// Missing or malformed files are skipped; unknown names pass through.
package alias

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	User  string            `json:"user"`
	Host  string            `json:"host"`
	SCM   string            `json:"scm"`
	Repos map[string]string `json:"repos"`
}

type Aliases map[string]Entry

// DefaultPaths returns the list of config paths from KK_GOTO_CONFIG (colon-separated)
// or defaults to ~/.config/koolkit/goto.json.
func DefaultPaths() []string {
	if v := os.Getenv("KK_GOTO_CONFIG"); v != "" {
		return strings.Split(v, ":")
	}
	cfg := os.Getenv("XDG_CONFIG_HOME")
	if cfg == "" {
		home, _ := os.UserHomeDir()
		cfg = filepath.Join(home, ".config")
	}
	return []string{filepath.Join(cfg, "koolkit", "goto.json")}
}

// Load merges alias files; earlier paths win on key conflict. Never errors;
// missing and malformed files are skipped silently.
func Load(paths []string) Aliases {
	out := Aliases{}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var a Aliases
		if json.Unmarshal(data, &a) != nil {
			continue
		}
		for k, v := range a {
			if _, seen := out[k]; !seen {
				out[k] = v
			}
		}
	}
	return out
}

// ResolveUser returns the mapped user and host for the given name,
// or the literal name unchanged if no alias exists.
func (a Aliases) ResolveUser(name string) (string, string) {
	if e, ok := a[name]; ok && e.User != "" {
		return e.User, e.Host
	}
	return name, ""
}

// ResolveRepo returns the mapped repo for the given user alias and repo,
// or the literal repo unchanged if no alias exists.
func (a Aliases) ResolveRepo(userAlias, repo string) string {
	if e, ok := a[userAlias]; ok {
		if r, ok := e.Repos[repo]; ok {
			return r
		}
	}
	return repo
}
