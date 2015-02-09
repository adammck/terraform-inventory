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

func main() {
	flag.Parse()
	file := flag.Arg(0)

	if *version == true {
		fmt.Printf("%s version %d\n", os.Args[0], versionInfo())
		return
	}

	if file == "" {
		fmt.Printf("Usage: %s [options] path\n", os.Args[0])
		os.Exit(1)
	}

	if !*list && *host == "" {
		fmt.Println("Either --host or --list must be specified")
		os.Exit(1)
	}

	path, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("Invalid file: %s\n", err)
		os.Exit(1)
	}

	stateFile, err := os.Open(path)
	defer stateFile.Close()
	if err != nil {
		fmt.Printf("Error opening tfstate file: %s\n", err)
		os.Exit(1)
	}

	var s state
	err = s.read(stateFile)
	if err != nil {
		fmt.Printf("Error reading tfstate file: %s\n", err)
		os.Exit(1)
	}

	if *list {
		os.Exit(cmdList(os.Stdout, os.Stderr, &s))

	} else if *host != "" {
		os.Exit(cmdHost(os.Stdout, os.Stderr, &s, *host))

	}
}
