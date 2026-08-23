package telemetry

import "time"

// DroppedItem records a candidate config or message that failed extraction or validation
type DroppedItem struct {
	Source    string    `json:"source"`
	Candidate string    `json:"candidate"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// ChannelStat records performance and yield for a single scraped source
type ChannelStat struct {
	Name          string        `json:"name"`
	Type          string        `json:"type"` // "channel" or "subscription"
	StatusCode    int           `json:"status_code"`
	Duration      time.Duration `json:"duration_ms"`
	MessagesCount int           `json:"messages_scanned"`
	ConfigsYield  int           `json:"configs_yield"`
	Error         string        `json:"error,omitempty"`
}

// DuplicateEntry records a single instance of a config
type DuplicateEntry struct {
	Source string `json:"source"`
	Raw    string `json:"raw"`
}

// DuplicateGroup holds a retained config alongside all its dropped duplicate copies
type DuplicateGroup struct {
	Protocol    string           `json:"protocol"`
	Fingerprint string           `json:"fingerprint"`
	Retained    DuplicateEntry   `json:"retained"`
	Duplicates  []DuplicateEntry `json:"duplicates"`
}

// DuplicatesReport provides full visibility into every merged config group
type DuplicatesReport struct {
	Timestamp          time.Time        `json:"timestamp"`
	TotalExtracted     int              `json:"total_extracted"`
	TotalUnique        int              `json:"total_unique"`
	TotalDuplicates    int              `json:"total_duplicates_dropped"`
	DuplicationRatePct float64          `json:"duplication_rate_percent"`
	DuplicateGroups    []DuplicateGroup `json:"duplicate_groups"`
}

// SummaryReport holds aggregated execution metrics, dropped candidates, and deduplication stats
type SummaryReport struct {
	Timestamp          time.Time         `json:"timestamp"`
	DurationSeconds    float64           `json:"duration_seconds"`
	TotalChannels      int               `json:"total_channels"`
	ActiveChannels     int               `json:"active_channels"`
	FailedChannels     int               `json:"failed_channels"`
	TotalExtracted     int               `json:"total_extracted"`
	TotalUnique        int               `json:"total_unique"`
	TotalDuplicates    int               `json:"total_duplicates"`
	DuplicationRatePct float64           `json:"duplication_rate_percent"`
	ProtocolBreakdown  map[string]int    `json:"protocol_breakdown"`
	DroppedCount       int               `json:"dropped_candidates_count"`
	DroppedItems       []DroppedItem     `json:"sample_dropped_items,omitempty"`
	ChannelStats       []ChannelStat     `json:"channel_stats"`
}
