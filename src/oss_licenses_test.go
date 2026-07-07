package main

import (
	"strings"
	"testing"
)

func TestOSSLicensesAreCompleteAndUnique(t *testing.T) {
	licenses := ossLicenses()
	if len(licenses) == 0 {
		t.Fatal("ossLicenses returned no entries")
	}

	seen := map[string]bool{}
	for _, license := range licenses {
		if strings.TrimSpace(license.Name) == "" {
			t.Fatalf("license has empty name: %+v", license)
		}
		if seen[license.Name] {
			t.Fatalf("duplicate OSS license entry: %s", license.Name)
		}
		seen[license.Name] = true
		if strings.TrimSpace(license.License) == "" {
			t.Fatalf("%s has empty license", license.Name)
		}
		if strings.TrimSpace(license.Copyright) == "" {
			t.Fatalf("%s has empty copyright", license.Name)
		}
		if strings.TrimSpace(license.URL) == "" {
			t.Fatalf("%s has empty URL", license.Name)
		}
		if strings.TrimSpace(license.Text) == "" {
			t.Fatalf("%s has empty license text", license.Name)
		}
		if !strings.Contains(license.Text, "SOFTWARE IS PROVIDED") {
			t.Fatalf("%s license text does not include disclaimer", license.Name)
		}
	}

	required := []string{
		"ClipForVRChat",
		"@vitejs/plugin-vue",
		"Vite",
		"Vue.js",
		"Wails",
		"cloudflare/circl",
		"flock",
		"go-ansi-parser",
		"go-arg",
		"go-scalar",
		"go-webview2",
		"golang.design/x/clipboard",
		"golang.org/x/crypto",
		"golang.org/x/image",
		"golang.org/x/sys",
		"golang.org/x/text",
		"golang.org/x/xerrors",
		"gozxing",
		"imaging",
		"pkg/errors",
		"ProtonMail/go-crypto",
		"rivo/uniseg",
		"slicer",
		"Spout2",
		"u",
	}
	for _, name := range required {
		if !seen[name] {
			t.Fatalf("required OSS license entry %q is missing", name)
		}
	}

	if seen["go-qrcode"] {
		t.Fatal("go-qrcode is test-only and should not be shown in app OSS licenses")
	}
}

func TestOSSLicenseAuditCorrections(t *testing.T) {
	licenses := map[string]OSSLicense{}
	for _, license := range ossLicenses() {
		licenses[license.Name] = license
	}

	checks := map[string]string{
		"go-arg":               "BSD-2-Clause",
		"flock":                "BSD-2-Clause",
		"ProtonMail/go-crypto": "BSD-3-Clause",
	}
	for name, want := range checks {
		if got := licenses[name].License; got != want {
			t.Fatalf("%s license = %q, want %q", name, got, want)
		}
	}
}
