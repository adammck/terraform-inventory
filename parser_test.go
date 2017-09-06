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
			"outputs": {
					"olddatacenter": "<0.7_format",
					"datacenter": {
						"sensitive": false,
						"type": "string",
						"value": "mydc"
					},
					"ids": {
						"type": "list",
						"value": [1, 2, 3, 4]
					},
					"map": {
						"type": "map",
						"value": {
							"key": "value"
						}
					}
			},
			"resources": {
				"aws_instance.one.0": {
					"type": "aws_instance",
					"primary": {
						"id": "i-aaaaaaaa",
						"attributes": {
							"id": "i-aaaaaaaa",
							"private_ip": "10.0.0.1",
							"tags.%": "1",
							"tags.Role": "Web"
						}
					}
				},
				"aws_instance.dup.0": {
					"type": "aws_instance",
					"primary": {
						"id": "i-aaaaaaaa",
						"attributes": {
							"id": "i-aaaaaaaa",
							"private_ip": "10.0.0.1",
							"tags.#": "1",
							"tags.Role": "Web"
						}
					}
				},
				"aws_instance.one.1": {
					"type": "aws_instance",
					"primary": {
						"id": "i-a1a1a1a1",
						"attributes": {
							"id": "i-a1a1a1a1",
							"private_ip": "10.0.1.1"
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
							"ipv4_address": "192.168.0.3",
							"tags.#": "2",
							"tags.1": "staging",
							"tags.2": "webserver"
						}
					}
				},
				"cloudstack_instance.four": {
					"type": "cloudstack_instance",
					"primary": {
						"id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
						"attributes": {
							"id": "aaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
							"ipaddress": "10.2.1.5",
							"name": "terraform-inventory-4",
							"service_offering": "small",
							"template": "centos-7-0-x64",
							"zone": "nyc2"
						}
					}
				},
				"vsphere_virtual_machine.five": {
					"type": "vsphere_virtual_machine",
					"primary": {
						"id": "aaaaaaaa",
						"attributes": {
							"custom_configuration_parameters.%": "1",
							"custom_configuration_parameters.role": "rrrrrrrr",
							"datacenter": "dddddddd",
							"host": "hhhhhhhh",
							"id": "aaaaaaaa",
							"image": "Ubunutu 14.04 LTS",
							"network_interface.0.ipv4_address": "10.20.30.40",
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
							"access_ip_v6": "",
							"metadata.status": "superServer",
							"metadata.#": "very bad",
							"metadata_toes": "faada2142412jhb1j2"
						}
					}
				},
				"softlayer_virtual_guest.seven": {
					"type": "softlayer_virtual_guest",
					"primary": {
						"id": "12345678",
						"attributes": {
							"id":"12345678",
							"ipv4_address_private":"10.0.0.7",
							"ipv4_address":""
						}
					}
				},
				"google_compute_instance.eight": {
					"type": "123456789",
					"primary": {
						"id": "123456789",
						"attributes": {
							"network_interface.0.access_config.0.assigned_nat_ip": "10.0.0.8",
							"network_interface.0.access_config.0.nat_ip": "10.0.0.8",
							"network_interface.0.address": "10.0.0.8",
							"tags.#": "1",
							"tags.1201918879": "database",
							"tags_fingerprint": "AqbISNuzJIs=",
							"zone": "europe-west1-d"
						}
					}
				},
				"exoscale_compute.nine": {
					"type": "exoscale_compute",
					"depends_on": [
						"x",
						"y"
					],
					"primary": {
						"id": "123456789",
						"attributes": {
							"diskSize": "10",
							"id": "123456789",
							"keypair": "kp",
							"name": "xyz",
							"networks.#": "1",
							"networks.0.%": "5",
							"networks.0.default": "true",
							"networks.0.ip4address": "10.0.0.9",
							"networks.0.ip6address": "",
							"networks.0.networkname": "guestNetworkForBasicZone",
							"networks.0.type": "Shared",
							"securitygroups.#": "1",
							"securitygroups.0": "xyz",
							"size": "zzz",
							"state": "Running",
							"template": "Linux CoreOS stable 1298 64-bit",
							"userdata": "",
							"zone": "ch-gva-2"
						},
						"meta": {},
						"tainted": false
					},
					"deposed": [],
					"provider": ""
				}
			}
		}
	]
}
`

const expectedListOutput = `
{
	"all":	 {
		"hosts": [
			"10.0.0.1",
			"10.0.0.7",
			"10.0.0.8",
			"10.0.0.9",
			"10.0.1.1",
			"10.120.0.226",
			"10.2.1.5",
			"10.20.30.40",
			"192.168.0.3",
			"50.0.0.1"
		]
	},
	"one":   ["10.0.0.1", "10.0.1.1"],
	"dup":   ["10.0.0.1"],
	"two":   ["50.0.0.1"],
	"three": ["192.168.0.3"],
	"four":  ["10.2.1.5"],
	"five":  ["10.20.30.40"],
	"six":   ["10.120.0.226"],
	"seven": ["10.0.0.7"],
	"eight": ["10.0.0.8"],
	"nine": ["10.0.0.9"],

	"one.0":   ["10.0.0.1"],
	"dup.0":   ["10.0.0.1"],
	"one.1":   ["10.0.1.1"],
	"two.0":   ["50.0.0.1"],
	"three.0": ["192.168.0.3"],
	"four.0":  ["10.2.1.5"],
	"five.0":  ["10.20.30.40"],
	"six.0":   ["10.120.0.226"],
	"seven.0": ["10.0.0.7"],
	"eight.0": ["10.0.0.8"],
	"nine.0":  ["10.0.0.9"],

	"type_aws_instance":                  ["10.0.0.1", "10.0.1.1", "50.0.0.1"],
	"type_digitalocean_droplet":          ["192.168.0.3"],
	"type_cloudstack_instance":           ["10.2.1.5"],
	"type_vsphere_virtual_machine":       ["10.20.30.40"],
	"type_openstack_compute_instance_v2": ["10.120.0.226"],
	"type_softlayer_virtual_guest":       ["10.0.0.7"],
	"type_exoscale_compute":              ["10.0.0.9"],
	"type_google_compute_instance":       ["10.0.0.8"],

	"role_rrrrrrrr": ["10.20.30.40"],
	"role_web": ["10.0.0.1"],
	"webserver": ["192.168.0.3"],
	"staging": ["192.168.0.3"],
	"status_superserver": ["10.120.0.226"],
	"database": ["10.0.0.8"],

	"_meta": {"hostvars": {}}
}
`

const expectedHostOneOutput = `
{
	"id":"i-aaaaaaaa",
	"private_ip":"10.0.0.1",
	"tags.#": "1",
	"tags.Role": "Web"
}
`

func TestListCommand(t *testing.T) {
	var s state
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.NoError(t, err)

	// Decode expectation as JSON
	var exp interface{}
	err = json.Unmarshal([]byte(expectedListOutput), &exp)
	assert.NoError(t, err)

	// Run the command, capture the output
	var stdout, stderr bytes.Buffer
	exitCode := cmdList(&stdout, &stderr, &s)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "", stderr.String())

	// Decode the output to compare
	var act interface{}
	err = json.Unmarshal([]byte(stdout.String()), &act)
	assert.NoError(t, err)

	assert.Equal(t, exp, act)
}

func TestHostCommand(t *testing.T) {
	var s state
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.NoError(t, err)

	// Decode expectation as JSON
	var exp interface{}
	err = json.Unmarshal([]byte(expectedHostOneOutput), &exp)
	assert.NoError(t, err)

	// Run the command, capture the output
	var stdout, stderr bytes.Buffer
	exitCode := cmdHost(&stdout, &stderr, &s, "10.0.0.1")
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "", stderr.String())

	// Decode the output to compare
	var act interface{}
	err = json.Unmarshal([]byte(stdout.String()), &act)
	assert.NoError(t, err)

	assert.Equal(t, exp, act)
}
