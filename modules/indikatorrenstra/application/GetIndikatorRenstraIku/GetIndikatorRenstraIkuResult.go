package application

import domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"

type GetIndikatorRenstraIkuResult struct {
	IdIndikatorRenstra uint64                                       `json:"id_indikator_renstra"`
	UUID               string                                       `json:"uuid"`
	IkuIds             []uint64                                     `json:"iku_ids"`
	Mappings           []domainindikatorrenstra.IndikatorRenstraIku `json:"mappings"`
}
