package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("Swot.EmptyData", "Data SWOT tidak ditemukan")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("Swot.InvalidData", "Data SWOT tidak valid")
}

func InvalidKategori() domain.Error {
	return domain.NotFoundError("Swot.InvalidKategori", "Kategori SWOT harus diisi (Kekuatan, Kelemahan, Peluang, Ancaman)")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("Swot.NotFound", fmt.Sprintf("Data SWOT dengan ID %s tidak ditemukan", id))
}
