package internal

import (
	"fmt"
	"strings"
)

var (
	version = "unknown"
	build   = ""
	hash    = ""
	// Release contains the current release
	Release MigrateRelease
)

type MigrateRelease struct {
	vendor  string
	version string
	build   string
	hash    string
}

func init() {
	Release = MigrateRelease{
		vendor:  "Thomas Meitz",
		version: version,
		build:   build,
		hash:    hash,
	}
}

func (r *MigrateRelease) shortenHash() string {
	if r.hash == "" {
		return "unknown"
	}
	if len(r.hash) >= 8 {
		return r.hash[:7]
	}
	return r.hash
}

func (r *MigrateRelease) String() string {
	splittedVersion := strings.Split(r.version, "-")
	if len(splittedVersion) > 1 {
		return fmt.Sprintf("Copyright by %v\n%v %v", r.vendor, r.version, r.build)
	}
	return fmt.Sprintf("Copyright by %v\n%v-%v %v", r.vendor, r.version, r.shortenHash(), r.build)
}
