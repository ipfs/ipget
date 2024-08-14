package main

import (
	_ "embed"
	"encoding/json"
	"runtime/debug"
)

var version string

//go:embed version.json
var versionJSON []byte

func init() {
	// Read version from embedded JSON file.
	var verMap map[string]string
	json.Unmarshal(versionJSON, &verMap)
	version = verMap["version"]

	// If running from a module, try to get the build info.
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// Append the revision to the version.
	for i := range bi.Settings {
		if bi.Settings[i].Key == "vcs.revision" {
			version += "-" + bi.Settings[i].Value
			break
		}
	}
}
