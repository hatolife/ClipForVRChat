package main

import (
	"path/filepath"
	"testing"
)

func TestAcquireInstanceLockPreventsSecondLock(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	first, err := acquireInstanceLock(configPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = first.Unlock()
	}()

	second, err := acquireInstanceLock(configPath)
	if err == nil {
		_ = second.Unlock()
		t.Fatal("expected second lock to fail")
	}
}

func TestAcquireInstanceLockAllowsAfterUnlock(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	first, err := acquireInstanceLock(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Unlock(); err != nil {
		t.Fatal(err)
	}

	second, err := acquireInstanceLock(configPath)
	if err != nil {
		t.Fatal(err)
	}
	_ = second.Unlock()
}
