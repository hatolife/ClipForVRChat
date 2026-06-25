//go:build !windows

package main

import "fmt"

func revealFileInExplorer(path string) error {
	return fmt.Errorf("Explorerでの表示はWindowsでのみ利用できます: %s", path)
}
