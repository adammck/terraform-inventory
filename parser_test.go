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
							"private_ip": "10.0.0.2"
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
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	exp := state{
		Modules: []moduleState{
			moduleState{
				Resources: map[string]resourceState{
					"aws_instance.one": resourceState{
						Type: "aws_instance",
						Primary: instanceState{
							ID: "i-aaaaaaaa",
							Attributes: map[string]string{
								"id":         "i-aaaaaaaa",
								"private_ip": "10.0.0.1",
							},
						},
					},
					"aws_instance.two": resourceState{
						Type: "aws_instance",
						Primary: instanceState{
							ID: "i-bbbbbbbb",
							Attributes: map[string]string{
								"id":         "i-bbbbbbbb",
								"private_ip": "10.0.0.2",
							},
						},
					},
					"aws_security_group.example": resourceState{
						Type: "aws_security_group",
						Primary: instanceState{
							ID: "sg-cccccccc",
							Attributes: map[string]string{
								"id":          "sg-cccccccc",
								"description": "Whatever",
							},
						},
					},
					"digitalocean_droplet.three": resourceState{
						Type: "digitalocean_droplet",
						Primary: instanceState{
							ID: "ddddddd",
							Attributes: map[string]string{
								"id":           "ddddddd",
								"ipv4_address": "192.168.0.3",
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, exp, s)
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
