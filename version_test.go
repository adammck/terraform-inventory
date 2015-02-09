package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionInfo(t *testing.T) {
	assert.Equal(t, "unknown", versionInfo())

	build_version = "vXYZ"
	assert.Equal(t, "vXYZ", versionInfo())
}
