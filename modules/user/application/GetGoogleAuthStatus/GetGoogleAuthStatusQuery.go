package application

type GetGoogleAuthStatusQuery struct{}

type GetGoogleAuthStatusResult struct {
	HasSMTP           bool   `json:"has_smtp"`
	SMTPEmail         string `json:"smtp_email"`
	SMTPHost          string `json:"smtp_host"`
	PreferredProvider string `json:"preferred_provider"`
	IsReady           bool   `json:"is_ready"`
}
