package domain

type SearchFilter struct {
    Field    string `json:"field"`
    Operator string `json:"operator"`
    Value    *string `json:"value"`
}
