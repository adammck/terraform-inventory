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
	"fourteen_0":	 ["fourteen"],
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
						"id": "5730e2c9-765f-46e1-aa77-81c94f56ce5d",
						"attributes": {
							"diskSize": "10",
							"id": "5730e2c9-765f-46e1-aa77-81c94f56ce5d",
							"keypair": "kp",
							"name": "xyz",
							"gateway": "10.0.0.1",
							"ip4": "true",
							"ip6": "false",
							"ip6_address": "",
							"ip6_cidr": "",
							"ip_address": "10.0.0.9",
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
							"tags.%": "1",
							"tags.Role": "nine",
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
				"opentelekomcloud_compute_instance_v2.nineteen": {
					"type": "opentelekomcloud_compute_instance_v2",
					"primary": {
						"id": "00000000-0000-0000-0000-000000000015",
						"attributes": {
							"id": "00000000-0000-0000-0000-000000000015",
							"access_ip_v4": "10.0.0.19",
							"access_ip_v6": "",
							"tag.%": "1",
							"tag.tfinventory": "rocks"
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
				},
				"linode_instance.eighteen": {
					"type": "linode_instance",
					"depends_on": [],
					"primary": {
						"id": "123456789",
						"attributes": {
							"ip_address": "80.80.100.124",
							"private_ip": "true",
							"private_ip_address": "192.168.167.23",
							"tags.#": "0"
						}
					}
                },
				"hcloud_server.twenty": {
					"type": "hcloud_server",
					"depends_on": [],
					"primary": {
						"id": "42",
						"attributes": {
							"backup_window": "",
							"backups": "false",
							"datacenter": "fsn1-dc14",
							"id": "42",
							"image": "1",
							"ipv4_address": "10.0.0.20",
							"keep_disk": "false",
							"labels.%": "1",
							"labels.testlabel": "hcloud_test",
							"location": "fsn1",
							"name": "twenty",
							"server_type": "cx11",
							"ssh_keys.#": "1",
							"ssh_keys.0": "1337",
							"status": "running"
						},
						"meta": {},
						"tainted": false
					},
					"deposed": [],
					"provider": "provider.hcloud"
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
			"10.0.0.19",
			"10.0.0.20",
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
			"80.80.100.124",
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
	"nineteen": ["10.0.0.19"],
	"sixteen": ["10.0.0.16"],
	"seventeen": ["50.0.0.17"],
	"eighteen": ["80.80.100.124"],
	"twenty": ["10.0.0.20"],

	"one_0":   ["10.0.0.1"],
	"dup_0":   ["10.0.0.1"],
	"one_1":   ["10.0.1.1"],
	"two_0":   ["50.0.0.1"],
	"three_0": ["192.168.0.3"],
	"four_0":  ["10.2.1.5"],
	"five_0":  ["10.20.30.40"],
	"six_0":   ["10.120.0.226"],
	"seven_0": ["10.0.0.7"],
	"eight_0": ["10.0.0.8"],
	"nine_0":  ["10.0.0.9"],
	"ten_0":   ["10.0.0.10"],
	"eleven_0": ["10.0.0.11"],
	"twelve_0": ["10.20.30.50"],
	"thirteen_0": ["10.0.0.13"],
	"fourteen_0": ["192.168.102.14"],
	"nineteen_0": ["10.0.0.19"],
	"sixteen_0": ["10.0.0.16"],
	"seventeen_0": ["50.0.0.17"],
	"eighteen_0": ["80.80.100.124"],
	"twenty_0": ["10.0.0.20"],

	"type_aws_instance":                  ["10.0.0.1", "10.0.1.1", "50.0.0.1"],
	"type_digitalocean_droplet":          ["192.168.0.3"],
	"type_cloudstack_instance":           ["10.2.1.5"],
	"type_vsphere_virtual_machine":       ["10.20.30.40", "10.20.30.50"],
	"type_openstack_compute_instance_v2": ["10.120.0.226"],
	"type_opentelekomcloud_compute_instance_v2": ["10.0.0.19"],
	"type_profitbricks_server":           ["10.0.0.16"],
	"type_hcloud_server":                 ["10.0.0.20"],
	"type_softlayer_virtual_guest":       ["10.0.0.7"],
	"type_exoscale_compute":              ["10.0.0.9"],
	"type_google_compute_instance":       ["10.0.0.8"],
	"type_triton_machine":                ["10.0.0.10"],
	"type_scaleway_server":               ["10.0.0.11"],
	"type_packet_device":                 ["10.0.0.13"],
	"type_libvirt_domain":                ["192.168.102.14"],
	"type_aws_spot_instance_request":			["50.0.0.17"],
	"type_linode_instance":               ["80.80.100.124"],

	"role_nine": ["10.0.0.9"],
	"role_rrrrrrrr": ["10.20.30.40"],
	"role_web": ["10.0.0.1"],
	"role_test": ["10.0.0.10"],
	"role_worker": ["50.0.0.17"],
	"webserver": ["192.168.0.3"],
	"staging": ["192.168.0.3"],
	"status_superserver": ["10.120.0.226"],
	"testlabel_hcloud_test": ["10.0.0.20"],
	"database": ["10.0.0.8"],
	"scw_test": ["10.0.0.11"],
	"tfinventory_rocks": ["10.0.0.19"]
}
`

const expectedInventoryOutput = `[all]
10.0.0.1
10.0.0.10
10.0.0.11
10.0.0.13
10.0.0.16
10.0.0.19
10.0.0.20
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
80.80.100.124
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

[dup_0]
10.0.0.1

[eight]
10.0.0.8

[eight_0]
10.0.0.8

[eighteen]
80.80.100.124

[eighteen_0]
80.80.100.124

[eleven]
10.0.0.11

[eleven_0]
10.0.0.11

[five]
10.20.30.40

[five_0]
10.20.30.40

[four]
10.2.1.5

[four_0]
10.2.1.5

[fourteen]
192.168.102.14

[fourteen_0]
192.168.102.14

[nine]
10.0.0.9

[nine_0]
10.0.0.9

[nineteen]
10.0.0.19

[nineteen_0]
10.0.0.19

[one]
10.0.0.1
10.0.1.1

[one_0]
10.0.0.1

[one_1]
10.0.1.1

[role_nine]
10.0.0.9

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

[seven_0]
10.0.0.7

[seventeen]
50.0.0.17

[seventeen_0]
50.0.0.17

[six]
10.120.0.226

[six_0]
10.120.0.226

[sixteen]
10.0.0.16

[sixteen_0]
10.0.0.16

[staging]
192.168.0.3

[status_superserver]
10.120.0.226

[ten]
10.0.0.10

[ten_0]
10.0.0.10

[testTag1]
10.20.30.50

[testlabel_hcloud_test]
10.0.0.20

[tfinventory_rocks]
10.0.0.19

[thirteen]
10.0.0.13

[thirteen_0]
10.0.0.13

[three]
192.168.0.3

[three_0]
192.168.0.3

[twelve]
10.20.30.50

[twelve_0]
10.20.30.50

[twenty]
10.0.0.20

[twenty_0]
10.0.0.20

[two]
50.0.0.1

[two_0]
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

[type_hcloud_server]
10.0.0.20

[type_libvirt_domain]
192.168.102.14

[type_linode_instance]
80.80.100.124

[type_openstack_compute_instance_v2]
10.120.0.226

[type_opentelekomcloud_compute_instance_v2]
10.0.0.19

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
	var s stateAnyTerraformVersion
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.NoError(t, err)

	assert.Equal(t, TerraformVersionPre0dot12, s.TerraformVersion)

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
	var s stateAnyTerraformVersion
	r := strings.NewReader(exampleStateFileEnvHostname)
	err := s.read(r)
	assert.NoError(t, err)

	assert.Equal(t, TerraformVersionPre0dot12, s.TerraformVersion)

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
	var s stateAnyTerraformVersion
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.NoError(t, err)

	assert.Equal(t, TerraformVersionPre0dot12, s.TerraformVersion)

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
	var s stateAnyTerraformVersion
	r := strings.NewReader(exampleStateFile)
	err := s.read(r)
	assert.NoError(t, err)

	assert.Equal(t, TerraformVersionPre0dot12, s.TerraformVersion)

	// Run the command, capture the output
	var stdout, stderr bytes.Buffer
	exitCode := cmdInventory(&stdout, &stderr, &s)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "", stderr.String())

	assert.Equal(t, expectedInventoryOutput, stdout.String())
}

//
// Terraform 0.12 BEGIN
//

const exampleStateFileTerraform0dot12 = `
{
	"format_version": "0.1",
	"terraform_version": "0.12.1",
	"values": {
		"outputs": {
			"my_endpoint": {
				"sensitive": false,
				"value": "a.b.c.d.example.com"
			},
			"my_password": {
				"sensitive": true,
				"value": "1234"
			},
			"map": {
				"sensitive": false,
				"value": {
					"first": "a",
					"second": "b"
				}
			}
		},
		"root_module": {
			"resources": [
				{
					"address": "aws_instance.one",
					"type": "aws_instance",
					"name": "one",
					"provider_name": "aws",
					"schema_version": 1,
					"values": {
						"ami": "ami-00000000000000000",
						"id": "i-11111111111111111",
						"private_ip": "10.0.0.1",
						"public_ip": "35.159.25.34",
						"tags": {
							"Name": "one-aws-instance"
						},
						"volume_tags": {
							"Ignored": "stuff"
						}
					}
				},
				{
					"address": "vsphere_tag.bar",
					"mode": "managed",
					"type": "vsphere_tag",
					"name": "bar",
					"provider_name": "vsphere",
					"schema_version": 0,
					"values": {
						"category_id": "urn:vmomi:InventoryServiceCategory:dc032379-bc2c-4fe5-bd8a-77040e3f4bc8:GLOBAL",
						"description": "",
						"id": "urn:vmomi:InventoryServiceTag:c70f4a73-f744-458a-b2ef-595e3c7c7c28:GLOBAL",
						"name": "bar"
					},
					"depends_on": [
						"vsphere_tag_category.foo"
					]
				},
				{
					"address": "vsphere_tag_category.foo",
					"mode": "managed",
					"type": "vsphere_tag_category",
					"name": "foo",
					"provider_name": "vsphere",
					"schema_version": 0,
					"values": {
						"associable_types": [
							"VirtualMachine"
						],
						"cardinality": "SINGLE",
						"description": "",
						"id": "urn:vmomi:InventoryServiceCategory:dc032379-bc2c-4fe5-bd8a-77040e3f4bc8:GLOBAL",
						"name": "foo"
					}
				},
				{
					"address": "vsphere_virtual_machine.vm",
					"mode": "managed",
					"type": "vsphere_virtual_machine",
					"name": "vm",
					"provider_name": "vsphere",
					"schema_version": 3,
					"values": {
						"alternate_guest_name": "",
						"annotation": "",
						"boot_delay": 0,
						"boot_retry_delay": 10000,
						"boot_retry_enabled": false,
						"cdrom": [],
						"change_version": "2019-08-24T17:27:59.706242Z",
						"clone": [],
						"cpu_hot_add_enabled": false,
						"cpu_hot_remove_enabled": false,
						"cpu_limit": -1,
						"cpu_performance_counters_enabled": false,
						"cpu_reservation": 0,
						"cpu_share_count": 4000,
						"cpu_share_level": "normal",
						"custom_attributes": {},
						"datastore_cluster_id": null,
						"datastore_id": "datastore-1",
						"default_ip_address": "12.34.56.78",
						"disk": [
							{
								"attach": false,
								"datastore_id": "datastore-1",
								"device_address": "scsi:0:0",
								"disk_mode": "persistent",
								"disk_sharing": "sharingNone",
								"eagerly_scrub": false,
								"io_limit": -1,
								"io_reservation": 0,
								"io_share_count": 1000,
								"io_share_level": "normal",
								"keep_on_remove": true,
								"key": 2000,
								"label": "disk0",
								"name": "",
								"path": "foo/bar.vmdk",
								"size": 4,
								"thin_provisioned": true,
								"unit_number": 0,
								"uuid": "6000C292-6cff-cc74-87e5-37ce78a22b57",
								"write_through": false
							}
						],
						"efi_secure_boot_enabled": false,
						"enable_disk_uuid": false,
						"enable_logging": false,
						"ept_rvi_mode": "automatic",
						"extra_config": {},
						"firmware": "bios",
						"folder": "",
						"force_power_off": true,
						"guest_id": "debian8_64Guest",
						"guest_ip_addresses": [
							"12.34.56.78"
						],
						"host_system_id": "host-764",
						"hv_mode": "hvAuto",
						"id": "42361f05-2e60-752c-5999-6c592f0a3904",
						"ignored_guest_ips": null,
						"imported": false,
						"latency_sensitivity": "normal",
						"memory": 2048,
						"memory_hot_add_enabled": false,
						"memory_limit": -1,
						"memory_reservation": 0,
						"memory_share_count": 20480,
						"memory_share_level": "normal",
						"migrate_wait_timeout": 30,
						"moid": "vm-827",
						"name": "vm",
						"nested_hv_enabled": false,
						"network_interface": [
							{
								"adapter_type": "vmxnet3",
								"bandwidth_limit": -1,
								"bandwidth_reservation": 0,
								"bandwidth_share_count": 100,
								"bandwidth_share_level": "high",
								"device_address": "pci:0:7",
								"key": 4000,
								"mac_address": "00:50:56:b3:af:02",
								"network_id": "dvportgroup-837",
								"use_static_mac": false
							}
						],
						"num_cores_per_socket": 1,
						"num_cpus": 4,
						"reboot_required": false,
						"resource_pool_id": "resgroup-768",
						"run_tools_scripts_after_power_on": true,
						"run_tools_scripts_after_resume": true,
						"run_tools_scripts_before_guest_reboot": false,
						"run_tools_scripts_before_guest_shutdown": true,
						"run_tools_scripts_before_guest_standby": true,
						"scsi_bus_sharing": "noSharing",
						"scsi_controller_count": 1,
						"scsi_type": "pvscsi",
						"shutdown_wait_timeout": 3,
						"swap_placement_policy": "inherit",
						"sync_time_with_host": false,
						"tags": [
							"urn:vmomi:InventoryServiceTag:c70f4a73-f744-458a-b2ef-595e3c7c7c28:GLOBAL"
						],
						"uuid": "42361f05-2e60-752c-5999-6c592f0a3904",
						"vapp": [],
						"vapp_transport": [],
						"vmware_tools_status": "guestToolsRunning",
						"vmx_path": "foo/bar.vmx",
						"wait_for_guest_ip_timeout": 0,
						"wait_for_guest_net_routable": true,
						"wait_for_guest_net_timeout": 5
					},
					"depends_on": [
						"vsphere_tag.bar"
					]
				}
			],
			"child_modules": [
				{
					"resources": [
						{
							"address": "aws_instance.host",
							"type": "aws_instance",
							"name": "host",
							"values": {
								"ami": "ami-00000000000000001",
								"id": "i-22222222222222222",
								"private_ip": "10.0.0.2",
								"public_ip": "",
								"tags": {
									"Name": "two-aws-instance"
								}
							}
						}
					],
					"address": "module.my-module-two"
				},
				{
					"resources": [
						{
							"address": "aws_instance.host",
							"type": "aws_instance",
							"name": "host",
							"index": 0,
							"values": {
								"ami": "ami-00000000000000001",
								"id": "i-33333333333333333",
								"private_ip": "10.0.0.3",
								"public_ip": "",
								"tags": {
									"Name": "three-aws-instance"
								}
							}
						},
						{
							"address": "aws_instance.host",
							"type": "aws_instance",
							"name": "host",
							"index": 1,
							"values": {
								"ami": "ami-00000000000000001",
								"id": "i-11133333333333333",
								"private_ip": "10.0.1.3",
								"public_ip": "",
								"tags": {
									"Name": "three-aws-instance"
								}
							}
						}
					],
					"address": "module.my-module-three"
				}
			]
		}
	}
}
`

const expectedListOutputTerraform0dot12 = `
{
	"all": {
		"hosts": [
			"10.0.0.2",
			"10.0.0.3",
			"10.0.1.3",
			"35.159.25.34",
			"12.34.56.78"
		],
		"vars": {
			"my_endpoint": "a.b.c.d.example.com",
			"my_password": "1234",
			"map": {"first": "a", "second": "b"}
		}
	},
	"one_0": ["35.159.25.34"],
	"one": ["35.159.25.34"],
	"module_my-module-two_host_0": ["10.0.0.2"],
	"module_my-module-two_host": ["10.0.0.2"],
	"module_my-module-three_host_0": ["10.0.0.3"],
	"module_my-module-three_host_1": ["10.0.1.3"],
	"module_my-module-three_host": ["10.0.0.3", "10.0.1.3"],

	"type_aws_instance": ["10.0.0.2", "10.0.0.3", "10.0.1.3", "35.159.25.34"],

	"name_one-aws-instance": ["35.159.25.34"],
	"name_two-aws-instance": ["10.0.0.2"],
	"name_three-aws-instance": ["10.0.0.3", "10.0.1.3"],

	"foo_bar": ["12.34.56.78"],
	"type_vsphere_virtual_machine": ["12.34.56.78"],
	"vm_0": ["12.34.56.78"],
	"vm": ["12.34.56.78"]
}
`

const expectedInventoryOutputTerraform0dot12 = `[all]
10.0.0.2
10.0.0.3
10.0.1.3
35.159.25.34
12.34.56.78

[all:vars]
map={"first":"a","second":"b"}
my_endpoint="a.b.c.d.example.com"
my_password="1234"

[foo_bar]
12.34.56.78

[module_my-module-three_host]
10.0.0.3
10.0.1.3

[module_my-module-three_host_0]
10.0.0.3

[module_my-module-three_host_1]
10.0.1.3

[module_my-module-two_host]
10.0.0.2

[module_my-module-two_host_0]
10.0.0.2

[name_one-aws-instance]
35.159.25.34

[name_three-aws-instance]
10.0.0.3
10.0.1.3

[name_two-aws-instance]
10.0.0.2

[one]
35.159.25.34

[one_0]
35.159.25.34

[type_aws_instance]
10.0.0.2
10.0.0.3
10.0.1.3
35.159.25.34

[type_vsphere_virtual_machine]
12.34.56.78

[vm]
12.34.56.78

[vm_0]
12.34.56.78

`

const expectedHostOneOutputTerraform0dot12 = `
{
	"ami": "ami-00000000000000000",
	"ansible_host": "35.159.25.34",
	"id":"i-11111111111111111",
	"private_ip":"10.0.0.1",
	"public_ip": "35.159.25.34",
	"tags.#": "1",
	"tags.Name": "one-aws-instance",
	"volume_tags.#":"1",
	"volume_tags.Ignored":"stuff"
}
`

func TestListCommandTerraform0dot12(t *testing.T) {
	var s stateAnyTerraformVersion
	r := strings.NewReader(exampleStateFileTerraform0dot12)
	err := s.read(r)
	assert.NoError(t, err)

	assert.Equal(t, TerraformVersion0dot12, s.TerraformVersion)

	// Decode expectation as JSON
	var exp interface{}
	err = json.Unmarshal([]byte(expectedListOutputTerraform0dot12), &exp)
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

func TestHostCommandTerraform0dot12(t *testing.T) {
	var s stateAnyTerraformVersion
	r := strings.NewReader(exampleStateFileTerraform0dot12)
	err := s.read(r)
	assert.NoError(t, err)

	assert.Equal(t, TerraformVersion0dot12, s.TerraformVersion)

	// Decode expectation as JSON
	var exp interface{}
	err = json.Unmarshal([]byte(expectedHostOneOutputTerraform0dot12), &exp)
	assert.NoError(t, err)

	// Run the command, capture the output
	var stdout, stderr bytes.Buffer
	exitCode := cmdHost(&stdout, &stderr, &s, "35.159.25.34")
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "", stderr.String())

	// Decode the output to compare
	var act interface{}
	err = json.Unmarshal([]byte(stdout.String()), &act)
	assert.NoError(t, err)

	assert.Equal(t, exp, act)
}

func TestInventoryCommandTerraform0dot12(t *testing.T) {
	var s stateAnyTerraformVersion
	r := strings.NewReader(exampleStateFileTerraform0dot12)
	err := s.read(r)
	assert.NoError(t, err)

	assert.Equal(t, TerraformVersion0dot12, s.TerraformVersion)

	// Run the command, capture the output
	var stdout, stderr bytes.Buffer
	exitCode := cmdInventory(&stdout, &stderr, &s)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "", stderr.String())

	assert.Equal(t, expectedInventoryOutputTerraform0dot12, stdout.String())
}

//
// Terraform 0.12 END
//
