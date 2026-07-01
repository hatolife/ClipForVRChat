//go:build windows && embeddedspout

package appcore

import _ "embed"

const embeddedSpoutAvailable = true

//go:embed embedded_spout_assets/spout-capture.exe
var embeddedSpoutHelperEXE []byte

//go:embed embedded_spout_assets/SpoutLibrary.dll
var embeddedSpoutLibraryDLL []byte
