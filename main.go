// kk-goto/main.go
package main

import (
	"os"

	"github.com/localshred/kk-goto/cmd"
	"github.com/localshred/koolkit/pkg/plugin"
)

func main() {
	meta := plugin.Meta{Describe: "Navigate to SCM provider URLs (repos, issues, PRs)", Requires: []string{}}
	if handled, code := plugin.Handle(os.Args[1:], meta, os.Stdout); handled {
		os.Exit(code)
	}
	os.Exit(cmd.Execute())
}
