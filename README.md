# `chgo` Change Go

Running `chgo VERSION` will install and link to requested version of Go to `$HOME/bin/go`.

Examples:
- `chgo go1.17.6`
- `chgo gotip`
- `chgo latest` - installs the highest supported stable version from [go.dev/dl](go.dev/dl)

## List Mode

`chgo` can also list available versions of Go by using the `--list` flag. This mode will not install and link Go versions.

Examples:
- `chgo --list` - lists all availble versions
- `chgo --list go1.17` list all versions of Go 1.17
- `chgo --list latest` lists the latest supported stable versions