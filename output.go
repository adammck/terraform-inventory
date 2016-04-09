package main

import (
	"fmt"
)

type Output struct {

	// The keyName and value of the output
	keyName string
	value   string
}

func NewOutput(keyName string, value string) (*Output, error) {

	// TODO: Warn instead of silently ignore error?
	if len(keyName) == 0 {
		return nil, fmt.Errorf("couldn't parse keyName: %s", keyName)
	}

	return &Output{
		keyName: keyName,
		value:   value,
	}, nil
}
