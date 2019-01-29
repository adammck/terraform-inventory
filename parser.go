package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type state struct {
	Modules []moduleState `json:"modules"`
}

// read populates the state object from a statefile.
func (s *state) read(stateFile io.Reader) error {

	// read statefile contents
	b, err := ioutil.ReadAll(stateFile)
	if err != nil {
		return err
	}

	// parse into struct
	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	return nil
}

// outputs returns a slice of the Outputs found in the statefile.
func (s *state) outputs() []*Output {
	inst := make([]*Output, 0)

	for _, m := range s.Modules {
		for k, v := range m.Outputs {
			var o *Output
			switch v := v.(type) {
			case map[string]interface{}:
				o, _ = NewOutput(k, v["value"])
			case string:
				o, _ = NewOutput(k, v)
			default:
				o, _ = NewOutput(k, "<error>")
			}

			inst = append(inst, o)
		}
	}

	return inst
}

// map of resource ID -> resource Name
func (s *state) mapResourceIDNames() map[string]string {
	t := map[string]string{}

	for _, m := range s.Modules {
		for _, k := range m.resourceKeys() {
			if m.ResourceStates[k].Primary.ID != "" && m.ResourceStates[k].Primary.Attributes["name"] != "" {
				kk := strings.ToLower(m.ResourceStates[k].Primary.ID)
				t[kk] = m.ResourceStates[k].Primary.Attributes["name"]
			}
		}
	}
	return t
}

// resources returns a slice of the Resources found in the statefile.
func (s *state) resources() []*Resource {
	inst := make([]*Resource, 0)

	for _, m := range s.Modules {
		for _, k := range m.resourceKeys() {
			if strings.HasPrefix(k, "data.") {
				// This does not represent a host (e.g. AWS AMI)
				continue
			}

			// If a module is used, the resource key may not be unique, for instance:
			//
			// The module cannot use dynamic resource naming and thus has to use some hardcoded name:
			//
			//     resource "aws_instance" "host" { ... }
			//
			// The main file then uses the module twice:
			//
			//     module "application1" { source = "./modules/mymodulename" }
			//     module "application2" { source = "./modules/mymodulename" }
			//
			// Avoid key clashes by prepending module name to the key. If path is ["root"], don't
			// prepend anything.
			//
			// In the above example: `aws_instance.host` -> `aws_instance.application1_host`
			fullKey := k
			resourceNameIndex := strings.Index(fullKey, ".") + 1
			if len(m.Path) > 1 && resourceNameIndex > 0 {
				for i := len(m.Path) - 1; i >= 1; i-- {
					fullKey = fullKey[:resourceNameIndex] + strings.Replace(m.Path[i], ".", "_", -1) + "_" + fullKey[resourceNameIndex:]
				}
			}

			// Terraform stores resources in a name->map map, but we need the name to
			// decide which groups to include the resource in. So wrap it in a higher-
			// level object with both properties.
			r, err := NewResource(fullKey, m.ResourceStates[k])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse resource with key %s: %s\n", k, err)
				continue
			}
			if r.IsSupported() {
				inst = append(inst, r)
			}
		}
	}

	return inst
}

type moduleState struct {
	Path           []string                 `json:"path"`
	ResourceStates map[string]resourceState `json:"resources"`
	Outputs        map[string]interface{}   `json:"outputs"`
}

// resourceKeys returns a sorted slice of the key names of the resources in this
// module. Do this instead of range over ResourceStates, to ensure that the
// output is consistent.
func (ms *moduleState) resourceKeys() []string {
	lk := len(ms.ResourceStates)
	keys := make([]string, lk, lk)
	i := 0

	for k := range ms.ResourceStates {
		keys[i] = k
		i += 1
	}

	sort.Strings(keys)
	return keys
}

type resourceState struct {

	// Populated from statefile
	Type    string        `json:"type"`
	Primary instanceState `json:"primary"`
}

type instanceState struct {
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
