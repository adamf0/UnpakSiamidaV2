package application

type SendEmailRenstraResult struct {
	Message        string   `json:"message"`
	QueuedCount    int      `json:"queuedCount"`
	QueuedUuids    []string `json:"queuedUuids"`
	OutboxIds      []string `json:"outboxIds"`
	SkippedCount   int      `json:"skippedCount"`
	SkippedTargets []string `json:"skippedTargets"`
}
