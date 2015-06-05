package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const exampleStateFile = `
{
	"version": 1,
	"serial": 1,
	"modules": [
		{
			"path": [
				"root"
			],
			"outputs": {},
			"resources": {
				"aws_instance.one": {
					"type": "aws_instance",
					"primary": {
						"id": "i-aaaaaaaa",
						"attributes": {
							"id": "i-aaaaaaaa",
							"private_ip": "10.0.0.1"
						}
					}
				},
				"aws_instance.two": {
					"type": "aws_instance",
					"primary": {
						"id": "i-bbbbbbbb",
						"attributes": {
							"id": "i-bbbbbbbb",
							"private_ip": "10.0.0.2",
							"public_ip": "50.0.0.1"
						}
					}
				},
				"aws_security_group.example": {
					"type": "aws_security_group",
					"primary": {
						"id": "sg-cccccccc",
						"attributes": {
							"id": "sg-cccccccc",
							"description": "Whatever"
						}
					}
				},
				"digitalocean_droplet.three": {
					"type": "digitalocean_droplet",
					"primary": {
						"id": "ddddddd",
						"attributes": {
							"id": "ddddddd",
							"ipv4_address": "192.168.0.3"
						}
					}
				}
			}
		}
	]
}
`

func TestStateRead(t *testing.T) {
	var s state
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.Nil(t, err)
	assert.Equal(t, "aws_instance", s.Modules[0].Resources["aws_instance.one"].Type)
}

func TestResources(t *testing.T) {
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	inst := s.resources()
	assert.Equal(t, 3, len(inst))
	assert.Equal(t, "aws_instance", inst["one"].Type)
	assert.Equal(t, "aws_instance", inst["two"].Type)
	assert.Equal(t, "digitalocean_droplet", inst["three"].Type)
}

func TestAddress(t *testing.T) {
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	inst := s.resources()
	assert.Equal(t, 3, len(inst))
	assert.Equal(t, "10.0.0.1", inst["one"].Address())
	assert.Equal(t, "50.0.0.1", inst["two"].Address())
	assert.Equal(t, "192.168.0.3", inst["three"].Address())
}

func TestIsSupported(t *testing.T) {
	r := resourceState{
		Type: "something",
	}
	assert.Equal(t, false, r.isSupported())

	r = resourceState{
		Type: "aws_instance",
		Primary: instanceState{
			Attributes: map[string]string{
				"private_ip": "10.0.0.2",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())

	r = resourceState{
		Type: "digitalocean_droplet",
		Primary: instanceState{
			Attributes: map[string]string{
				"ipv4_address": "192.168.0.3",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())
}
