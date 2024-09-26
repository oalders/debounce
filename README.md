# debounce

<p align="center">
  <img src="logo.jpeg" alt="debounce logo" />
</p>

<!-- vim-markdown-toc GFM -->

- [Introduction](#introduction)
- [Installation](#installation)
  - [Go Install](#go-install)
  - [Using Ubi](#using-ubi)
  - [Download a Release](#download-a-release)
- [Examples](#examples)
  - [A command without arguments](#a-command-without-arguments)
  - [A command with arguments](#a-command-with-arguments)
  - [Using Shell Variables](#using-shell-variables)
  - [More Complex Commands](#more-complex-commands)
- [Available Flags](#available-flags)
  - [--cache-dir](#--cache-dir)
  - [--local](#--local)
    - [Example Usage](#example-usage)
  - [--status](#--status)
    - [Resetting the Cache](#resetting-the-cache)
  - [--version](#--version)
  - [--help](#--help)
- [Caveats](#caveats)

<!-- vim-markdown-toc -->

## Introduction

`debounce` is a simple utility to limit the rate at which a command can fire.

The command format is:

```bash
debounce <integer> <unit> <command>
```

Available units are:

- seconds (s)
- minutes (m)
- hours (h)
- days (d)

The following are equivalent:

```bash
debounce 1 s some-command
debounce 1 second some-command
debounce 1 seconds some-command
```

```bash
debounce 1 m some-command
debounce 1 minute some-command
debounce 1 minutes some-command
```

```bash
debounce 1 h some-command
debounce 1 hour some-command
debounce 1 hours some-command
```

```bash
debounce 1 d some-command
debounce 1 day some-command
debounce 1 days some-command
```

Note that `day` is shorthand for `24 hours`. Daylight Saving Time is not taken
into account.

## Installation

Choose from the following options to install `debounce`.

### Go Install

```bash
go install github.com/oalders/debounce@latest
```

or for a specific version:

```bash
go install github.com/oalders/debounce@v0.5.0
```

### Using Ubi

You can use [ubi](https://github.com/houseabsolute/ubi) to install `debounce`.

```bash
#!/usr/bin/env bash

set -eux -o pipefail

# Choose a directory in your $PATH
dir="$HOME/local/bin"

if [ ! "$(command -v ubi)" ]; then
    curl --silent --location \
        https://raw.githubusercontent.com/houseabsolute/ubi/master/bootstrap/bootstrap-ubi.sh |
        TARGET=$dir sh
fi

ubi --project oalders/debounce --in "$dir"
```

### Download a Release

You can download the latest release directly from
[here](https://github.com/oalders/debounce/releases).

## Examples

### A command without arguments

This runs the command without any parameters but prevents repeated execution
within the set time window.

```bash
$ debounce 2 seconds date
Mon Aug  5 23:09:09 EDT 2024
$ debounce 2 seconds date
🚥 will not run date more than once every 2 seconds
```

### A command with arguments

This example uses [ubi](https://github.com/houseabsolute/ubi) to install
[is](https://github.com/oalders/is) into the current directory. The command
won't be executed more than once every 8 hours.

```bash
$ debounce 8 hours ubi --verbose --project oalders/is --in .
[ubi::installer][INFO] Installed binary into ./is
$ debounce 8 hours ubi --verbose --project oalders/is --in .
🚥 will not run "ubi --verbose --project oalders/is --in ." more than once every 8 hours
```

### Using Shell Variables

Remember to single quote variables which shouldn't be expanded until the command
is run.

```bash
debounce 10 s zsh -c 'echo $PWD'
```

### More Complex Commands

You can use `&&` and `||` in your commands. You'll want to quote your command to
ensure that the entire command is passed to `debounce`.

```bash
debounce 2 s bash -c 'sleep 2 && date'
```

## Available Flags

It's important to add `debounce` flags before other command arguments to avoid
confusion between `debounce` flags and flags meant for your command.

Good: ✅

```shell
debounce --debug 90 s curl https://www.olafalders.com
```

Bad: 💥

```shell
$ debounce 90 s curl https://www.prettygoodping.com --debug
🚀 Running command: curl https://www.prettygoodping.com --debug
curl: option --debug: is unknown
```

You could be explicit about this by using `--` as a visual indicator that flag
parsing has ended.

```shell
debounce --debug 90 s -- curl https://www.prettygoodping.com
```

### --cache-dir

Specify an alternate cache directory to use. The directory must already exist.

```bash
debounce --cache-dir /tmp 30 s date
```

### --local

The `--local` flag modifies the behavior of `debounce` by including the full
path to the current working directory in the cache key. This ensures that the
cache is specific to the directory from which the command is run.

By default, `debounce` works globally. It generates a cache key based on the
command and its arguments. This means that if you run the same command from
different directories, `debounce` will treat them as the same command and use
the same cache file. This can be useful if you want to limit the execution of a
command globally, regardless of the directory.

With the `--local` flag, `debounce` includes the full path to the current
working directory in the cache key. This means that the cache is specific to the
directory from which the command is run. If you run the same command from
different directories, `debounce` will treat them as different commands and use
separate cache files. This is useful if you want to limit the execution of a
command on a per-directory basis.

#### Example Usage

```shell
# Run a command without the --local flag
debounce 10 s echo "Hello, World!"
# Run the same command again within 10 seconds
debounce 10 s echo "Hello, World!"
# Output: 🚥 will not run "echo Hello, World!" more than once every 10 seconds

# Run a command with the --local flag
debounce --local 10 s echo "Hello, Local World!"
# Run the same command again within 10 seconds from the same directory
debounce --local 10 s echo "Hello, Local World!"
# Output: 🚥 will not run "echo Hello, Local World!" more than once every 10 seconds

# Run the same command from a different directory
cd /another/directory
debounce --local 10 s echo "Hello, Local World!"
# Output: Hello, Local World!
```

### --status

Print debounce status information for a command.

```bash
debounce --status 30 s date
📁 cache location: /Users/olaf/.cache/debounce/0e87632cd46bd4907c516317eb6d81fe0f921a23c7643018f21292894b470681
🚧 cache last modified: Thu, 19 Sep 2024 08:28:20 EDT
⏲️ debounce interval: 00:00:30
🕰️ cache age: 00:00:12
⏳ time remaining: 00:00:17
```

#### Resetting the Cache

Since the cache is just a file, you can `rm` the cache location file whenever
you'd like to start fresh.

```shell
rm /Users/olaf/.cache/debounce/0e87632cd46bd4907c516317eb6d81fe0f921a23c7643018f21292894b470681
```

### --version

Prints current version.

```bash
debounce --version
0.5.0
```

### --help

Displays usage instructions.

```text
debounce --help
Usage: debounce <quantity> <unit> <command> ... [flags]

limit the rate at which a command can fire

Arguments:
  <quantity>       Quantity of time
  <unit>           s,second,seconds,m,minute,minutes,h,hour,hours,d,day,days
  <command> ...    Command to run

Flags:
  -h, --help                Show context-sensitive help.
      --debug               Print debugging info to screen
      --version             Print version to screen
      --status              Print cache information for a command without running it
      --cache-dir=STRING    Override the default cache directory
```

## Caveats

Under the hood, `debounce` creates or updates a cache file to track when a
command was run successfully. This means that, under the right conditions, it's
entirely possible to kick off two long-running tasks in parallel without
`debounce` knowing about it.

Additionally, if a command fails, the cache file will not be created or updated.

I've created this tool in a way that meets my needs. I will consider pull
requests for additional functionality to address issues like these. Please get
in touch with me first to discuss your feature if you'd like to add something.
