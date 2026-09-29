# kk-goto

kk-goto is a [koolkit](https://github.com/localshred/koolkit) (`kk`)
plugin that builds and opens SCM navigation URLs — a user or org page, a
repo, a repo's issues or pull requests list, or a specific issue or pull
request. User and repo names can be short aliases, resolved from a local
JSON config file before the URL is built.

## Install

kk-goto is installed and run through `kk`, the koolkit plugin manager.
See the [koolkit repo](https://github.com/localshred/koolkit) for
installing `kk` itself.

```
kk install goto
```

## Usage

```
kk goto -u octocat -r Hello-World
```

This resolves `octocat`/`Hello-World` to a URL and opens it with the
`open` command. Add `-n` with `-i` or `-p` to target a specific issue or
pull request, or `-y`/`--dry` to print the URL instead of opening it:

```
kk goto -u octocat -r Hello-World -n 42 -p -y
```

See [docs/usage.md](docs/usage.md) for the full flag reference, the
positional-argument shorthand, and the alias file format.

## Requirements

kk-goto declares no external tools in its plugin metadata
(`plugin.Meta.Requires` is empty), so `kk doctor` has nothing to check
for it. At runtime it shells out to the `open` command to launch URLs,
except when run with `--dry`.
