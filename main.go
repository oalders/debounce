package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/oalders/debounce/run"
	"github.com/oalders/debounce/types"
)

func main() {
	debug := flag.Bool("debug", false, "Print debugging info to screen")
	version := flag.Bool("version", false, "Print version to screen")
	flag.Parse()
	if *version {
		fmt.Println("debounce 0.1.0")
		os.Exit(0)
	}

	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println(`Usage: debounce <integer> <hours|minutes|seconds> <command>
eg: debounce 2 hours date
eg: debounce 10 minutes zsh -c 'echo $PWD'
eg: debounce 5 seconds bash -c 'sleep 2 && date'`)
		os.Exit(1)
	}

	runContext := types.DebounceCommand{
		Quantity: flag.Args()[0],
		Unit:     flag.Args()[1],
		Command:  flag.Args()[2:],
		Debug:    *debug,
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
