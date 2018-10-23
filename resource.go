package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// keyNames contains the names of the keys to check for in each resource in the
// state file. This allows us to support multiple types of resource without too
// much fuss.
var keyNames []string
var nameParser *regexp.Regexp

func init() {
	keyNames = []string{
		"ipv4_address",                                        // DO and SoftLayer
		"public_ip",                                           // AWS
		"public_ipv6",                                         // Scaleway
		"private_ip",                                          // AWS
		"ipaddress",                                           // CS
		"ip_address",                                          // VMware, Docker
		"network_interface.0.ipv4_address",                    // VMware
		"default_ip_address",                                  // provider.vsphere v1.1.1
		"access_ip_v4",                                        // OpenStack
		"floating_ip",                                         // OpenStack
		"network_interface.0.access_config.0.nat_ip",          // GCE
		"network_interface.0.access_config.0.assigned_nat_ip", // GCE
		"network_interface.0.address",                         // GCE
		"ipv4_address_private",                                // SoftLayer
		"networks.0.ip4address",                               // Exoscale
		"primaryip",                                           // Joyent Triton
		"network_interface.0.addresses.0",                     // Libvirt
		"network.0.address",                                   // Packet
		"primary_ip",                                          // Profitbricks
	}

	// type.name.0
	nameParser = regexp.MustCompile(`^(\w+)\.([\w\-]+)(?:\.(\d+))?$`)
}

type Resource struct {

	// The state (as unmarshalled from the statefile) which this resource wraps.
	// Everything which Terraform knows about the resource can be found in here.
	State resourceState

	// The key name of the resource, provided to the constructor. Unfortunately,
	// it seems like the counter index can only be found here.
	keyName string

	// Extracted from keyName
	resourceType string
	baseName     string
	counter      int
}

func NewResource(keyName string, state resourceState) (*Resource, error) {
	m := nameParser.FindStringSubmatch(keyName)

	// This should not happen unless our regex changes.
	// TODO: Warn instead of silently ignore error?
	if len(m) != 4 {
		return nil, fmt.Errorf("couldn't parse keyName: %s", keyName)
	}

	var c int
	var err error
	if m[3] != "" {

		// The third section should be the index, if it's present. Not sure what
		// else we can do other than panic (which seems highly undesirable) if that
		// isn't the case.
		c, err = strconv.Atoi(m[3])
		if err != nil {
			return nil, err
		}
	}

	return &Resource{
		State:        state,
		keyName:      keyName,
		resourceType: m[1],
		baseName:     m[2],
		counter:      c,
	}, nil
}

func (r Resource) IsSupported() bool {
	return r.Address() != ""
}

// Tags returns a map of arbitrary key/value pairs explicitly associated with
// the resource. Different providers have different mechanisms for attaching
// these.
func (r Resource) Tags() map[string]string {
	t := map[string]string{}

	switch r.resourceType {
	case "openstack_compute_instance_v2":
		for k, v := range r.Attributes() {
			parts := strings.SplitN(k, ".", 2)
			// At some point Terraform changed the key for counts of attributes to end with ".%"
			// instead of ".#". Both need to be considered as Terraform still supports state
			// files using the old format.
			if len(parts) == 2 && parts[0] == "metadata" && parts[1] != "#" && parts[1] != "%" {
				kk := strings.ToLower(parts[1])
				vv := strings.ToLower(v)
				t[kk] = vv
			}
		}
	case "aws_instance":
		for k, v := range r.Attributes() {
			parts := strings.SplitN(k, ".", 2)
			// At some point Terraform changed the key for counts of attributes to end with ".%"
			// instead of ".#". Both need to be considered as Terraform still supports state
			// files using the old format.
			if len(parts) == 2 && parts[0] == "tags" && parts[1] != "#" && parts[1] != "%" {
				kk := strings.ToLower(parts[1])
				vv := strings.ToLower(v)
				t[kk] = vv
			}
		}
	case "aws_spot_instance_request":
		for k, v := range r.Attributes() {
			parts := strings.SplitN(k, ".", 2)
			// At some point Terraform changed the key for counts of attributes to end with ".%"
			// instead of ".#". Both need to be considered as Terraform still supports state
			// files using the old format.
			if len(parts) == 2 && parts[0] == "tags" && parts[1] != "#" && parts[1] != "%" {
				kk := strings.ToLower(parts[1])
				vv := strings.ToLower(v)
				t[kk] = vv
			}
		}
	case "vsphere_virtual_machine":
		for k, v := range r.Attributes() {
			parts := strings.SplitN(k, ".", 2)

			if len(parts) == 2 && parts[0] == "custom_configuration_parameters" && parts[1] != "#" && parts[1] != "%" {
				kk := strings.ToLower(parts[1])
				vv := strings.ToLower(v)
				t[kk] = vv
			}
			if len(parts) == 2 && parts[0] == "tags" && parts[1] != "#" && parts[1] != "%" {
				kk := strings.ToLower(parts[1])
				vv := strings.ToLower(v)
				t[kk] = vv
			}
		}
	case "digitalocean_droplet", "google_compute_instance", "scaleway_server":
		for k, v := range r.Attributes() {
			parts := strings.SplitN(k, ".", 2)
			if len(parts) == 2 && parts[0] == "tags" && parts[1] != "#" {
				vv := strings.ToLower(v)
				t[vv] = ""
			}
		}
	case "triton_machine":
		for k, v := range r.Attributes() {
			parts := strings.SplitN(k, ".", 2)
			if len(parts) == 2 && parts[0] == "tags" && parts[1] != "%" {
				kk := strings.ToLower(parts[1])
				vv := strings.ToLower(v)
				t[kk] = vv
			}
		}
	}
	return t
}

// Attributes returns a map containing everything we know about this resource.
func (r Resource) Attributes() map[string]string {
	return r.State.Primary.Attributes
}

// NameWithCounter returns the resource name with its counter. For resources
// created without a 'count=' attribute, this will always be zero.
func (r Resource) NameWithCounter() string {
	return fmt.Sprintf("%s.%d", r.baseName, r.counter)
}

// Hostname returns the hostname of this resource.
func (r Resource) Hostname() string {
	if keyName := os.Getenv("TF_HOSTNAME_KEY_NAME"); keyName != "" {
		if ip := r.State.Primary.Attributes[keyName]; ip != "" {
			return ip
		}
	}

	return r.Address()
}

// Address returns the IP address of this resource.
func (r Resource) Address() string {
	if keyName := os.Getenv("TF_KEY_NAME"); keyName != "" {
		if ip := r.State.Primary.Attributes[keyName]; ip != "" {
			return ip
		}
	} else {
		for _, key := range keyNames {
			if ip := r.State.Primary.Attributes[key]; ip != "" {
				return ip
			}
		}
	}

	return ""
}
