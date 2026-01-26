package application

type UpdateKtsCommand struct {
	Uuid                   string
	NomorLaporan           *string
	TanggalLaporan         *string
	UraianKetidaksesuaianP *string //ckeditor
	UraianKetidaksesuaianL *string //ckeditor
	UraianKetidaksesuaianO *string //ckeditor
	UraianKetidaksesuaianR *string //ckeditor
	AkarMasalah            *string //ckeditor
	TindakanKoreksi        *string //ckeditor
	Acc                    string

	StatusAccAuditee *string //step2
	// accAuditee					string
	KeteranganTolak   *string //step2
	TindakanPerbaikan *string //step2

	TanggalPenyelesaian *string //step3

	TinjauanTindakanPerbaikan *string //step4
	TanggalClosing            *string //step4
	// accFinal						string

	TanggalClosingFinal *string //step5
	WmmUpmfUpmps        *string //step5
	// closingBy

	Tahun string
	Step  string
}
