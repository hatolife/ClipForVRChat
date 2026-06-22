package main

import (
	"runtime/debug"
	"strings"
)

var (
	version  = "develop"
	revision = "unknown"
)

const githubURL = "https://github.com/hatolife/ClipForVRChat"

func init() {
	if strings.TrimSpace(revision) == "" || revision == "unknown" {
		revision = buildInfoRevision()
	}
}

func appVersion() string {
	if strings.TrimSpace(revision) == "" || revision == "unknown" {
		return version
	}
	return version + "." + revision
}

func buildInfoRevision() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		rev := ""
		modified := false
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				rev = setting.Value
				if len(rev) > 7 {
					rev = rev[:7]
				}
			case "vcs.modified":
				modified = setting.Value == "true"
			}
		}
		if rev != "" {
			if modified {
				return "develop+" + rev
			}
			return rev
		}
	}
	return "unknown"
}
