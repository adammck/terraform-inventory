package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var version = flag.Bool("version", false, "print version information and exit")
var list = flag.Bool("list", false, "list mode")
var host = flag.String("host", "", "host mode")
var inventory = flag.Bool("inventory", false, "inventory mode")

func main() {
	flag.Parse()
	file := flag.Arg(0)

	if *version == true {
		fmt.Printf("%s version %s\n", os.Args[0], versionInfo())
		return
	}

	// not given on the command line? try ENV.
	if file == "" {
		file = os.Getenv("TF_STATE")
	}

	// also try the old ENV name.
	if file == "" {
		file = os.Getenv("TI_TFSTATE")
	}

	// check for a file named terraform.tfstate in the pwd
	if file == "" {
		fn := "terraform.tfstate"
		_, err := os.Stat(fn)
		if err == nil {
			file = fn
		}
	}

	if file == "" {
		fmt.Printf("Usage: %s [options] path\n", os.Args[0])
		os.Exit(1)
	}

	if !*list && *host == "" && !*inventory {
		fmt.Fprint(os.Stderr, "Either --host or --list must be specified")
		os.Exit(1)
	}

	path, err := filepath.Abs(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid file: %s\n", err)
		os.Exit(1)
	}

	stateFile, err := os.Open(path)
	defer stateFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening tfstate file: %s\n", err)
		os.Exit(1)
	}

	var s state
	err = s.read(stateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading tfstate file: %s\n", err)
		os.Exit(1)
	}

	if *list {
		os.Exit(cmdList(os.Stdout, os.Stderr, &s))

	} else if *inventory {
		os.Exit(cmdInventory(os.Stdout, os.Stderr, &s))

	} else if *host != "" {
		os.Exit(cmdHost(os.Stdout, os.Stderr, &s, *host))

	}
}
