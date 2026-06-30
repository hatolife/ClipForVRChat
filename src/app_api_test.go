package main

import (
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestFrontendAPICallsHaveAppMethods(t *testing.T) {
	data, err := os.ReadFile("frontend/src/main.js")
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`api\??\.([A-Z][A-Za-z0-9_]*)`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	calls := map[string]bool{}
	for _, match := range matches {
		calls[match[1]] = true
	}
	appType := reflect.TypeOf(&App{})
	missing := make([]string, 0)
	for call := range calls {
		if _, ok := appType.MethodByName(call); !ok {
			missing = append(missing, call)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("frontend calls missing App methods: %s", strings.Join(missing, ", "))
	}
}
