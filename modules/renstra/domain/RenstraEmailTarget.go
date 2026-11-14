package domain

type RenstraEmailTarget struct {
	UUID         string  `gorm:"column:uuid"`
	Tahun        string  `gorm:"column:tahun"`
	KodeAkses    *string `gorm:"column:kodeAkses"`
	TargetAudit  string  `gorm:"column:target_audit"`
	AuditeeName  *string `gorm:"column:auditee_name"`
	AuditeeEmail *string `gorm:"column:auditee_email"`
	Auditor1Name *string `gorm:"column:auditor1_name"`
	Auditor1Mail *string `gorm:"column:auditor1_email"`
	Auditor2Name *string `gorm:"column:auditor2_name"`
	Auditor2Mail *string `gorm:"column:auditor2_email"`
}
