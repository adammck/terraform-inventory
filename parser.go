package main

import (
	"io"
	"io/ioutil"
	"encoding/json"
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

// hosts returns a map of name to instanceState, for each of the aws_instance
// resources found in the statefile.
func (s *state) instances() map[string]instanceState {
	inst := make(map[string]instanceState)

	for _, m := range s.Modules {
		for k, r := range m.Resources {
			if r.Type == "aws_instance" {
				inst[k] = r.Primary
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

type instanceState struct {
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
