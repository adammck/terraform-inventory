package main

import (
	"fmt"
	"os"
	"flag"
	"path/filepath"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	flag.Parse()
	file := flag.Arg(0)

	if file == "" {
		fmt.Printf("Usage: %s PATH\n", os.Args[0])
		os.Exit(1)
	}

	path, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("Invalid file: %s\n", err)
		os.Exit(1)
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening tfstate file: %s\n", err)
		os.Exit(1)
	}

	state, err := terraform.ReadState(f)
	if err != nil {
		fmt.Printf("Error reading state: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("[production]")

	for _, resource := range state.Resources {
		if resource.Type == "aws_instance" {
			fmt.Printf("%s serverName=%s\n", resource.Attributes["public_ip"], resource.Attributes["id"])
		}
	}
}
