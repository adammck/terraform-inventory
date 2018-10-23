package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleStateFileEnvHostname = `
{
	"version": 1,
	"serial": 1,
	"modules": [
		{
			"resources": {
				"libvirt_domain.fourteen": {
					"type": "libvirt_domain",
					"primary": {
						"id": "824c29be-2164-44c8-83e0-787705571d95",
						"attributes": {
							"name": "fourteen",
							"network_interface.#": "1",
							"network_interface.0.addresses.#": "1",
							"network_interface.0.addresses.0": "192.168.102.14",
							"network_interface.0.mac": "96:EE:4D:BD:B2:45"
						}
					}
				}
			}
		}
	]
}`

const expectedListOutputEnvHostname = `
{
	"all":	 {
		"hosts": [
			"fourteen"
		],
		"vars": {
		}
	},
	"fourteen":	 ["fourteen"],
	"fourteen.0":	 ["fourteen"],
	"type_libvirt_domain": ["fourteen"]
}`

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
				},
				"triton_machine.ten": {
					"type": "triton_machine",
					"depends_on": [],
					"primary": {
						"id": "123456789",
						"attributes": {
							"administrator_pw": "",
							"cloud_config": "",
							"dataset": "dset1",
							"disk": "25600",
							"domain_names.#": "",
							"domain_names.0": "",
							"domain_names.1": "",
							"domain_names.2": "",
							"domain_names.3": "",
							"firewall_enabled": "true",
							"id": "123456789",
							"image": "",
							"ips.#": "1",
							"ips.0": "10.0.0.10",
							"memory": "1024",
							"metadata.%": "0",
							"name": "triton_ten",
							"networks.#": "2",
							"networks.0": "1",
							"networks.1": "2",
							"nic.#": "1",
							"nic.1.gateway": "",
							"nic.1.ip": "10.0.0.10",
							"nic.1.mac": "",
							"nic.1.netmask": "",
							"nic.1.network": "",
							"nic.1.primary": "true",
							"nic.1.state": "running",
							"package": "g4-highcpu-1G",
							"primaryip": "10.0.0.10",
							"tags.%": "1",
							"tags.Role": "test",
							"type": "smartmachine",
							"user_data": ""
						},
						"meta": {},
						"tainted": false
					},
					"deposed": [],
					"provider": ""
				},
				"scaleway_server.eleven": {
					"type": "scaleway_server",
					"depends_on": [],
					"primary": {
						"id": "490c369b-e062-4951-b1c5-f9a8ccee8a38",
						"attributes": {
							"enable_ipv6": "true",
							"id": "490c369b-e062-4951-b1c5-f9a8ccee8a38",
							"image": "ab8fbe9e-b13c-46a5-9139-ae7ae73569f0",
							"name": "eleven",
							"private_ip": "11.0.0.11",
							"public_ip": "10.0.0.11",
							"public_ipv6": "2001:bc8:4400:2500::e:800",
							"security_group": "92a62362-72ab-4864-a94e-f02557153218",
							"state": "running",
							"state_detail": "booted",
							"tags.#": "1",
							"tags.0": "scw_test",
							"type": "VC1S"
						},
						"meta": {},
						"tainted": false
					},
					"deposed": [],
					"provider": "provider.scaleway"
				},
				"vsphere_virtual_machine.twelve": {
					"type": "vsphere_virtual_machine",
					"primary": {
						"id": "422cfa4a-c6bb-3405-0335-2d9b2034405f",
						"attributes": {
							"default_ip_address": "10.20.30.50",
							"tags.#": "1",
							"tags.1357913579": "urn:vmomi:InventoryServiceTag:00000000-0001-4957-81fa-1234567890ab:GLOBAL"
						}
					}
				},
				"data.vsphere_tag.testTag1": {
					"type": "vsphere_tag",
					"primary": {
						"id": "urn:vmomi:InventoryServiceTag:00000000-0001-4957-81fa-1234567890ab:GLOBAL",
						"attributes": {
							"name": "testTag1"
						}
					}
				},
				"packet_device.thirteen": {
					"type": "packet_device",
					"depends_on": [],
					"primary": {
						"id": "e35816e2-b9b4-4ef3-9317-a32b98f6cb44",
						"attributes": {
							"billing_cycle": "hourly",
							"created": "2018-04-02T14:52:34Z",
							"facility": "ewr1",
							"hostname": "sa-test-1",
							"id": "e35816e2-b9b4-4ef3-9317-a32b98f6cb44",
							"locked": "false",
							"network.#": "3",
							"network.0.address": "10.0.0.13",
							"network.0.cidr": "31",
							"network.0.family": "4",
							"network.0.gateway": "10.0.0.254",
							"network.0.public": "true",
							"operating_system": "ubuntu_16_04",
							"plan": "baremetal_0",
							"project_id": "123456d5-087a-4976-877a-45b86584b786",
							"state": "active",
							"tags.#": "0",
							"updated": "2018-04-02T14:57:13Z"
						},
						"meta": {},
						"tainted": false
					},
					"deposed": [],
					"provider": ""
				},
				"libvirt_domain.fourteen": {
					"type": "libvirt_domain",
					"primary": {
						"id": "824c29be-2164-44c8-83e0-787705571d95",
						"attributes": {
							"network_interface.#": "1",
							"network_interface.0.addresses.#": "1",
							"network_interface.0.addresses.0": "192.168.102.14",
							"network_interface.0.mac": "96:EE:4D:BD:B2:45"
						}
					}
				},
				"profitbricks_server.sixteen": {
					"type": "profitbricks_server",
					"primary": {
						"id": "12345678",
						"attributes": {
							"primary_ip": "10.0.0.16"
						}
					}
				},
				"aws_spot_instance_request.seventeen": {
					"type": "aws_spot_instance_request",
					"primary": {
						"id": "i-a1a1a1a1",
						"attributes": {
							"id": "sir-a1a1a1a1",
							"public_ip": "50.0.0.17",
							"tags.%": "1",
							"tags.Role": "worker"
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
	"all":	 {
		"hosts": [
			"10.0.0.1",
			"10.0.0.10",
			"10.0.0.11",
			"10.0.0.13",
			"10.0.0.16",
			"10.0.0.7",
			"10.0.0.8",
			"10.0.0.9",
			"10.0.1.1",
			"10.120.0.226",
			"10.2.1.5",
			"10.20.30.40",
			"192.168.0.3",
			"192.168.102.14",
			"50.0.0.1",
			"50.0.0.17",
			"10.20.30.50"
		],
		"vars": {
			"datacenter": "mydc",
			"olddatacenter": "<0.7_format",
			"ids": [1, 2, 3, 4],
			"map": {"key": "value"}
		}
	},
	"one":	 ["10.0.0.1", "10.0.1.1"],
	"dup":	 ["10.0.0.1"],
	"two":	 ["50.0.0.1"],
	"three": ["192.168.0.3"],
	"four":  ["10.2.1.5"],
	"five":  ["10.20.30.40"],
	"six":	 ["10.120.0.226"],
	"seven": ["10.0.0.7"],
	"eight": ["10.0.0.8"],
	"nine": ["10.0.0.9"],
	"ten": ["10.0.0.10"],
	"eleven": ["10.0.0.11"],
	"twelve": ["10.20.30.50"],
	"testTag1": ["10.20.30.50"],
	"thirteen": ["10.0.0.13"],
	"fourteen": ["192.168.102.14"],
	"sixteen": ["10.0.0.16"],
	"seventeen": ["50.0.0.17"],

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
	"ten.0":   ["10.0.0.10"],
	"eleven.0": ["10.0.0.11"],
	"twelve.0": ["10.20.30.50"],
	"thirteen.0": ["10.0.0.13"],
	"fourteen.0": ["192.168.102.14"],
	"sixteen.0": ["10.0.0.16"],
	"seventeen.0": ["50.0.0.17"],

	"type_aws_instance":                  ["10.0.0.1", "10.0.1.1", "50.0.0.1"],
	"type_digitalocean_droplet":          ["192.168.0.3"],
	"type_cloudstack_instance":           ["10.2.1.5"],
	"type_vsphere_virtual_machine":       ["10.20.30.40", "10.20.30.50"],
	"type_openstack_compute_instance_v2": ["10.120.0.226"],
	"type_profitbricks_server":           ["10.0.0.16"],
	"type_softlayer_virtual_guest":       ["10.0.0.7"],
	"type_exoscale_compute":              ["10.0.0.9"],
	"type_google_compute_instance":       ["10.0.0.8"],
	"type_triton_machine":                ["10.0.0.10"],
	"type_scaleway_server":               ["10.0.0.11"],
	"type_packet_device":                 ["10.0.0.13"],
	"type_libvirt_domain":                ["192.168.102.14"],
	"type_aws_spot_instance_request":			["50.0.0.17"],

	"role_rrrrrrrr": ["10.20.30.40"],
	"role_web": ["10.0.0.1"],
	"role_test": ["10.0.0.10"],
	"role_worker": ["50.0.0.17"],
	"webserver": ["192.168.0.3"],
	"staging": ["192.168.0.3"],
	"status_superserver": ["10.120.0.226"],
	"database": ["10.0.0.8"],
	"scw_test": ["10.0.0.11"]
}
`

const expectedInventoryOutput = `[all]
10.0.0.1
10.0.0.10
10.0.0.11
10.0.0.13
10.0.0.16
10.0.0.7
10.0.0.8
10.0.0.9
10.0.1.1
10.120.0.226
10.2.1.5
10.20.30.40
192.168.0.3
192.168.102.14
50.0.0.1
50.0.0.17
10.20.30.50

[all:vars]
datacenter="mydc"
ids=[1,2,3,4]
map={"key":"value"}
olddatacenter="\u003c0.7_format"

[database]
10.0.0.8

[dup]
10.0.0.1

[dup.0]
10.0.0.1

[eight]
10.0.0.8

[eight.0]
10.0.0.8

[eleven]
10.0.0.11

[eleven.0]
10.0.0.11

[five]
10.20.30.40

[five.0]
10.20.30.40

[four]
10.2.1.5

[four.0]
10.2.1.5

[fourteen]
192.168.102.14

[fourteen.0]
192.168.102.14

[nine]
10.0.0.9

[nine.0]
10.0.0.9

[one]
10.0.0.1
10.0.1.1

[one.0]
10.0.0.1

[one.1]
10.0.1.1

[role_rrrrrrrr]
10.20.30.40

[role_test]
10.0.0.10

[role_web]
10.0.0.1

[role_worker]
50.0.0.17

[scw_test]
10.0.0.11

[seven]
10.0.0.7

[seven.0]
10.0.0.7

[seventeen]
50.0.0.17

[seventeen.0]
50.0.0.17

[six]
10.120.0.226

[six.0]
10.120.0.226

[sixteen]
10.0.0.16

[sixteen.0]
10.0.0.16

[staging]
192.168.0.3

[status_superserver]
10.120.0.226

[ten]
10.0.0.10

[ten.0]
10.0.0.10

[testTag1]
10.20.30.50

[thirteen]
10.0.0.13

[thirteen.0]
10.0.0.13

[three]
192.168.0.3

[three.0]
192.168.0.3

[twelve]
10.20.30.50

[twelve.0]
10.20.30.50

[two]
50.0.0.1

[two.0]
50.0.0.1

[type_aws_instance]
10.0.0.1
10.0.1.1
50.0.0.1

[type_aws_spot_instance_request]
50.0.0.17

[type_cloudstack_instance]
10.2.1.5

[type_digitalocean_droplet]
192.168.0.3

[type_exoscale_compute]
10.0.0.9

[type_google_compute_instance]
10.0.0.8

[type_libvirt_domain]
192.168.102.14

[type_openstack_compute_instance_v2]
10.120.0.226

[type_packet_device]
10.0.0.13

[type_profitbricks_server]
10.0.0.16

[type_scaleway_server]
10.0.0.11

[type_softlayer_virtual_guest]
10.0.0.7

[type_triton_machine]
10.0.0.10

[type_vsphere_virtual_machine]
10.20.30.40
10.20.30.50

[webserver]
192.168.0.3

`

const expectedHostOneOutput = `
{
	"ansible_host": "10.0.0.1",
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

func TestListCommandEnvHostname(t *testing.T) {
	var s state
	r := strings.NewReader(exampleStateFileEnvHostname)
	err := s.read(r)
	assert.NoError(t, err)

	// Decode expectation as JSON
	var exp interface{}
	err = json.Unmarshal([]byte(expectedListOutputEnvHostname), &exp)
	assert.NoError(t, err)

	// Run the command, capture the output
	var stdout, stderr bytes.Buffer
	os.Setenv("TF_HOSTNAME_KEY_NAME", "name")
	exitCode := cmdList(&stdout, &stderr, &s)
	os.Unsetenv("TF_HOSTNAME_KEY_NAME")
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

func TestInventoryCommand(t *testing.T) {
	var s state
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.NoError(t, err)

	// Run the command, capture the output
	var stdout, stderr bytes.Buffer
	exitCode := cmdInventory(&stdout, &stderr, &s)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "", stderr.String())

	assert.Equal(t, expectedInventoryOutput, stdout.String())
}
