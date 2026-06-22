package appcore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsNewerRelease(t *testing.T) {
	tests := []struct {
		name             string
		currentVersion   string
		currentPublished string
		latestVersion    string
		latestPublished  string
		want             bool
	}{
		{name: "newer version", currentVersion: "v0.1.4.abcdef0", latestVersion: "v0.1.5", want: true},
		{name: "same version same release", currentVersion: "v0.1.4.abcdef0", currentPublished: "2026-06-22T05:00:00Z", latestVersion: "v0.1.4", latestPublished: "2026-06-22T05:03:00Z", want: false},
		{name: "same version newer release time", currentVersion: "v0.1.4.abcdef0", currentPublished: "2026-06-22T05:00:00Z", latestVersion: "v0.1.4", latestPublished: "2026-06-22T05:10:00Z", want: true},
		{name: "older version", currentVersion: "v0.1.4", latestVersion: "v0.1.3", want: false},
		{name: "develop version", currentVersion: "develop", latestVersion: "v0.1.5", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNewerRelease(tt.currentVersion, tt.currentPublished, tt.latestVersion, tt.latestPublished); got != tt.want {
				t.Fatalf("IsNewerRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckLatestRelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Fatal("missing user-agent")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name":"v0.1.5","html_url":"https://github.com/hatolife/ClipForVRChat/releases/tag/v0.1.5","published_at":"2026-06-22T05:00:00Z"}`))
	}))
	defer server.Close()

	client := server.Client()
	oldTransport := client.Transport
	if oldTransport == nil {
		oldTransport = http.DefaultTransport
	}
	client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = "http"
		req.URL.Host = server.Listener.Addr().String()
		return oldTransport.RoundTrip(req)
	})

	got, err := CheckLatestRelease(context.Background(), client, "v0.1.4.abcdef0", "2026-06-22T04:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Available || got.LatestVersion != "v0.1.5" || got.LatestReleasePublished != "2026-06-22T05:00:00Z" || got.URL == "" {
		t.Fatalf("update info = %+v", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
