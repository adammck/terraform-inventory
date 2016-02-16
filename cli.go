package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func cmdList(stdout io.Writer, stderr io.Writer, s *state) int {
	groups := make(map[string]interface{}, 0)
	for _, res := range s.resources() {
		for _, grp := range res.Groups() {
			tmpGroup := []string{}

			_, ok := groups[grp]
			if ok {
				tmpGroup = groups[grp].([]string)
			}

			tmpGroup = append(tmpGroup, res.Address())
			groups[grp] = tmpGroup
		}
	}

	groups["all"] = make(map[string]string, 0)
	for _, out := range s.outputs() {
		groups["all"].(map[string]string)[out.keyName] = out.value
	}

	return output(stdout, stderr, groups)
}

func cmdHost(stdout io.Writer, stderr io.Writer, s *state, hostname string) int {
	for _, res := range s.resources() {
		if hostname == res.Address() {
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
