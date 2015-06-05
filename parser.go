package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
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

// resources returns a map of name to resourceState, for any supported resources
// found in the statefile.
func (s *state) resources() map[string]resourceState {
	typeRemover := regexp.MustCompile(`^[\w_]+\.`)
	inst := make(map[string]resourceState)

	for _, m := range s.Modules {
		for k, r := range m.Resources {
			if r.isSupported() {
				name := typeRemover.ReplaceAllString(k, "")
				inst[name] = r
			}
		}
	}

	return inst
}

type moduleState struct {
	Resources map[string]resourceState `json:"resources"`
}

type resourceState struct {
	Type    string        `json:"type"`
	Primary instanceState `json:"primary"`
}

// isSupported returns true if terraform-inventory supports this resource.
func (s *resourceState) isSupported() bool {
	return s.Address() != ""
}

// Address returns the IP address of this resource.
func (s *resourceState) Address() string {
	switch s.Type {
	case "aws_instance":
		return s.Primary.Attributes["private_ip"]

	case "digitalocean_droplet":
		return s.Primary.Attributes["ipv4_address"]

	default:
		return ""
	}
}

// Attributes returns a map containing everything we know about this resource.
func (s *resourceState) Attributes() map[string]string {
	return s.Primary.Attributes
}

type instanceState struct {
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
