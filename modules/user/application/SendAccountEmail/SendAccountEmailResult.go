package application

type SendAccountEmailResult struct {
	Message    string   `json:"message"`
	TotalQueued int      `json:"totalQueued"`
	UserUUIDs  []string `json:"userUuids"`
	OutboxID   string   `json:"outboxId"`
}
