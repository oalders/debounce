# debounce

<p align="center">
  <img src="logo.jpeg" />
</p>


<!-- vim-markdown-toc GFM -->

* [Introduction](#introduction)
* [Examples](#examples)
  * [A command without arguments](#a-command-without-arguments)
  * [A command with arguments](#a-command-with-arguments)
  * [Using Shell Variables](#using-shell-variables)
  * [More Complex Commands](#more-complex-commands)
* [Installation](#installation)

<!-- vim-markdown-toc -->

## Introduction

debounce is a simple utility to limit the rate at which a command can fire.

The arguments are:

```bash
debounce <integer> <unit> <command>
```

Available units are:

* seconds (s)
* minutes (m)
* hours (h)

The following are equivalent:

```bash
debounce 1 s date
debounce 1 second date
debounce 1 seconds date
```

```bash
debounce 1 m date
debounce 1 minute date
debounce 1 minutes date
```

```bash
debounce 1 h date
debounce 1 hour date
debounce 1 hours date
```

## Examples

### A command without arguments

```bash
$ debounce 2 seconds date
Mon Aug  5 23:09:09 EDT 2024
$ debounce 2 seconds date
🚥 will not run date more than once every 2 seconds
```

### A command with arguments

This command uses <https://github.com/houseabsolute/ubi> to install
<https://github.com/oalders/is> into the current directory.  The command will
not be run more than once every 8 hours.

```bash
$ debounce 8 hours ubi --verbose --project oalders/is --in .
[ubi::installer][INFO] Installed binary into ./is
$ debounce 8 hours ubi --verbose --project oalders/is --in .
🚥 will not run "ubi --verbose --project oalders/is --in ." more than once every 8 hours
```

### Using Shell Variables

Remember to single quote variables which shouldn't be expanded until the
command is run.

```bash
debounce 10 s zsh -c 'echo $PWD'
```

### More Complex Commands

You can use `&&` and `||` in your commands. You'll want to quote your command
to ensure that the entire command is passed to `debounce`.

```bash
debounce 2 s bash -c 'sleep 2 && date'
```

## Installation

Choose from the following options to install `debounce`.

1. [Download a release](https://github.com/oalders/debounce/releases)
1. Use `go install`
  * `go install github.com/oalders/debounce@latest`
  * `go install github.com/oalders/debounce@v0.1.0`
1. Use [ubi](https://github.com/houseabsolute/ubi)

```bash
#!/usr/bin/env bash

set -e -u -x -o pipefail

# Or choose a different dir in your $PATH
dir="$HOME/local/bin"

if [ ! "$(command -v ubi)" ]; then
    curl --silent --location \
        https://raw.githubusercontent.com/houseabsolute/ubi/master/bootstrap/bootstrap-ubi.sh |
        TARGET=$dir sh
fi

ubi --project oalders/debounce --in "$dir"
```
