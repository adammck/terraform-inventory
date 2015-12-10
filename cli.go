package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func cmdList(stdout io.Writer, stderr io.Writer, s *state) int {
	groups := make(map[string][]string, 0)

	// Add each instance name as a pseudo-group, so they can be provisioned
	// individually where necessary.
	for _, res := range s.resources() {
		_, ok := groups[res.Name]
		if !ok {
			groups[res.Name] = []string{}
		}

		// Add the instance by name. There can be many instances with the same name,
		// created using the count parameter.
		groups[res.Name] = append(groups[res.Name], res.Address())

		groups[res.NameWithCounter()] = []string{res.Address()}
	}

	return output(stdout, stderr, groups)
}

func cmdHost(stdout io.Writer, stderr io.Writer, s *state, hostname string) int {
	for _, res := range s.resources() {
		if hostname == res.Name {
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
