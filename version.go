package main

// Deliberately uninitialized. See below.
var build_version string

// versionInfo returns a string containing the version information of the
// current build. It's empty by default, but can be included as part of the
// build process by setting the main.build_version variable.
func versionInfo() string {
	if build_version != "" {
		return build_version
	} else {
		return "unknown"
	}
}
