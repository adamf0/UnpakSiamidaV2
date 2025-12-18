package domain

import (
	"github.com/google/uuid"
)

type TemplateDokumenTambahanDefault struct {
	ID          		uint           
	UUID        		uuid.UUID      
	Tahun       		string
	JenisFileID 		uint
	JenisFileUuid       uuid.UUID     
	JenisFile   		string  
	FakultasProdiUnit	string         
	Pertanyaan  		string         
	Klasifikasi 		string
	Tugas       		string
}