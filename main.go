package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/oalders/debounce/run"
	"github.com/oalders/debounce/types"
)

const vstring = "0.2.0"

func main() {
	quantity, unit, command, debug, status, cacheDir := parseAndValidateFlags()

	runContext := types.DebounceCommand{
		Quantity: quantity,
		Unit:     unit,
		Command:  command,
		Debug:    debug,
		Status:   status,
		CacheDir: cacheDir,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	success, output, err := run.Run(&runContext, home)
	fmt.Print(string(output))
	if err != nil {
		fmt.Println("there was an error")
		fmt.Println(err)
		os.Exit(1)
	}

	if success {
		os.Exit(0)
	}
	os.Exit(1)
}

func parseAndValidateFlags() (string, string, []string, bool, bool, string) {
	debug := flag.Bool("debug", false, "Print debugging info to screen")
	version := flag.Bool("version", false, "Print version to screen")
	status := flag.Bool("status", false, "Print cache information for a command without running it")
	cacheDir := flag.String("cache-dir", "", "Override the default cache directory")
	flag.Parse()
	if *version {
		fmt.Printf("debounce %s\n", vstring)
		os.Exit(0)
	}

	args := os.Args[1:]

	if len(args) < 3 {
		printUsage()
		os.Exit(1)
	}

	quantity := flag.Args()[0]
	unit := flag.Args()[1]
	command := flag.Args()[2:]

	if _, err := strconv.Atoi(quantity); err != nil {
		fmt.Printf("quantity %s is not a valid integer\n", quantity)
		printUsage()
		os.Exit(1)
	}

	switch unit {
	case "minutes", "minute", "m":
		unit = "m"
	case "seconds", "second", "s":
		unit = "s"
	case "hours", "hour", "h":
		unit = "h"
	default:
		fmt.Printf("unit \"%s\" is invalid. Valid units are: hour(s), minute(s) and second(s)\n", unit)
		fmt.Println("You may also use the shorthand forms: h, m and s")
		printUsage()
		os.Exit(1)
	}

	return quantity, unit, command, *debug, *status, *cacheDir
}

func printUsage() {
	fmt.Println(`Usage: debounce <integer> <hours|minutes|seconds> <command>
eg: debounce 2 hours date
eg: debounce 10 minutes zsh -c 'echo $PWD'
eg: debounce 5 seconds bash -c 'sleep 2 && date'`)
}
