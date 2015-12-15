package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func cmdList(stdout io.Writer, stderr io.Writer, s *state) int {
	groups := make(map[string][]string, 0)

	// add each instance as a pseudo-group, so they can be provisioned
	// individually where necessary.
	for name, res := range s.resources() {
		g := res.Groups()
		for _, i := range g {
			if i != "" {
				groups[i] = append(groups[i], res.Address())
			}
		}
		groups[name] = []string{res.Address()}
	}

	return output(stdout, stderr, groups)
}

func cmdHost(stdout io.Writer, stderr io.Writer, s *state, hostname string) int {
	for name, res := range s.resources() {
		if hostname == name {
			return output(stdout, stderr, res.Attributes())
		}
	}

	fmt.Fprintf(stderr, "No such host: %s\n", hostname)
	return 1
}

// output marshals an arbitrary JSON object and writes it to stdout, or writes
// an error to stderr, then returns the appropriate exit code.
func output(stdout io.Writer, stderr io.Writer, whatever interface{}) int {
	b, err := json.Marshal(whatever)
	if err != nil {
		fmt.Fprintf(stderr, "Error encoding JSON: %s\n", err)
		return 1
	}

	_, err = stdout.Write(b)
	if err != nil {
		fmt.Fprintf(stderr, "Error writing JSON: %s\n", err)
		return 1
	}

	return 0
}
