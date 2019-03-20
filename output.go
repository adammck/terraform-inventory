package main

import (
	"fmt"
	"strings"
)

type Output struct {

	// The keyName and value of the output
	keyName    string
	value      interface{}
	modulePath string
}

func NewOutput(keyName string, value interface{}, path []string) (*Output, error) {

	// TODO: Warn instead of silently ignore error?
	if len(keyName) == 0 {
		return nil, fmt.Errorf("couldn't parse keyName: %s", keyName)
	}

	return &Output{
		keyName:    keyName,
		value:      value,
		modulePath: strings.Join(path, "."),
	}, nil
}
