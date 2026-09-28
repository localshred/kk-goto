// Package cmd implements the kk-goto command.
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/localshred/kk-goto/internal/alias"
	"github.com/localshred/kk-goto/internal/url"
	"github.com/localshred/koolkit/pkg/run"
	"github.com/spf13/cobra"
)

func Execute() int {
	var (
		user, repo string
		number     int
		issues     bool
		pulls      bool
		dry        bool
	)
	root := &cobra.Command{
		Use:           "kk-goto",
		Short:         "Navigate to SCM provider URLs (repos, issues, PRs)",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			o := url.Opts{User: strings.ToLower(user), Repo: strings.ToLower(repo),
				Number: number, Issues: issues, Pulls: pulls}
			uAlias := o.User
			if len(args) > 0 {
				o, uAlias = parsePositional(o, args)
			}
			a := alias.Load(alias.DefaultPaths())
			o = resolve(a, o, uAlias)
			u, ok := url.BuildURL(o)
			if !ok {
				return fmt.Errorf("not enough info to build a URL — give at least -u/--user")
			}
			if code := openURL(run.Default, u, dry, cmd.OutOrStdout()); code != 0 {
				os.Exit(code)
			}
			return nil
		},
	}
	f := root.Flags()
	f.StringVarP(&user, "user", "u", "", "User or org to open")
	f.StringVarP(&repo, "repo", "r", "", "Repo to open")
	f.IntVarP(&number, "number", "n", 0, "Issue or PR number")
	f.BoolVarP(&issues, "issues", "i", false, "Navigate to repo issues")
	f.BoolVarP(&pulls, "pull-requests", "p", false, "Navigate to repo pull requests")
	f.BoolVarP(&dry, "dry", "y", false, "Print the URL and exit without opening")
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "kk-goto:", err)
		return 1
	}
	return 0
}

// parsePositional fills user/repo/number from `owner repo number` or
// `owner/repo#number` positional args (mx parse-from-args behavior).
func parsePositional(o url.Opts, args []string) (url.Opts, string) {
	fields := strings.FieldsFunc(strings.Join(args, " "), func(r rune) bool {
		return r == ' ' || r == '/' || r == '#'
	})
	if len(fields) > 0 {
		o.User = strings.ToLower(fields[0])
	}
	if len(fields) > 1 {
		o.Repo = strings.ToLower(fields[1])
	}
	if len(fields) > 2 {
		if n, err := strconv.Atoi(fields[2]); err == nil {
			o.Number = n
			if !o.Pulls {
				o.Issues = true
			}
		}
	}
	return o, o.User
}

func resolve(a alias.Aliases, o url.Opts, uAlias string) url.Opts {
	user, host := a.ResolveUser(uAlias)
	o.User = user
	if host != "" {
		o.Host = host
	}
	if o.Repo != "" {
		o.Repo = a.ResolveRepo(uAlias, o.Repo)
	}
	return o
}

func openURL(r run.Runner, u string, dry bool, out io.Writer) int {
	if dry {
		fmt.Fprintf(out, "[DRY] Goto URL %s\n", u)
		return 0
	}
	fmt.Fprintf(out, "Opening %s\n", u)
	res, err := r.Run("open", u)
	if err != nil {
		fmt.Fprintln(out, res.Stderr)
		return 1
	}
	return 0
}
