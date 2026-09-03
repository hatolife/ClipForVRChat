package main

import (
	"context"
	"encoding/binary"
	"math"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hatolife/ClipForVRChat/internal/appcore"
)

func TestOSCIntegrationReceivesAvatarBeaconBuildsBasisAndForwards(t *testing.T) {
	forwardReceiver := listenUDPForTest(t)
	receivePort := freeUDPPortForTest(t)
	configPath := filepath.Join(t.TempDir(), "config.json")
	cfg := appcore.DefaultConfig()
	cfg.Normalize()
	cfg.AutoCapture.OSC.Host = "127.0.0.1"
	cfg.AutoCapture.OSC.ReceivePort = receivePort
	cfg.AutoCapture.OSC.Forward.Enabled = true
	cfg.AutoCapture.OSC.Forward.Mode = appcore.OSCForwardModeAll
	cfg.AutoCapture.OSC.Forward.Targets = []appcore.OSCForwardTarget{{
		Host: "127.0.0.1",
		Port: forwardReceiver.LocalAddr().(*net.UDPAddr).Port,
	}}
	cfg.AutoCapture.PlayerLocal.BasisSource = appcore.PlayerLocalBasisSourceAvatarOSC

	app := NewApp(configPath, appcore.UIState{Mode: appcore.ModeResults, Config: cfg})
	app.ctx = context.Background()
	app.mu.Lock()
	app.restartCameraPoseReceiverLocked(cfg)
	app.mu.Unlock()
	defer stopCameraPoseReceiverForTest(app)
	waitForOSCReceiverStart(t, appcore.DiagnosticLogPath(configPath))

	target := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: receivePort}
	sender, err := net.DialUDP("udp", nil, target)
	if err != nil {
		t.Fatal(err)
	}
	defer sender.Close()

	samples := []struct {
		address string
		value   float32
	}{
		{address: "/avatar/parameters/avatar_beacon/coord/x", value: 0.75},
		{address: "/avatar/parameters/avatar_beacon/coord/y", value: 0.25},
		{address: "/avatar/parameters/avatar_beacon/coord/z", value: 0.5},
		{address: "/avatar/parameters/avatar_beacon/forward/x", value: 1},
		{address: "/avatar/parameters/avatar_beacon/forward/y", value: 0.5},
		{address: "/avatar/parameters/avatar_beacon/forward/z", value: 0.5},
	}
	for _, sample := range samples {
		packet := buildOSCFloatPacketForIntegration(sample.address, sample.value)
		if _, err := sender.Write(packet); err != nil {
			t.Fatal(err)
		}
	}

	forwardedAddresses := make(map[string]bool, len(samples))
	for range samples {
		packet := readUDPForTest(t, forwardReceiver)
		address, _, _, ok := appcore.ParseOSCPacket(packet)
		if !ok {
			t.Fatalf("転送されたOSC packetを解析できません: %v", packet)
		}
		forwardedAddresses[address] = true
	}
	for _, sample := range samples {
		if !forwardedAddresses[sample.address] {
			t.Fatalf("OSC packetが転送されていません: %s", sample.address)
		}
	}

	snapshot := waitForAvatarOSCBasisReady(t, app)
	if snapshot.RawSampleCount != len(samples) {
		t.Fatalf("RawSampleCount = %d, want %d", snapshot.RawSampleCount, len(samples))
	}
	if math.Abs(snapshot.Pose.Position.X-500) > 0.000001 ||
		math.Abs(snapshot.Pose.Position.Y+500) > 0.000001 ||
		math.Abs(snapshot.Pose.Position.Z) > 0.000001 {
		t.Fatalf("position = %+v, want (500, -500, 0)", snapshot.Pose.Position)
	}
	if math.Abs(snapshot.Pose.Rotation.Y-90) > 0.000001 {
		t.Fatalf("yaw = %v, want 90", snapshot.Pose.Rotation.Y)
	}

	entries := app.GetOSCLogEntries()
	var received int
	var forwarded int
	for _, entry := range entries {
		switch entry.Direction {
		case "receive":
			received++
		case "forward":
			if entry.Status == "ok" {
				forwarded++
			}
		}
	}
	if received != len(samples) || forwarded != len(samples) {
		t.Fatalf("OSC log entries: receive=%d forward=%d, want %d each; entries=%+v", received, forwarded, len(samples), entries)
	}
}

func waitForOSCReceiverStart(t *testing.T, logPath string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(logPath)
		if err == nil && strings.Contains(string(data), "auto-capture osc receiver start") {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("OSC receiverの起動を確認できません: %s", logPath)
}

func waitForAvatarOSCBasisReady(t *testing.T, app *App) PlayerLocalBasisSnapshot {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	var snapshot PlayerLocalBasisSnapshot
	for time.Now().Before(deadline) {
		snapshot = app.GetAvatarOSCBasisStatus()
		if snapshot.Status == "ready" && snapshot.Fresh {
			return snapshot
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("AvatarBeacon basisがreadyになりません: %+v", snapshot)
	return PlayerLocalBasisSnapshot{}
}

func buildOSCFloatPacketForIntegration(address string, value float32) []byte {
	packet := appendOSCStringForIntegration(nil, address)
	packet = appendOSCStringForIntegration(packet, ",f")
	var encoded [4]byte
	binary.BigEndian.PutUint32(encoded[:], math.Float32bits(value))
	return append(packet, encoded[:]...)
}

func appendOSCStringForIntegration(packet []byte, value string) []byte {
	packet = append(packet, value...)
	packet = append(packet, 0)
	for len(packet)%4 != 0 {
		packet = append(packet, 0)
	}
	return packet
}
