package application

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
)

type RawRabItem struct {
	ID          string      `json:"id"`
	Keterangan  string      `json:"keterangan"`
	QtySatuan   interface{} `json:"qty_satuan"`
	QtyPaket    interface{} `json:"qty_paket"`
	HargaSatuan interface{} `json:"harga_satuan"`
	Total       interface{} `json:"total"`
	BuktiUrl    string      `json:"bukti_url"`
	BuktiName   string      `json:"bukti_name"`
}

func parseRabFloat(val interface{}) float64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		var f float64
		_ = json.Unmarshal([]byte(v), &f)
		return f
	default:
		return 0
	}
}

type GetDashboardIkuAnggaranQueryHandler struct {
	Repo domainindikatorrenstra.IIndikatorRenstraRepository
}

func (h *GetDashboardIkuAnggaranQueryHandler) Handle(
	ctx context.Context,
	q GetDashboardIkuAnggaranQuery,
) (*GetDashboardIkuAnggaranResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	filterTahun := strings.TrimSpace(q.Tahun)
	filterUnit := strings.TrimSpace(q.Unit)

	rawData, err := h.Repo.GetDashboardIkuRawData(ctx, filterTahun, filterUnit)
	if err != nil {
		return nil, err
	}

	ikuToIndikators := make(map[uint64][]domainindikatorrenstra.IkuMapRow)
	for _, mr := range rawData.MapRows {
		ikuToIndikators[mr.IdIku] = append(ikuToIndikators[mr.IdIku], mr)
	}

	indikatorToProkers := make(map[uint64][]domainindikatorrenstra.ProkerDbRow)
	for _, pr := range rawData.ProkerRows {
		indikatorToProkers[pr.IdIndikator] = append(indikatorToProkers[pr.IdIndikator], pr)
	}

	prokerToAktivitas := make(map[uint64][]domainindikatorrenstra.AktivitasDbRow)
	for _, ar := range rawData.AktivitasRows {
		prokerToAktivitas[ar.IdMatrikProker] = append(prokerToAktivitas[ar.IdMatrikProker], ar)
	}

	actToTimelines := make(map[uint64][]domainindikatorrenstra.TimelineDbRow)
	for _, tr := range rawData.TimelineRows {
		actToTimelines[tr.IdAktivitasV2] = append(actToTimelines[tr.IdAktivitasV2], tr)
	}

	var items []DashboardIkuItem
	var summary DashboardIkuSummary
	summary.TotalIku = len(rawData.MasterIkus)

	for _, iku := range rawData.MasterIkus {
		indRows := ikuToIndikators[iku.ID]
		isMapped := len(indRows) > 0
		if isMapped {
			summary.TotalIkuMapped++
		}

		item := DashboardIkuItem{
			IkuID:          iku.ID,
			Iku:            iku.Iku,
			Kategori:       iku.Kategori,
			Satuan:         iku.Satuan,
			JenisIku:       iku.JenisIku,
			Target:         iku.Target,
			IsMapped:       isMapped,
			TotalIndikator: len(indRows),
			RabDetails:     make([]DashboardRabItemDetail, 0),
			IndikatorItems: make([]DashboardIndikatorFlowItem, 0),
		}

		seenProker := make(map[uint64]bool)
		seenAct := make(map[uint64]bool)

		for _, indRow := range indRows {
			indItem := DashboardIndikatorFlowItem{
				ID:            indRow.IdIndikatorRenstra,
				Indikator:     indRow.Indikator,
				Tahun:         indRow.Tahun,
				KegiatanItems: make([]DashboardKegiatanFlowItem, 0),
			}

			prokers := indikatorToProkers[indRow.IdIndikatorRenstra]
			seenIndAct := make(map[uint64]bool)

			for _, p := range prokers {
				if !seenProker[p.ID] {
					seenProker[p.ID] = true
					item.TotalProker++
				}

				acts := prokerToAktivitas[p.ID]
				for _, a := range acts {
					if seenIndAct[a.ID] {
						continue
					}
					seenIndAct[a.ID] = true

					kegItem := DashboardKegiatanFlowItem{
						ID:           a.ID,
						Aktivitas:    a.Aktivitas,
						UnitNama:     p.UnitNama,
						PIC:          p.PIC,
						TargetProker: p.TargetUnpak,
						Anggaran:     a.Anggaran,
						RabItems:     make([]DashboardRabItemDetail, 0),
					}

					if a.Realisasi != nil {
						kegItem.Realisasi = *a.Realisasi
					} else {
						kegItem.Realisasi = a.Anggaran
					}

					if !seenAct[a.ID] {
						seenAct[a.ID] = true
						item.TotalAktivitas++
						item.Anggaran += a.Anggaran
						if a.Realisasi != nil {
							item.Realisasi += *a.Realisasi
						}
					}

					tls := actToTimelines[a.ID]
					for _, tl := range tls {
						item.TotalTimeline++

						var rabList []RawRabItem
						if strings.TrimSpace(tl.RAB) != "" {
							_ = json.Unmarshal([]byte(tl.RAB), &rabList)
						}

						if len(rabList) > 0 {
							for _, r := range rabList {
								rTotal := parseRabFloat(r.Total)
								if rTotal == 0 {
									rTotal = parseRabFloat(r.QtySatuan) * parseRabFloat(r.QtyPaket) * parseRabFloat(r.HargaSatuan)
								}
								kegItem.TotalRab += rTotal

								isCair := strings.TrimSpace(r.BuktiUrl) != "" || strings.EqualFold(tl.Status, "Terlaksana")
								if isCair {
									kegItem.PencairanRab += rTotal
								}

								rabDetail := DashboardRabItemDetail{
									ID:            r.ID,
									Keterangan:    r.Keterangan,
									QtySatuan:     parseRabFloat(r.QtySatuan),
									QtyPaket:      parseRabFloat(r.QtyPaket),
									HargaSatuan:   parseRabFloat(r.HargaSatuan),
									Total:         rTotal,
									BuktiUrl:      r.BuktiUrl,
									BuktiName:     r.BuktiName,
									TimelineNama:  tl.NamaAktivitas,
									StatusTahapan: tl.Status,
									AktivitasNama: a.Aktivitas,
									ProkerTarget:  p.TargetUnpak,
									UnitNama:      p.UnitNama,
								}
								kegItem.RabItems = append(kegItem.RabItems, rabDetail)
							}
						} else if tl.Anggaran > 0 {
							kegItem.TotalRab += tl.Anggaran
							if strings.EqualFold(tl.Status, "Terlaksana") {
								kegItem.PencairanRab += tl.Anggaran
							}
						}
					}

					kegSisa := kegItem.Realisasi - kegItem.PencairanRab
					if kegSisa < 0 {
						kegSisa = 0
					}
					kegItem.SisaDana = kegSisa
					if kegItem.Realisasi > 0 {
						kegItem.Persentase = (kegItem.PencairanRab / kegItem.Realisasi) * 100
					} else if kegItem.Anggaran > 0 {
						kegItem.Persentase = (kegItem.PencairanRab / kegItem.Anggaran) * 100
					}

					indItem.Anggaran += kegItem.Anggaran
					indItem.Realisasi += kegItem.Realisasi
					indItem.TotalRab += kegItem.TotalRab
					indItem.PencairanRab += kegItem.PencairanRab
					indItem.KegiatanItems = append(indItem.KegiatanItems, kegItem)
				}
			}

			indSisa := indItem.Realisasi - indItem.PencairanRab
			if indSisa < 0 {
				indSisa = 0
			}
			indItem.SisaDana = indSisa
			if indItem.Realisasi > 0 {
				indItem.Persentase = (indItem.PencairanRab / indItem.Realisasi) * 100
			} else if indItem.Anggaran > 0 {
				indItem.Persentase = (indItem.PencairanRab / indItem.Anggaran) * 100
			}

			item.IndikatorItems = append(item.IndikatorItems, indItem)
		}

		// Agregasikan detail RAB unik ke level IKU
		seenRabDetailId := make(map[string]bool)
		for _, indItem := range item.IndikatorItems {
			for _, kItem := range indItem.KegiatanItems {
				for _, r := range kItem.RabItems {
					rKey := r.ID + "_" + r.TimelineNama
					if !seenRabDetailId[rKey] {
						seenRabDetailId[rKey] = true
						item.RabDetails = append(item.RabDetails, r)
						item.TotalItemRab++
						item.TotalRab += r.Total
						isCair := strings.TrimSpace(r.BuktiUrl) != "" || strings.EqualFold(r.StatusTahapan, "Terlaksana")
						if isCair {
							item.TotalItemDicairkan++
							item.PencairanRab += r.Total
						}
					}
				}
				if len(kItem.RabItems) == 0 && kItem.TotalRab > 0 {
					item.TotalRab += kItem.TotalRab
					item.PencairanRab += kItem.PencairanRab
					if kItem.PencairanRab > 0 {
						item.TotalItemDicairkan++
					}
				}
			}
		}

		if item.Anggaran > 0 {
			item.PersentaseRealisasi = (item.Realisasi / item.Anggaran) * 100
		}
		if item.Realisasi > 0 {
			item.PersentasePencairan = (item.PencairanRab / item.Realisasi) * 100
		} else if item.Anggaran > 0 {
			item.PersentasePencairan = (item.PencairanRab / item.Anggaran) * 100
		}

		summary.TotalAnggaran += item.Anggaran
		summary.TotalRealisasi += item.Realisasi
		summary.TotalRab += item.TotalRab
		summary.TotalPencairanRab += item.PencairanRab

		items = append(items, item)
	}

	if summary.TotalAnggaran > 0 {
		summary.PersentaseRealisasi = (summary.TotalRealisasi / summary.TotalAnggaran) * 100
		summary.PersentasePencairanAnggaran = (summary.TotalPencairanRab / summary.TotalAnggaran) * 100
	}
	if summary.TotalRealisasi > 0 {
		summary.PersentasePencairanRealisasi = (summary.TotalPencairanRab / summary.TotalRealisasi) * 100
	}

	return &GetDashboardIkuAnggaranResult{
		Summary: summary,
		Items:   items,
		Tahun:   filterTahun,
	}, nil
}
