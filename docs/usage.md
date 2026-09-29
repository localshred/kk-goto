# kk-goto usage

kk-goto is a single-command plugin: there are no subcommands. It is
invoked as `kk goto [flags]` or `kk goto <positional args>`, and it
prints its metadata for `kk` via the `--kk-describe` / `--kk-requires`
discovery flags (handled before any of the flags below).

## Flags

| Flag                    | Type   | Default | Meaning                        |
|--------------------------|--------|---------|---------------------------------|
| `-u`, `--user`           | string | `""`    | User or org to open             |
| `-r`, `--repo`           | string | `""`    | Repo to open                    |
| `-n`, `--number`         | int    | `0`     | Issue or PR number              |
| `-i`, `--issues`         | bool   | `false` | Navigate to repo issues         |
| `-p`, `--pull-requests`  | bool   | `false` | Navigate to repo pull requests  |
| `-y`, `--dry`            | bool   | `false` | Print the URL, don't open it    |

`--user` and `--repo` are lowercased before alias resolution.

Examples:

```
kk goto -u octocat
kk goto -u octocat -r Hello-World
kk goto -u octocat -r Hello-World -i
kk goto -u octocat -r Hello-World -p
kk goto -u octocat -r Hello-World -n 42 -i
kk goto -u octocat -r Hello-World -n 42 -p
kk goto -u octocat -r Hello-World -n 42 -y
```

If `--number` is set but neither `-i` nor `-p` is given, kk-goto
defaults the target to issues.

If no user can be resolved (from a flag, a positional arg, or an
alias), kk-goto exits with:

```
kk-goto: not enough info to build a URL — give at least -u/--user
```

## Positional arguments

Instead of flags, kk-goto accepts `owner`, `repo`, and `number` as
positional arguments, in either of these forms:

```
kk goto octocat Hello-World 42
kk goto octocat/Hello-World#42
```

Positional args are joined and split on spaces, `/`, and `#`. The first
field becomes the user, the second the repo, and the third (if it
parses as an integer) the number. As with the `--number` flag, a
positional number defaults the target to issues unless `-p` was also
given.

## Alias resolution

User and repo names can be short aliases defined in a JSON config file
instead of literal names.

### File location

kk-goto reads `KK_GOTO_CONFIG`, a colon-separated list of file paths
(empty segments are ignored). If unset, it falls back to a single
default path: `$XDG_CONFIG_HOME/koolkit/goto.json`, or
`~/.config/koolkit/goto.json` if `XDG_CONFIG_HOME` is unset.

When multiple paths are configured, they are merged; if the same alias
key appears in more than one file, the entry from the earlier path
wins. Missing or malformed files are skipped silently — kk-goto never
errors on config loading.

### File format

The file is a JSON object keyed by alias name:

```json
{
  "oc": {
    "user": "octocat",
    "host": "https://github.com",
    "scm": "github",
    "repos": {
      "hw": "Hello-World"
    }
  }
}
```

- `user` — the real user/org name the alias resolves to.
- `host` — optional; overrides the default `https://github.com` host
  for URLs built for this alias.
- `scm` — informational; not read by the current URL-building logic.
- `repos` — a map of repo alias to real repo name, scoped to this user
  alias.

With the file above:

```
kk goto -u oc -r hw
```

resolves to `https://github.com/octocat/Hello-World`.

### Resolution rules

- The `--user` (or positional owner) value is looked up as an alias
  key. If found and its `user` field is non-empty, the alias's `user`
  and `host` replace the input; otherwise the input is used unchanged
  as a literal name (default host `https://github.com`).
- The `--repo` (or positional repo) value is looked up in the
  *original* user alias's `repos` map (not the resolved user's). If
  found, the mapped repo name is used; otherwise the input is used
  unchanged.
- Unknown aliases always pass through as literal names rather than
  causing an error.

## Selection

kk-goto's alias resolution is a direct, deterministic lookup keyed on
the alias config file — there is no interactive or fuzzy-match
selection step in the current code. (`go-fuzzyfinder` appears in
`go.mod` only as an indirect/transitive dependency pulled in through
koolkit; kk-goto's own code never invokes it.)

## URL construction

Given the resolved user, repo, host, number, and issues/pulls flags,
kk-goto builds one of:

| Condition                          | URL shape                              |
|-------------------------------------|-----------------------------------------|
| issues + number                     | `{host}/{user}/{repo}/issues/{number}`  |
| pull requests + number              | `{host}/{user}/{repo}/pull/{number}`    |
| issues, no number                   | `{host}/{user}/{repo}/issues`           |
| pull requests, no number            | `{host}/{user}/{repo}/pulls`            |
| repo set, no issues/pulls           | `{host}/{user}/{repo}`                  |
| no repo                             | `{host}/{user}`                         |

`{host}` defaults to `https://github.com` unless an alias supplied a
different host.

## Dry run and opening

Without `--dry`, kk-goto prints `Opening <url>` and shells out to the
`open` command with the URL. If that command fails, it prints the
command's stderr and exits with status 1.

With `-y`/`--dry`, kk-goto prints `[DRY] Goto URL <url>` and exits 0
without shelling out at all:

```
kk goto -u octocat -r Hello-World -y
[DRY] Goto URL https://github.com/octocat/Hello-World
```
