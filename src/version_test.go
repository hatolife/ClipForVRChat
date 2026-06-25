package main

import "testing"

func TestAppVersionIncludesRevision(t *testing.T) {
	oldVersion := version
	oldRevision := revision
	oldBuildChannel := buildChannel
	t.Cleanup(func() {
		version = oldVersion
		revision = oldRevision
		buildChannel = oldBuildChannel
	})

	version = "v1.2.3"
	revision = "abcdef0"
	buildChannel = "release"

	if got := appVersion(); got != "v1.2.3.abcdef0" {
		t.Fatalf("appVersion() = %q, want v1.2.3.abcdef0", got)
	}
}

func TestAppVersionIncludesDevelopChannel(t *testing.T) {
	oldVersion := version
	oldRevision := revision
	oldBuildChannel := buildChannel
	t.Cleanup(func() {
		version = oldVersion
		revision = oldRevision
		buildChannel = oldBuildChannel
	})

	version = "v1.2.3"
	revision = "abcdef0"
	buildChannel = "develop"

	if got := appVersion(); got != "v1.2.3.abcdef0.develop" {
		t.Fatalf("appVersion() = %q, want v1.2.3.abcdef0.develop", got)
	}
}
