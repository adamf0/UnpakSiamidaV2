package DeleteMataProgramV2

type DeleteMataProgramV2Command struct {
	UUID  string
	Scope string // "single" or "group" (default)
}

type DeleteMataProgramV2Result struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	RowsAffected int64  `json:"rows_affected"`
}
