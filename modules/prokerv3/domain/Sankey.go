package domain

type SankeyNode struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color,omitempty"`
}

type SankeyLink struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Value  float64 `json:"value"`
}

type SankeyDataResult struct {
	Tahun        string                   `json:"tahun"`
	Monitoring   []map[string]interface{} `json:"monitoring"`
	MasterIkus   []map[string]interface{} `json:"master_ikus"`
	FakultasUnit []map[string]interface{} `json:"fakultas_unit"`
}
