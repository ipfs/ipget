package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"
)

//go:embed version.json
var versionJSON []byte

var version = buildVersion()

func buildVersion() string {
	// Read version from embedded JSON file.
	var v struct {
		Version string `json:"version"`
	}
	json.Unmarshal(versionJSON, &v)
	release := v.Version

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return release + " dev-build"
	}

	var dirty bool
	var day, revision string

	// Append the revision to the version.
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			revision = kv.Value[:7]
		case "vcs.time":
			t, _ := time.Parse(time.RFC3339, kv.Value)
			day = t.UTC().Format("2006-01-02")
		case "vcs.modified":
			dirty = kv.Value == "true"
		}
	}
	if dirty {
		revision += "-dirty"
	}
	if revision != "" {
		return fmt.Sprintf("%s %s-%s", release, day, revision)
	}
	return release + " dev-build"
}
