package main

import (
	"fmt"
)

type Output struct {

	// The keyName and value of the output
	keyName string
	value   interface{}
}

func NewOutput(keyName string, value interface{}) (*Output, error) {

	// TODO: Warn instead of silently ignore error?
	if len(keyName) == 0 {
		return nil, fmt.Errorf("couldn't parse keyName: %s", keyName)
	}

	return &Output{
		keyName: keyName,
		value:   value,
	}, nil
}
