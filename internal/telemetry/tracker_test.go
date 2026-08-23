package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTracker_RecordAndExport(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "telemetry.json")

	tracker := NewTracker()
	tracker.RecordSource(ChannelStat{
		Name:          "test_channel",
		Type:          "channel",
		StatusCode:    200,
		Duration:      100 * time.Millisecond,
		MessagesCount: 15,
		ConfigsYield:  5,
	})
	tracker.RecordSource(ChannelStat{
		Name:       "bad_sub",
		Type:       "subscription",
		StatusCode: 404,
		Duration:   50 * time.Millisecond,
		Error:      "received status code 404",
	})

	tracker.RecordDropped("test_channel", "customproto://12345", "unsupported or unrecognized proxy scheme")
	tracker.RecordDropped("test_channel", "vmess://short", "vmess payload too short or malformed")

	uniqueCounts := map[string]int{
		"vless": 3,
		"vmess": 2,
	}

	if err := tracker.ExportReport(reportPath, uniqueCounts); err != nil {
		t.Fatalf("Failed to export telemetry report: %v", err)
	}

	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read report file: %v", err)
	}

	var report SummaryReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("Report JSON is invalid: %v", err)
	}

	if report.TotalChannels != 2 {
		t.Errorf("Expected 2 total channels, got %d", report.TotalChannels)
	}
	if report.ActiveChannels != 1 {
		t.Errorf("Expected 1 active channel, got %d", report.ActiveChannels)
	}
	if report.FailedChannels != 1 {
		t.Errorf("Expected 1 failed channel, got %d", report.FailedChannels)
	}
	if report.TotalUnique != 5 {
		t.Errorf("Expected 5 total unique, got %d", report.TotalUnique)
	}
	if report.DroppedCount != 2 {
		t.Errorf("Expected 2 dropped items, got %d", report.DroppedCount)
	}
}
