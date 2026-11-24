package RestoreMataProgramV2

type RestoreMataProgramV2Command struct {
	UUID string
}

type RestoreMataProgramV2Result struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	RowsAffected int64  `json:"rows_affected"`
}
