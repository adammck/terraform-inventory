package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
				},
				"cloudstack_instance.four": {
					"type": "cloudstack_instance",
					"primary": {
						"id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"attributes": {
							"id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
							"ipaddress": "10.2.1.5"
						}
					}
				},
				"vsphere_virtual_machine.five": {
					"type": "vsphere_virtual_machine",
					"primary": {
						"id": "aaaaaaaa",
						"attributes": {
							"datacenter": "dddddddd",
							"host": "hhhhhhhh",
							"id": "aaaaaaaa",
							"image": "Ubunutu 14.04 LTS",
							"ip_address": "10.20.30.40",
							"linked_clone": "false",
							"name": "nnnnnn",
							"power_on": "true"
						}
					}
				},
				"openstack_compute_instance_v2.six": {
					"type": "openstack_compute_instance_v2",
					"primary": {
						"id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"attributes": {
							"id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
							"access_ip_v4": "10.120.0.226",
							"access_ip_v6": ""
						}
					}
				}
			}
		}
	]
}
`

const expectedListOutput = `
{
	"one":   ["10.0.0.1"],
	"two":   ["50.0.0.1"],
	"three": ["192.168.0.3"],
	"four":  ["10.2.1.5"],
	"five":  ["10.20.30.40"],
	"six":   ["10.120.0.226"],

	"one.0":   ["10.0.0.1"],
	"two.0":   ["50.0.0.1"],
	"three.0": ["192.168.0.3"],
	"four.0":  ["10.2.1.5"],
	"five.0":  ["10.20.30.40"],
	"six.0":   ["10.120.0.226"]
}
`

func TestIntegration(t *testing.T) {

	var s state
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.Nil(t, err)

	// Run the command, capture the output
	var stdout, stderr bytes.Buffer
	exitCode := cmdList(&stdout, &stderr, &s)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "", stderr.String())

	var exp, act interface{}
	json.Unmarshal([]byte(expectedListOutput), &exp)
	json.Unmarshal([]byte(stdout.String()), &act)
	assert.Equal(t, exp, act)
}

func TestStateRead(t *testing.T) {
	var s state
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.Nil(t, err)
	assert.Equal(t, "aws_instance", s.Modules[0].Resources["aws_instance.one"].Type)
	assert.Equal(t, "aws_instance", s.Modules[0].Resources["aws_instance.two"].Type)
}

func TestResources(t *testing.T) {
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	inst := s.resources()
	assert.Equal(t, 6, len(inst))
	assert.Equal(t, "aws_instance", inst["aws_instance.one"].Type)
	assert.Equal(t, "aws_instance", inst["aws_instance.two"].Type)
	assert.Equal(t, "digitalocean_droplet", inst["digitalocean_droplet.three"].Type)
	assert.Equal(t, "cloudstack_instance", inst["cloudstack_instance.four"].Type)
	assert.Equal(t, "vsphere_virtual_machine", inst["vsphere_virtual_machine.five"].Type)
	assert.Equal(t, "openstack_compute_instance_v2", inst["openstack_compute_instance_v2.six"].Type)
}

func TestAddress(t *testing.T) {
	r := strings.NewReader(exampleStateFile)

	var s state
	err := s.read(r)
	assert.Nil(t, err)

	inst := s.resources()
	assert.Equal(t, 6, len(inst))
	assert.Equal(t, "10.0.0.1", inst["aws_instance.one"].Address())
	assert.Equal(t, "50.0.0.1", inst["aws_instance.two"].Address())
	assert.Equal(t, "192.168.0.3", inst["digitalocean_droplet.three"].Address())
	assert.Equal(t, "10.2.1.5", inst["cloudstack_instance.four"].Address())
	assert.Equal(t, "10.20.30.40", inst["vsphere_virtual_machine.five"].Address())
	assert.Equal(t, "10.120.0.226", inst["openstack_compute_instance_v2.six"].Address())
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

	r = resourceState{
		Type: "cloudstack_instance",
		Primary: instanceState{
			Attributes: map[string]string{
				"ipaddress": "10.2.1.5",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())

	r = resourceState{
		Type: "vsphere_virtual_machine",
		Primary: instanceState{
			Attributes: map[string]string{
				"ip_address": "10.20.30.40",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())

	r = resourceState{
		Type: "openstack_compute_instance_v2",
		Primary: instanceState{
			Attributes: map[string]string{
				"ip_address": "10.120.0.226",
			},
		},
	}
	assert.Equal(t, true, r.isSupported())
}
