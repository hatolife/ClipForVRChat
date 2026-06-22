package main

import "testing"

func TestAppVersionIncludesRevision(t *testing.T) {
	oldVersion := version
	oldRevision := revision
	t.Cleanup(func() {
		version = oldVersion
		revision = oldRevision
	})

	version = "v1.2.3"
	revision = "abcdef0"

	if got := appVersion(); got != "v1.2.3.abcdef0" {
		t.Fatalf("appVersion() = %q, want v1.2.3.abcdef0", got)
	}
}

func TestAppVersionOmitsUnknownRevision(t *testing.T) {
	oldVersion := version
	oldRevision := revision
	t.Cleanup(func() {
		version = oldVersion
		revision = oldRevision
	})

	version = "develop"
	revision = "unknown"

	if got := appVersion(); got != "develop" {
		t.Fatalf("appVersion() = %q, want develop", got)
	}
}
