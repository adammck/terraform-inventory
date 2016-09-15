package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type allGroup struct {
	Hosts []string          `json:"hosts"`
	Vars  map[string]string `json:"vars"`
}

func gatherResources(s *state) map[string]interface{} {
	groups := make(map[string]interface{}, 0)
	all_group := allGroup{Vars: make(map[string]string)}
	groups["all"] = &all_group

	for _, res := range s.resources() {
		for _, grp := range res.Groups() {

			_, ok := groups[grp]
			if !ok {
				groups[grp] = []string{}
			}

			groups[grp] = append(groups[grp].([]string), res.Address())
			all_group.Hosts = append(all_group.Hosts, res.Address())
		}
	}

	if len(s.outputs()) > 0 {
		for _, out := range s.outputs() {
			all_group.Vars[out.keyName] = out.value.(string)
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

		switch grp := res.(type) {
			case []string:
				writeLn("["+group+"]", stdout, stderr)
			    for _, item := range grp {
					writeLn(item, stdout, stderr)
				}

			case *allGroup:
				writeLn("["+group+"]", stdout, stderr)
			    for _, item := range grp.Hosts {
			    	writeLn(item, stdout, stderr)
				}
				writeLn("", stdout, stderr)
				writeLn("["+group+":vars]", stdout, stderr)
			    for key, item := range grp.Vars {
			    	writeLn(key+"="+item, stdout, stderr)
				}
		}

		writeLn("", stdout, stderr)
	}

	return 0
}

func writeLn(str string, stdout io.Writer, stderr io.Writer) {
	_, err := io.WriteString(stdout, str + "\n")
	checkErr(err, stderr)
}

func checkErr(err error, stderr io.Writer) int {
	if err != nil {
		fmt.Fprintf(stderr, "Error writing Invetory: %s\n", err)
		return 1
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
