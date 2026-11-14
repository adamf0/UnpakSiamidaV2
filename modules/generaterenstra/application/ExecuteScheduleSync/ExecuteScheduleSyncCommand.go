package ExecuteScheduleSync

type ExecuteScheduleSyncCommand struct {
	UUIDs        []string `json:"uuids"`
	Confirmation string   `json:"confirmation"`
}

type ExecuteScheduleSyncResult struct {
	Message          string   `json:"message"`
	ProcessedCount   int      `json:"processedCount"`
	ProcessedUUIDs   []string `json:"processedUuids"`
	ProcessedTargets []string `json:"processedTargets"`
}
