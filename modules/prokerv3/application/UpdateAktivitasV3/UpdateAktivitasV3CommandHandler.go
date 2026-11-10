package UpdateAktivitasV3

import (
	domain "UnpakSiamida/modules/prokerv3/domain"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type UpdateAktivitasV3CommandHandler struct {
	Repo domain.IProkerV3Repository
}

func (h *UpdateAktivitasV3CommandHandler) Handle(ctx context.Context, cmd UpdateAktivitasV3Command) (*domain.AktivitasV2, error) {
	actUUID, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return nil, errors.New("Invalid aktivitas UUID")
	}

	act, err := h.Repo.GetAktivitasEntityByUuid(ctx, actUUID)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, errors.New("Aktivitas tidak ditemukan")
	}

	rabModified := cmd.RAB != "" || (cmd.Anggaran > 0 && len(act.RabItems) == 0)
	origAnggaran := act.Anggaran
	origRAB := act.RAB
	origRabItems := act.RabItems

	if cmd.Aktivitas != "" {
		act.Aktivitas = strings.TrimSpace(cmd.Aktivitas)
	}
	if cmd.PIC != "" {
		act.PIC = strings.TrimSpace(cmd.PIC)
	}
	if rabModified && cmd.Anggaran >= 0 {
		act.Anggaran = cmd.Anggaran
	}
	if cmd.RAB != "" {
		act.RAB = cmd.RAB
		act.RabItems = nil
	} else if cmd.Anggaran > 0 && len(act.RabItems) == 0 {
		b, _ := json.Marshal([]map[string]interface{}{
			{
				"kategori":   "Umum",
				"rincian":    "Alokasi Anggaran Kegiatan",
				"keterangan": "Alokasi Anggaran Kegiatan",
				"qty_satuan": 1,
				"qty_paket":  1,
				"sbu":        cmd.Anggaran,
				"total":      cmd.Anggaran,
			},
		})
		act.RAB = string(b)
	}
	if cmd.Realisasi != nil {
		act.Realisasi = cmd.Realisasi
	}
	if cmd.CatatanRevisiAnggaran != nil {
		act.CatatanRevisiAnggaran = strings.TrimSpace(*cmd.CatatanRevisiAnggaran)
	}
	if cmd.Pencairan != nil {
		act.Pencairan = cmd.Pencairan
	}
	if cmd.Status != "" {
		act.Status = strings.TrimSpace(cmd.Status)
	}
	if cmd.Capaian != "" {
		act.Capaian = strings.TrimSpace(cmd.Capaian)
	}
	if cmd.Masalah != "" {
		act.Masalah = strings.TrimSpace(cmd.Masalah)
	}
	if cmd.PemecahanMasalah != "" {
		act.PemecahanMasalah = strings.TrimSpace(cmd.PemecahanMasalah)
	}

	if !rabModified {
		act.RAB = ""
		act.RabItems = nil
	}

	if err := h.Repo.UpdateAktivitas(ctx, act); err != nil {
		return nil, err
	}

	if !rabModified {
		act.Anggaran = origAnggaran
		act.RAB = origRAB
		act.RabItems = origRabItems
	}

	return act, nil
}
