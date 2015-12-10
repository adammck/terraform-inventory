package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

type state struct {
	Modules []moduleState `json:"modules"`
}

// keyNames contains the names of the keys to check for in each resource in the
// state file. This allows us to support multiple types of resource without too
// much fuss.
var keyNames []string
var nameParser *regexp.Regexp

func init() {
	keyNames = []string{
		"ipv4_address", // DO
		"public_ip",    // AWS
		"private_ip",   // AWS
		"ipaddress",    // CS
		"ip_address",   // VMware
		"access_ip_v4", // OPENSTACK
	}

	// type.name.0
	nameParser = regexp.MustCompile(`^(\w+)\.([\w\-]+)(?:\.(\d+))?$`)
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

// resources returns a map of name to resourceState, for any supported resources
// found in the statefile.
func (s *state) resources() map[string]resourceState {
	inst := make(map[string]resourceState)

	for _, m := range s.Modules {
		for k, r := range m.Resources {
			if r.isSupported() {

				_, name, counter := parseName(k)
				//fmt.Println(resType, name, counter)
				r.Name = name
				r.Counter = counter
				inst[k] = r
			}
		}
	}

	return inst
}

func parseName(name string) (string, string, int) {
	m := nameParser.FindStringSubmatch(name)

	// This should not happen unless our regex changes.
	// TODO: Warn instead of silently ignore error?
	if len(m) != 4 {
		//fmt.Printf("len=%d\n", len(m))
		return "", "", 0
	}

	var c int
	var err error
	if m[3] != "" {
		c, err = strconv.Atoi(m[3])
		if err != nil {
			fmt.Printf("err: %s\n", err)
			// ???
		}
	}

	return m[1], m[2], c
}

type moduleState struct {
	Resources map[string]resourceState `json:"resources"`
}

type resourceState struct {

	// Populated from statefile
	Type    string        `json:"type"`
	Primary instanceState `json:"primary"`

	// Extracted from key name, and injected by resources method
	Name    string
	Counter int
}

// isSupported returns true if terraform-inventory supports this resource.
func (s resourceState) isSupported() bool {
	return s.Address() != ""
}

// NameWithCounter returns the resource name with its counter. For resources
// created without a 'count=' attribute, this will always be zero.
func (s resourceState) NameWithCounter() string {
	return fmt.Sprintf("%s.%d", s.Name, s.Counter)
}

// Address returns the IP address of this resource.
func (s resourceState) Address() string {
	for _, key := range keyNames {
		if ip := s.Primary.Attributes[key]; ip != "" {
			return ip
		}
	}

	return ""
}

// Attributes returns a map containing everything we know about this resource.
func (s resourceState) Attributes() map[string]string {
	return s.Primary.Attributes
}

type instanceState struct {
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
