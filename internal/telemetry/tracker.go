package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Tracker collects metrics, stats, dropped candidate logs, and duplicate groups across the scraping pipeline
type Tracker struct {
	mu              sync.Mutex
	startTime       time.Time
	channelStats    []ChannelStat
	droppedItems    []DroppedItem
	maxDropped      int
	duplicateGroups map[string]*DuplicateGroup
	totalDuplicates int
}

// Global is the singleton tracker instance
var Global = NewTracker()

func NewTracker() *Tracker {
	return &Tracker{
		startTime:       time.Now(),
		channelStats:    make([]ChannelStat, 0),
		droppedItems:    make([]DroppedItem, 0),
		maxDropped:      200,
		duplicateGroups: make(map[string]*DuplicateGroup),
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

// RecordDuplicate logs a duplicate config event under its unique fingerprint group
func (t *Tracker) RecordDuplicate(proto, fp, retainedSource, retainedRaw, dupSource, dupRaw string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.totalDuplicates++
	key := proto + ":" + fp
	group, exists := t.duplicateGroups[key]
	if !exists {
		group = &DuplicateGroup{
			Protocol:    proto,
			Fingerprint: fp,
			Retained: DuplicateEntry{
				Source: retainedSource,
				Raw:    retainedRaw,
			},
			Duplicates: make([]DuplicateEntry, 0, 1),
		}
		t.duplicateGroups[key] = group
	}

	group.Duplicates = append(group.Duplicates, DuplicateEntry{
		Source: dupSource,
		Raw:    dupRaw,
	})
}

// ExportReport writes the aggregated telemetry summary to JSON
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

	duplicationRate := 0.0
	if totalYield > 0 {
		duplicationRate = (float64(t.totalDuplicates) / float64(totalYield)) * 100.0
	}

	report := SummaryReport{
		Timestamp:          time.Now(),
		DurationSeconds:    time.Since(t.startTime).Seconds(),
		TotalChannels:      len(t.channelStats),
		ActiveChannels:     active,
		FailedChannels:     failed,
		TotalExtracted:     totalYield,
		TotalUnique:        totalUnique,
		TotalDuplicates:    t.totalDuplicates,
		DuplicationRatePct: duplicationRate,
		ProtocolBreakdown:  uniqueCounts,
		DroppedCount:       len(t.droppedItems),
		DroppedItems:       t.droppedItems,
		ChannelStats:       t.channelStats,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ExportDuplicates writes full details of all dropped duplicates and their groups to JSON
func (t *Tracker) ExportDuplicates(path string, totalExtracted, totalUnique int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	groups := make([]DuplicateGroup, 0, len(t.duplicateGroups))
	for _, g := range t.duplicateGroups {
		groups = append(groups, *g)
	}

	if totalExtracted <= 0 {
		for _, cs := range t.channelStats {
			totalExtracted += cs.ConfigsYield
		}
	}

	duplicationRate := 0.0
	if totalExtracted > 0 {
		duplicationRate = (float64(t.totalDuplicates) / float64(totalExtracted)) * 100.0
	}

	report := DuplicatesReport{
		Timestamp:          time.Now(),
		TotalExtracted:     totalExtracted,
		TotalUnique:        totalUnique,
		TotalDuplicates:    t.totalDuplicates,
		DuplicationRatePct: duplicationRate,
		DuplicateGroups:    groups,
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

	totalYield := 0
	for _, cs := range t.channelStats {
		totalYield += cs.ConfigsYield
	}
	totalUnique := 0
	for _, count := range uniqueCounts {
		totalUnique += count
	}

	fmt.Println("\n================= Observability & Telemetry Summary =================")
	fmt.Printf("Total Sources Processed : %d\n", len(t.channelStats))
	fmt.Printf("Total Pipeline Duration : %.2fs\n", time.Since(t.startTime).Seconds())
	fmt.Printf("Total Configs Extracted : %d\n", totalYield)
	fmt.Printf("Unique Configs Retained : %d\n", totalUnique)
	if totalYield > 0 {
		fmt.Printf("Duplicates Deduplicated: %d (%.1f%% duplication rate)\n",
			t.totalDuplicates, (float64(t.totalDuplicates)/float64(totalYield))*100.0)
	}
	if len(t.droppedItems) > 0 {
		fmt.Printf("Logged Dropped Items    : %d\n", len(t.droppedItems))
	}
	exportEnv := strings.ToLower(os.Getenv("EXPORT_TELEMETRY"))
	if exportEnv == "true" || exportEnv == "1" {
		fmt.Printf("Reports Output          : saved to reports/telemetry.json & duplicates.json\n")
	}
	fmt.Println("---------------------------------------------------------------------")
	for proto, count := range uniqueCounts {
		fmt.Printf("  %-10s : %d unique\n", proto, count)
	}
	fmt.Println("=====================================================================")
}
