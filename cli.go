package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

type counterSorter struct {
	resources []*Resource
}

func (cs counterSorter) Len() int {
	return len(cs.resources)
}

func (cs counterSorter) Swap(i, j int) {
	cs.resources[i], cs.resources[j] = cs.resources[j], cs.resources[i]
}

func (cs counterSorter) Less(i, j int) bool {
	return cs.resources[i].counter < cs.resources[j].counter
}

type allGroup struct {
	Hosts []string               `json:"hosts"`
	Vars  map[string]interface{} `json:"vars"`
}

func appendUniq(strs []string, item string) []string {
	if len(strs) == 0 {
		strs = append(strs, item)
		return strs
	}
	sort.Strings(strs)
	i := sort.SearchStrings(strs, item)
	if i == len(strs) || (i < len(strs) && strs[i] != item) {
		strs = append(strs, item)
	}
	return strs
}

func gatherResources(s *state) map[string]interface{} {
	outputGroups := make(map[string]interface{})

	all := &allGroup{Hosts: make([]string, 0), Vars: make(map[string]interface{})}
	types := make(map[string][]string)
	individual := make(map[string][]string)
	ordered := make(map[string][]string)
	tags := make(map[string][]string)

	unsortedOrdered := make(map[string][]*Resource)

	resourceIDNames := s.mapResourceIDNames()
	for _, res := range s.resources() {
		// place in list of all resources
		all.Hosts = appendUniq(all.Hosts, res.Hostname())

		// place in list of resource types
		tp := fmt.Sprintf("type_%s", res.resourceType)
		types[tp] = appendUniq(types[tp], res.Hostname())

		unsortedOrdered[res.baseName] = append(unsortedOrdered[res.baseName], res)

		// store as invdividual host (eg. <name>.<count>)
		invdName := fmt.Sprintf("%s.%d", res.baseName, res.counter)
		if old, exists := individual[invdName]; exists {
			fmt.Fprintf(os.Stderr, "overwriting already existing individual key %s, old: %v, new: %v", invdName, old, res.Hostname())
		}
		individual[invdName] = []string{res.Hostname()}

		// inventorize tags
		for k, v := range res.Tags() {
			// Valueless
			tag := k
			if v != "" {
				tag = fmt.Sprintf("%s_%s", k, v)
			}
			// if v is a resource ID, then tag should be resource name
			if _, exists := resourceIDNames[v]; exists {
				tag = resourceIDNames[v]
			}
			tags[tag] = appendUniq(tags[tag], res.Hostname())
		}
	}

	// inventorize outputs as variables
	if len(s.outputs()) > 0 {
		for _, out := range s.outputs() {
			all.Vars[out.keyName] = out.value
		}
	}

	// sort the ordered groups
	for basename, resources := range unsortedOrdered {
		cs := counterSorter{resources}
		sort.Sort(cs)

		for i := range resources {
			ordered[basename] = append(ordered[basename], resources[i].Hostname())
		}
	}

	outputGroups["all"] = all
	for k, v := range individual {
		if old, exists := outputGroups[k]; exists {
			fmt.Fprintf(os.Stderr, "individual overwriting already existing output with key %s, old: %v, new: %v", k, old, v)
		}
		outputGroups[k] = v
	}
	for k, v := range ordered {
		if old, exists := outputGroups[k]; exists {
			fmt.Fprintf(os.Stderr, "ordered overwriting already existing output with key %s, old: %v, new: %v", k, old, v)
		}
		outputGroups[k] = v
	}
	for k, v := range types {
		if old, exists := outputGroups[k]; exists {
			fmt.Fprintf(os.Stderr, "types overwriting already existing output key %s, old: %v, new: %v", k, old, v)
		}
		outputGroups[k] = v
	}
	for k, v := range tags {
		if old, exists := outputGroups[k]; exists {
			fmt.Fprintf(os.Stderr, "tags overwriting already existing output key %s, old: %v, new: %v", k, old, v)
		}
		outputGroups[k] = v
	}

	return outputGroups
}

func cmdList(stdout io.Writer, stderr io.Writer, s *state) int {
	return output(stdout, stderr, gatherResources(s))
}

func cmdInventory(stdout io.Writer, stderr io.Writer, s *state) int {
	groups := gatherResources(s)
	group_names := []string{}
	for group, _ := range groups {
		group_names = append(group_names, group)
	}
	sort.Strings(group_names)
	for _, group := range group_names {

		switch grp := groups[group].(type) {
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
			vars := []string{}
			for key, _ := range grp.Vars {
				vars = append(vars, key)
			}
			sort.Strings(vars)
			for _, key := range vars {
				jsonItem, _ := json.Marshal(grp.Vars[key])
				itemLn := fmt.Sprintf("%s", string(jsonItem))
				writeLn(key+"="+itemLn, stdout, stderr)
			}
		}

		writeLn("", stdout, stderr)
	}

	return 0
}

func writeLn(str string, stdout io.Writer, stderr io.Writer) {
	_, err := io.WriteString(stdout, str+"\n")
	checkErr(err, stderr)
}

func checkErr(err error, stderr io.Writer) int {
	if err != nil {
		fmt.Fprintf(stderr, "Error writing inventory: %s\n", err)
		return 1
	}
	return 0
}

func cmdHost(stdout io.Writer, stderr io.Writer, s *state, hostname string) int {
	for _, res := range s.resources() {
		if hostname == res.Hostname() {
			attributes := res.Attributes()
			attributes["ansible_host"] = res.Address()
			return output(stdout, stderr, attributes)
		}
	}

	fmt.Fprintf(stdout, "{}")
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
