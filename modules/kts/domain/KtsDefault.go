package domain

import (
	"github.com/google/uuid"
)

type KtsDefault struct {
	Id          				uint
	UUID            			uuid.UUID  
	RenstraId              		*uint
	RenstraNilai              	*uint       
	DokumenTambahan             *uint       
	Status       				string     

	TemplateDokumen				*uint
	Pertanyaan					*string
	JenisFileId					*uint
	JenisFile					*string

	TemplateRenstra				*uint
	Standar						*string
	StandarId					*uint
	Indikator					*string
	IndikatorId					*uint

	Tahun						*string
	IdTarget					*uint
	Target 						*string

	//auditor (1)
	NomorLaporan       			*string     
	TanggalLaporan       		*string     
	Auditor       				*string
	Auditee       				*string     
	KetidaksesuaianP       		*string
	KetidaksesuaianL	 		*string
	KetidaksesuaianO	 		*string
	KetidaksesuaianR	 		*string
	AkarMasalah       			*string
	TindakanKoreksi       		*string
	AccAuditor       			*uint

	//auditee (2)
	StatusAccAuditee       		*uint
	AccAuditee       			*uint 
	KeteranganTolak       		*string
	TindakanPerbaikan       	*string
	
	//auditor (3)
	TanggalPenyelesaian       	*string

	//auditee (4)
	TinjauanTindakanPerbaikan  	*string
	TanggalClosing       		*string
	AccFinal       				*uint

	//auditor (5)
	TanggalClosingFinal  		*string
	WmmUpmfUpmps  				*string
	ClosingBy  					*uint
}