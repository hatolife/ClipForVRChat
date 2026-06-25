package main

import (
	"runtime/debug"
	"strings"
)

var (
	version      = "v0.1.7"
	revision     = "unknown"
	releaseTime  = ""
	buildChannel = "develop"
)

const githubURL = "https://github.com/hatolife/ClipForVRChat"

func init() {
	if strings.TrimSpace(revision) == "" || revision == "unknown" {
		revision = buildInfoRevision()
	}
}

func appVersion() string {
	if strings.TrimSpace(revision) == "" || revision == "unknown" {
		if strings.TrimSpace(buildChannel) == "develop" {
			return version + ".unknown.develop"
		}
		return version
	}
	appVersion := version + "." + revision
	if strings.TrimSpace(buildChannel) == "develop" {
		appVersion += ".develop"
	}
	return appVersion
}

func appReleaseTime() string {
	return strings.TrimSpace(releaseTime)
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
