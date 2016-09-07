package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func gatherResources(s *state) map[string]interface{} {
	groups := make(map[string]interface{}, 0)
	for _, res := range s.resources() {
		for _, grp := range res.Groups() {

			_, ok := groups[grp]
			if !ok {
				groups[grp] = []string{}
			}

			groups[grp] = append(groups[grp].([]string), res.Address())
		}
	}

	if len(s.outputs()) > 0 {
		groups["all"] = make(map[string]interface{}, 0)
		for _, out := range s.outputs() {
			groups["all"].(map[string]interface{})[out.keyName] = out.value
		}
	}
	return groups
}

func cmdList(stdout io.Writer, stderr io.Writer, s *state) int {
	return output(stdout, stderr, gatherResources(s))
}

func cmdInventory(stdout io.Writer, stderr io.Writer, s *state) int {
	groups := gatherResources(s)
	for group, res := range groups {

		_, err := io.WriteString(stdout, "["+group+"]\n")
		if err != nil {
			fmt.Fprintf(stderr, "Error writing Invetory: %s\n", err)
			return 1
		}

		for _, ress := range res.([]string) {

			_, err := io.WriteString(stdout, ress+"\n")
			if err != nil {
				fmt.Fprintf(stderr, "Error writing Invetory: %s\n", err)
				return 1
			}
		}

		_, err = io.WriteString(stdout, "\n")
		if err != nil {
			fmt.Fprintf(stderr, "Error writing Invetory: %s\n", err)
			return 1
		}
	}

	return 0
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
