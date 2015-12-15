package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// keyNames contains the names of the keys to check for in each resource in the
// state file. This allows us to support multiple types of resource without too
// much fuss.
var keyNames []string
var nameParser *regexp.Regexp

func init() {
	keyNames = []string{
		"ipv4_address", // DO
		"public_ip",    // AWS
		"private_ip",   // AWS
		"ipaddress",    // CS
		"ip_address",   // VMware
		"access_ip_v4", // OpenStack
		"floating_ip",  // OpenStack
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

// Groups returns the list of Ansible groups which this resource should be
// included in.
func (r Resource) Groups() []string {
	return []string{
		r.BaseName(),
		r.NameWithCounter(),
	}
}

// Attributes returns a map containing everything we know about this resource.
func (r Resource) Attributes() map[string]string {
	return r.State.Primary.Attributes
}

func (r Resource) BaseName() string {
	return r.baseName
}

// NameWithCounter returns the resource name with its counter. For resources
// created without a 'count=' attribute, this will always be zero.
func (r Resource) NameWithCounter() string {
	return fmt.Sprintf("%s.%d", r.baseName, r.counter)
}

// Address returns the IP address of this resource.
func (r Resource) Address() string {
	for _, key := range keyNames {
		if ip := r.State.Primary.Attributes[key]; ip != "" {
			return ip
		}
	}

	return ""
}
