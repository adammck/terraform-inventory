package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// TerraformVersion defines which version of Terraform state applies
type TerraformVersion int

const (
	// TerraformVersionUnknown means unknown version
	TerraformVersionUnknown TerraformVersion = 0

	// TerraformVersionPre0dot12 means < 0.12
	TerraformVersionPre0dot12 TerraformVersion = 1

	// TerraformVersion0dot12 means >= 0.12
	TerraformVersion0dot12 TerraformVersion = 2
)

type stateAnyTerraformVersion struct {
	StatePre0dot12   state
	State0dot12      stateTerraform0dot12
	TerraformVersion TerraformVersion
}

// Terraform < v0.12
type state struct {
	Modules []moduleState `json:"modules"`
}
type moduleState struct {
	Path           []string                 `json:"path"`
	ResourceStates map[string]resourceState `json:"resources"`
	Outputs        map[string]interface{}   `json:"outputs"`
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

// Terraform <= v0.12
type stateTerraform0dot12 struct {
	Values valuesStateTerraform0dot12 `json:"values"`
}
type valuesStateTerraform0dot12 struct {
	RootModule *moduleStateTerraform0dot12 `json:"root_module"`
	Outputs    map[string]interface{}      `json:"outputs"`
}
type moduleStateTerraform0dot12 struct {
	ResourceStates []resourceStateTerraform0dot12 `json:"resources"`
	ChildModules   []moduleStateTerraform0dot12   `json:"child_modules"`
	Address        string                         `json:"address"` // empty for root module, else e.g. `module.mymodulename`
}
type resourceStateTerraform0dot12 struct {
	Address   string                 `json:"address"`
	Index     *interface{}           `json:"index"` // only set by Terraform for counted resources
	Name      string                 `json:"name"`
	RawValues map[string]interface{} `json:"values"`
	Type      string                 `json:"type"`
}

// read populates the state object from a statefile.
func (s *stateAnyTerraformVersion) read(stateFile io.Reader) error {
	s.TerraformVersion = TerraformVersionUnknown

	b, readErr := ioutil.ReadAll(stateFile)
	if readErr != nil {
		return readErr
	}

	err0dot12 := json.Unmarshal(b, &(*s).State0dot12)
	if err0dot12 == nil && s.State0dot12.Values.RootModule != nil {
		s.TerraformVersion = TerraformVersion0dot12
	} else {
		errPre0dot12 := json.Unmarshal(b, &(*s).StatePre0dot12)
		if errPre0dot12 == nil && s.StatePre0dot12.Modules != nil {
			s.TerraformVersion = TerraformVersionPre0dot12
		} else {
			return fmt.Errorf("0.12 format error: %v; pre-0.12 format error: %v (nil error means no content/modules found in the respective format)", err0dot12, errPre0dot12)
		}
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

// outputs returns a slice of the Outputs found in the statefile.
func (s *stateTerraform0dot12) outputs() []*Output {
	inst := make([]*Output, 0)

	for k, v := range s.Values.Outputs {
		var o *Output
		switch v := v.(type) {
		case map[string]interface{}:
			o, _ = NewOutput(k, v["value"])
		default: // not expected
			o, _ = NewOutput(k, "<error>")
		}

		inst = append(inst, o)
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

// map of resource ID -> resource Name
func (s *stateTerraform0dot12) mapResourceIDNames() map[string]string {
	t := map[string]string{}

	for _, module := range s.getAllModules() {
		for _, resourceState := range module.ResourceStates {
			id, typeOk := resourceState.RawValues["id"].(string)
			if typeOk && id != "" && resourceState.Name != "" {
				k := strings.ToLower(id)

				if val, ok := resourceState.RawValues["category_id"]; ok && resourceState.Type == "vsphere_tag" {
					if categoryID, typeOk := val.(string); typeOk {
						if categoryName := s.getResourceIDName(categoryID); categoryName != "" {
							t[k] = fmt.Sprintf("%s_%s", s.getResourceIDName(categoryID), resourceState.Name)
							continue
						}
					}
				}

				t[k] = resourceState.Name
			}
		}
	}

	return t
}

func (s *stateTerraform0dot12) getResourceIDName(matchingID string) string {
	for _, module := range s.getAllModules() {
		for _, resourceState := range module.ResourceStates {
			id, typeOk := resourceState.RawValues["id"].(string)
			if typeOk && id == matchingID {
				return resourceState.Name
			}
		}
	}

	return ""
}

func (s *stateTerraform0dot12) getAllModules() []*moduleStateTerraform0dot12 {
	var allModules []*moduleStateTerraform0dot12
	allModules = append(allModules, s.Values.RootModule)
	addChildModules(&allModules, s.Values.RootModule)
	return allModules
}

// recursively adds all child modules to the slice
func addChildModules(out *[]*moduleStateTerraform0dot12, from *moduleStateTerraform0dot12) {
	for i := range from.ChildModules {
		addChildModules(out, &from.ChildModules[i])
		*out = append(*out, &from.ChildModules[i])
	}
}

// resources returns a slice of the Resources found in the statefile.
func (s *stateAnyTerraformVersion) resources() []*Resource {
	switch s.TerraformVersion {
	case TerraformVersionPre0dot12:
		return s.StatePre0dot12.resources()
	case TerraformVersion0dot12:
		return s.State0dot12.resources()
	case TerraformVersionUnknown:
	}
	panic("Unimplemented Terraform version enum")
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
				asJSON, _ := json.Marshal(m.ResourceStates[k])
				fmt.Fprintf(os.Stderr, "Warning: failed to parse resource %s (%v)\n", asJSON, err)
				continue
			}
			if r.IsSupported() {
				inst = append(inst, r)
			}
		}
	}

	return inst
}

func encodeTerraform0Dot12ValuesAsAttributes(rawValues *map[string]interface{}) map[string]string {
	ret := make(map[string]string)
	for k, v := range *rawValues {
		switch v := v.(type) {
		case map[string]interface{}:
			ret[k+".#"] = strconv.Itoa(len(v))
			for kk, vv := range v {
				if str, typeOk := vv.(string); typeOk {
					ret[k+"."+kk] = str
				} else {
					ret[k+"."+kk] = "<error>"
				}
			}
		case []interface{}:
			ret[k+".#"] = strconv.Itoa(len(v))
			for kk, vv := range v {
				switch o := vv.(type) {
				case string:
					ret[k+"."+strconv.Itoa(kk)] = o
				case map[string]interface{}:
					for kkk, vvv := range o {
						if str, typeOk := vvv.(string); typeOk {
							ret[k+"."+strconv.Itoa(kk)+"."+kkk] = str
						} else {
							ret[k+"."+strconv.Itoa(kk)+"."+kkk] = "<error>"
						}
					}
				default:
					ret[k+"."+strconv.Itoa(kk)] = "<error>"
				}
			}
		case string:
			ret[k] = v
		default:
			ret[k] = "<error>"
		}
	}
	return ret
}

// resources returns a slice of the Resources found in the statefile.
func (s *stateTerraform0dot12) resources() []*Resource {
	inst := make([]*Resource, 0)

	for _, module := range s.getAllModules() {
		for _, rs := range module.ResourceStates {
			id, typeOk := rs.RawValues["id"].(string)
			if !typeOk {
				continue
			}

			if strings.HasPrefix(rs.Address, "data.") {
				// This does not represent a host (e.g. AWS AMI)
				continue
			}

			modulePrefix := ""
			if module.Address != "" {
				modulePrefix = strings.Replace(module.Address, ".", "_", -1) + "_"
			}
			resourceKeyName := rs.Type + "." + modulePrefix + rs.Name
			if rs.Index != nil {
				i := *rs.Index
				switch v := i.(type) {
				case int:
					resourceKeyName += "." + strconv.Itoa(v)
				case float64:
					resourceKeyName += "." + strconv.Itoa(int(v))
				case string:
					resourceKeyName += "." + strings.Replace(v, ".", "_", -1)
				default:
					fmt.Fprintf(os.Stderr, "Warning: unknown index type %v\n", v)
				}
			}

			// Terraform stores resources in a name->map map, but we need the name to
			// decide which groups to include the resource in. So wrap it in a higher-
			// level object with both properties.
			//
			// Convert to the pre-0.12 structure for backwards compatibility of code.
			r, err := NewResource(resourceKeyName, resourceState{
				Type: rs.Type,
				Primary: instanceState{
					ID:         id,
					Attributes: encodeTerraform0Dot12ValuesAsAttributes(&rs.RawValues),
				},
			})
			if err != nil {
				asJSON, _ := json.Marshal(rs)
				fmt.Fprintf(os.Stderr, "Warning: failed to parse resource %s (%v)\n", asJSON, err)
				continue
			}
			if r.IsSupported() {
				inst = append(inst, r)
			}
		}
	}

	sort.Slice(inst, func(i, j int) bool {
		return inst[i].baseName < inst[j].baseName
	})

	return inst
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
		i++
	}

	sort.Strings(keys)
	return keys
}
