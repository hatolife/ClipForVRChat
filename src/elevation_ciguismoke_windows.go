//go:build windows && ciguismoke

package main

// rejectElevatedProcess is disabled only in the CI-only GUI smoke binary.
// Release builds never set the ciguismoke tag and keep the normal safety check.
func rejectElevatedProcess() error {
	return nil
}
