package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func InvalidKategori() domain.Error {
	return domain.NotFoundError("RiskAssessment.InvalidKategori", "Kategori risiko wajib diisi")
}

func InvalidPernyataanRisiko() domain.Error {
	return domain.NotFoundError("RiskAssessment.InvalidPernyataanRisiko", "Pernyataan risiko wajib diisi")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("RiskAssessment.NotFound", fmt.Sprintf("Data analisis risiko dengan UUID %s tidak ditemukan", id))
}

func InvalidData() domain.Error {
	return domain.NotFoundError("RiskAssessment.InvalidData", "Data analisis risiko tidak valid")
}

func EmptyData() domain.Error {
	return domain.NotFoundError("RiskAssessment.EmptyData", "Data analisis risiko tidak boleh kosong")
}
