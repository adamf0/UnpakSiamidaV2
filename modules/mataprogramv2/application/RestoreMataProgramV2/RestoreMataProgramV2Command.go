package RestoreMataProgramV2

type RestoreMataProgramV2Command struct {
	UUID  string
	Scope string // "single" (default) or "group"
}

type RestoreMataProgramV2Result struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	RowsAffected int64  `json:"rows_affected"`
}
