package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Tracker collects metrics, stats, and dropped candidate logs across the scraping pipeline
type Tracker struct {
	mu           sync.Mutex
	startTime    time.Time
	channelStats []ChannelStat
	droppedItems []DroppedItem
	maxDropped   int
}

// Global is the singleton tracker instance
var Global = NewTracker()

func NewTracker() *Tracker {
	return &Tracker{
		startTime:    time.Now(),
		channelStats: make([]ChannelStat, 0),
		droppedItems: make([]DroppedItem, 0),
		maxDropped:   200,
	}
}

func (t *Tracker) RecordSource(stat ChannelStat) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.channelStats = append(t.channelStats, stat)
}

func (t *Tracker) RecordDropped(source, candidate, reason string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.droppedItems) < t.maxDropped {
		t.droppedItems = append(t.droppedItems, DroppedItem{
			Source:    source,
			Candidate: candidate,
			Reason:    reason,
			Timestamp: time.Now(),
		})
	}
}

func (t *Tracker) ExportReport(path string, uniqueCounts map[string]int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	active := 0
	failed := 0
	totalYield := 0
	for _, cs := range t.channelStats {
		if cs.Error != "" || cs.StatusCode != 200 {
			failed++
		} else {
			active++
		}
		totalYield += cs.ConfigsYield
	}

	totalUnique := 0
	for _, count := range uniqueCounts {
		totalUnique += count
	}

	report := SummaryReport{
		Timestamp:         time.Now(),
		DurationSeconds:   time.Since(t.startTime).Seconds(),
		TotalChannels:     len(t.channelStats),
		ActiveChannels:    active,
		FailedChannels:    failed,
		TotalExtracted:    totalYield,
		TotalUnique:       totalUnique,
		ProtocolBreakdown: uniqueCounts,
		DroppedCount:      len(t.droppedItems),
		DroppedItems:      t.droppedItems,
		ChannelStats:      t.channelStats,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (t *Tracker) PrintConsoleSummary(uniqueCounts map[string]int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fmt.Println("\n================= Observability & Telemetry Summary =================")
	fmt.Printf("Total Sources Processed : %d\n", len(t.channelStats))
	fmt.Printf("Total Pipeline Duration : %.2fs\n", time.Since(t.startTime).Seconds())
	if len(t.droppedItems) > 0 {
		fmt.Printf("Logged Dropped Items    : %d (saved to sub/telemetry.json)\n", len(t.droppedItems))
	}
	fmt.Println("---------------------------------------------------------------------")
	for proto, count := range uniqueCounts {
		fmt.Printf("  %-10s : %d unique\n", proto, count)
	}
	fmt.Println("=====================================================================")
}
