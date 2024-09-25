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
  * [Available Flags](#available-flags)
  * [--cache-dir](#--cache-dir)
  * [--status](#--status)
    * [Resetting the Cache](#resetting-the-cache)
  * [--version](#--version)
  * [--help](#--help)
* [Installation](#installation)
* [Caveats](#caveats)

<!-- vim-markdown-toc -->

## Introduction

`debounce` is a simple utility to limit the rate at which a command can fire.

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

### Available Flags

Keep in mind that you should add `debounce` flags *before* any other arguments
to ensure that they are not interpreted as part of the command you are running.

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
./bin/debounce --cache-dir /tmp 30 s date
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
./bin/debounce --version
0.2.0
```

### --help

Displays usage instructions.

```text
./bin/debounce --help
Usage: debounce <quantity> <unit> <command> ... [flags]

limit the rate at which a command can fire

Arguments:
  <quantity>       Quantity of time
  <unit>           s,second,seconds,m,minute,minutes,h,hour,hours
  <command> ...    Command to run

Flags:
  -h, --help                Show context-sensitive help.
      --debug               Print debugging info to screen
      --version             Print version to screen
      --status              Print cache information for a command without running it
      --cache-dir=STRING    Override the default cache directory
```

## Installation

Choose from the following options to install `debounce`.

1. [Download a release](https://github.com/oalders/debounce/releases)
2. Use `go install`
  * `go install github.com/oalders/debounce@latest`
  * `go install github.com/oalders/debounce@v0.2.0`
3. Use [ubi](https://github.com/houseabsolute/ubi)

```bash
#!/usr/bin/env bash

set -eux -o pipefail

# Or choose a different dir in your $PATH
dir="$HOME/local/bin"

if [ ! "$(command -v ubi)" ]; then
    curl --silent --location \
        https://raw.githubusercontent.com/houseabsolute/ubi/master/bootstrap/bootstrap-ubi.sh |
        TARGET=$dir sh
fi

ubi --project oalders/debounce --in "$dir"
```

## Caveats

Under the hood, `debounce` creates or updates a cache file to track when a
command was run successfully. This means that under the right conditions, it's
entirely possibly to kick off two long running tasks in parallel without
debounce knowing about it.

Additionally, if a command fails, the cache file will not be created or
updated.

I've created this tool in a way that meets my needs. I will consider pull
requests for additional functionality to address issues like these. Please get
in touch with me first to discuss your feature if you'd like to add something.
