package appcore

import "testing"

func TestValidateClipboardPNGSizeRejectsTooLargeInput(t *testing.T) {
	maxBytes := clipboardPNGMaxBytes()
	if err := validateClipboardPNGSize(uintptr(maxBytes + 1)); err == nil {
		t.Fatal("validateClipboardPNGSize accepted oversized clipboard data")
	}
}

func TestValidateClipboardPNGSizeAllowsBoundedInput(t *testing.T) {
	maxBytes := clipboardPNGMaxBytes()
	if err := validateClipboardPNGSize(uintptr(maxBytes)); err != nil {
		t.Fatalf("validateClipboardPNGSize rejected bounded clipboard data: %v", err)
	}
}
