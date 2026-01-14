package application

type UpdateRenstraNilaiCommand struct {
    Uuid  			string
	UuidRenstra     string
	Tahun     		string
	Mode 			string
	Granted 		string
	Capaian         *string
	Catatan         *string
	LinkBukti      	*string

	CapaianAuditor	*string
	CatatanAuditor 	*string
}
