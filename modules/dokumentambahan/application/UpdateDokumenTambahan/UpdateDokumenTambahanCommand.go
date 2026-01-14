package application

type UpdateDokumenTambahanCommand struct {
    Uuid  			string
	UuidRenstra     string
	Tahun     		string
	Mode 			string
	Granted 		string
	Link      		*string

	CapaianAuditor	*string
	CatatanAuditor 	*string
}
