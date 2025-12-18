package domain

import (
	"github.com/google/uuid"
)

type TemplateRenstraDefault struct {
	ID           	   		uint       
	UUID         	   		uuid.UUID  
	Tahun              		string	
	IndikatorRenstraUuid    uuid.UUID  
	IndikatorRenstraID 		uint	
	Indikator          		string	  
	IsPertanyaan       		bool 	  
	FakultasUnit       		uint 	  
	Kategori           		string 	  
	Klasifikasi        		string	  
	Satuan             		*string	  
	Target             		*string	  
	TargetMin          		*string	  
	TargetMax          		*string	  
	Tugas              		string	  
}