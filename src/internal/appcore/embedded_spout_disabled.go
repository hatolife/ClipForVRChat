//go:build !windows || !embeddedspout

package appcore

const embeddedSpoutAvailable = false

var (
	embeddedSpoutHelperEXE  []byte
	embeddedSpoutLibraryDLL []byte
)
