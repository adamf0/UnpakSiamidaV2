package UpdateVerifikasiDokumenV3

type UpdateVerifikasiDokumenV3Command struct {
	DocUUID          string `json:"doc_uuid" form:"doc_uuid"`
	StatusVerifikasi string `json:"status_verifikasi" form:"status_verifikasi"`
	Status           string `json:"status" form:"status"`
	Catatan          string `json:"catatan" form:"catatan"`
	CatatanPenolakan string `json:"catatan_penolakan" form:"catatan_penolakan"`
}
