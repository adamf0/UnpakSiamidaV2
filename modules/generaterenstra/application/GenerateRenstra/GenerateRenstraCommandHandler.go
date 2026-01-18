package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domaingenerate "UnpakSiamida/modules/generaterenstra/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"

	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type GenerateRenstraCommandHandler struct {
	Repo        domaingenerate.IGenerateRenstraRepository
	RepoRenstra domainrenstra.IRenstraRepository

	RepoFakultasUnit            domainfakultasunit.IFakultasUnitRepository
	RepoTemplateRenstra         domaintemplaterenstra.ITemplateRenstraRepository
	RepoTemplateDokumenTambahan domaintemplatedokumentambahan.ITemplateDokumenTambahanRepository
}

func (h *GenerateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd GenerateRenstraCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// ===============================
	// PARSE UUID
	// ===============================
	renstraUUID, err := uuid.Parse(cmd.UuidRenstra)
	if err != nil {
		return "", domaingenerate.InvalidParsing("Renstra")
	}

	fakultasUnitUUID, err := uuid.Parse(cmd.UuidFakultasUnit)
	if err != nil {
		return "", domaingenerate.InvalidParsing("Fakultas Unit")
	}

	// ===============================
	// LOAD MASTER DATA (PARALLEL)
	// ===============================
	var (
		existingRenstra      *domainrenstra.Renstra
		existingFakultasUnit *domainfakultasunit.FakultasUnit
	)

	g1, gctx1 := errgroup.WithContext(ctx)

	g1.Go(func() error {
		r, err := h.RepoRenstra.GetByUuid(gctx1, renstraUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaingenerate.NotFoundRenstra(cmd.UuidRenstra)
			}
			return err
		}
		existingRenstra = r
		return nil
	})

	g1.Go(func() error {
		r, err := h.RepoFakultasUnit.GetDefaultByUuid(gctx1, fakultasUnitUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaingenerate.NotFoundFakultasUnit(cmd.UuidFakultasUnit)
			}
			return err
		}
		existingFakultasUnit = r
		return nil
	})

	if err := g1.Wait(); err != nil {
		return "", err
	}

	// ===============================
	// SAFETY GUARD (WAJIB)
	// ===============================
	if existingRenstra == nil {
		return "", domaingenerate.NotFoundRenstra(cmd.UuidRenstra)
	}

	if existingFakultasUnit == nil {
		return "", domaingenerate.NotFoundFakultasUnit(cmd.UuidFakultasUnit)
	}

	// ===============================
	// SNAPSHOT VALUE (ANTI-RACE)
	// ===============================
	fakultasUnitID := existingFakultasUnit.ID
	fakultasUnitNama := existingFakultasUnit.Nama
	fakultasUnitType := existingFakultasUnit.Type

	// ===============================
	// LOAD TEMPLATE & AUDIT (PARALLEL)
	// ===============================
	var (
		existingTemplateRenstra         []domaintemplaterenstra.TemplateRenstraDefault
		existingTemplateDokumenTambahan []domaintemplatedokumentambahan.TemplateDokumenTambahanDefault
		existingAuditRenstra            []domaingenerate.GenerateRenstraDefault
		existingAuditDokumenTambahan    []domaingenerate.GenerateDokumenTambahanDefault
	)

	g2, gctx2 := errgroup.WithContext(ctx)

	g2.Go(func() error {
		r, err := h.RepoTemplateRenstra.GetAllByTahunFakUnitDefault(
			gctx2,
			cmd.Tahun,
			fakultasUnitID,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaingenerate.NotFoundTemplate(cmd.Tahun, fakultasUnitNama)
			}
			return err
		}
		existingTemplateRenstra = r
		return nil
	})

	g2.Go(func() error {
		fakKey := fmt.Sprintf("%s#all", fakultasUnitType)

		r, err := h.RepoTemplateDokumenTambahan.GetAllByTahunFakUnitDefault(
			gctx2,
			cmd.Tahun,
			fakKey,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaingenerate.NotFoundTemplate(cmd.Tahun, fakKey)
			}
			return err
		}
		existingTemplateDokumenTambahan = r
		return nil
	})

	g2.Go(func() error {
		r, err := h.Repo.GetAllRenstraNilaiByTahunFakUnitDefault(
			gctx2,
			cmd.Tahun,
			fakultasUnitID,
		)
		if err != nil {
			return err
		}
		existingAuditRenstra = r
		return nil
	})

	g2.Go(func() error {
		r, err := h.Repo.GetAllDokumenTambahanByTahunFakUnitDefault(
			gctx2,
			cmd.Tahun,
			fakultasUnitID,
		)
		if err != nil {
			return err
		}
		existingAuditDokumenTambahan = r
		return nil
	})

	if err := g2.Wait(); err != nil {
		return "", err
	}

	// ===============================
	// MAP & DIFF LOGIC
	// ===============================
	templateRenstraMap := make(map[uint]domaintemplaterenstra.TemplateRenstraDefault)
	var templateRenstraKeys []MatchKey

	for _, tr := range existingTemplateRenstra {
		templateRenstraMap[tr.ID] = tr
		templateRenstraKeys = append(templateRenstraKeys, MatchKey{
			TemplateID: tr.ID,
			Tahun:      tr.Tahun,
			Tugas:      tr.Tugas,
		})
	}

	templateDokumenTambahanMap := make(map[uint]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault)
	var templateDokumenTambahanKeys []MatchKey

	for _, tr := range existingTemplateDokumenTambahan {
		templateDokumenTambahanMap[tr.ID] = tr
		templateDokumenTambahanKeys = append(templateDokumenTambahanKeys, MatchKey{
			TemplateID: tr.ID,
			Tahun:      tr.Tahun,
			Tugas:      tr.Tugas,
		})
	}

	auditRenstraMap := make(map[uint]domaingenerate.GenerateRenstraDefault)
	var auditRenstraKeys []MatchKey

	for _, a := range existingAuditRenstra {
		auditRenstraMap[a.TemplateRenstra] = a
		auditRenstraKeys = append(auditRenstraKeys, MatchKey{
			TemplateID: a.TemplateRenstra,
			Tahun:      a.TemplateTahun,
			Tugas:      a.Tugas,
		})
	}

	auditDokumenTambahanMap := make(map[uint]domaingenerate.GenerateDokumenTambahanDefault)
	var auditDokumenTambahanKeys []MatchKey

	for _, a := range existingAuditDokumenTambahan {
		auditDokumenTambahanMap[a.TemplateDokumenTambahan] = a
		auditDokumenTambahanKeys = append(auditDokumenTambahanKeys, MatchKey{
			TemplateID: a.TemplateDokumenTambahan,
			Tahun:      a.TemplateTahun,
			Tugas:      a.Tugas,
		})
	}
	// printExisting("existingTemplateDokumenTambahan",existingTemplateDokumenTambahan)
	// printExisting("existingAuditDokumenTambahan",existingAuditDokumenTambahan)

	insertRenstras := diffMatchKey(templateRenstraKeys, auditRenstraKeys)
	deleteRenstras := diffMatchKey(auditRenstraKeys, templateRenstraKeys)

	insertDokumenTambahans := diffMatchKey(templateDokumenTambahanKeys, auditDokumenTambahanKeys)
	deleteDokumenTambahans := diffMatchKey(auditDokumenTambahanKeys, templateDokumenTambahanKeys)

	// ===============================
	// TRANSACTION
	// ===============================
	tx, err := h.Repo.BeginTx(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ===============================
	// SYNC Renstra Nilai
	// ===============================
	for _, ins := range insertRenstras {
		tr := templateRenstraMap[ins.TemplateID]

		result := domaingenerate.NewGenerateRenstra(
			ins.Tahun,
			existingRenstra.Tahun,
			fakultasUnitID,
			existingRenstra.FakultasUnit,
			existingRenstra.ID,
			ins.TemplateID,
			tr.UUID.String(),
			tr.Indikator,
			ins.Tugas,
			"insert",
		)

		if !result.IsSuccess {
			tx.Rollback()
			return "", result.Error
		}

		if err := h.Repo.CreateRenstraNilai(ctx, tx, result.Value); err != nil {
			tx.Rollback()
			return "", err
		}
	}

	for _, del := range deleteRenstras {
		ar := auditRenstraMap[del.TemplateID]

		result := domaingenerate.NewGenerateRenstra(
			del.Tahun,
			existingRenstra.Tahun,
			fakultasUnitID,
			existingRenstra.FakultasUnit,
			existingRenstra.ID,
			del.TemplateID,
			ar.TemplateRenstraUuid.String(),
			ar.Indikator,
			del.Tugas,
			"delete",
		)

		if !result.IsSuccess {
			tx.Rollback()
			return "", result.Error
		}

		if err := h.Repo.DeleteRenstraNilai(ctx, tx, result.Value); err != nil {
			tx.Rollback()
			return "", err
		}
	}

	// ===============================
	// SYNC Dokumen Tambahan
	// ===============================
	for _, ins := range insertDokumenTambahans {
		tr := templateDokumenTambahanMap[ins.TemplateID]

		result := domaingenerate.NewGenerateDokumenTambahan(
			ins.Tahun,
			existingRenstra.Tahun,
			fakultasUnitID,
			existingRenstra.FakultasUnit,
			existingRenstra.ID,
			ins.TemplateID,
			tr.UUID.String(),
			tr.JenisFile, //tr.Indikator,
			ins.Tugas,
			"insert",
		)

		if !result.IsSuccess {
			tx.Rollback()
			return "", result.Error
		}

		if err := h.Repo.CreateDokumenTambahan(ctx, tx, result.Value); err != nil {
			tx.Rollback()
			return "", err
		}
	}

	for _, del := range deleteDokumenTambahans {
		ar := auditDokumenTambahanMap[del.TemplateID]

		result := domaingenerate.NewGenerateDokumenTambahan(
			del.Tahun,
			existingRenstra.Tahun,
			fakultasUnitID,
			existingRenstra.FakultasUnit,
			existingRenstra.ID,
			del.TemplateID,
			ar.TemplateDokumenTambahanUuid.String(),
			ar.JenisFile, //ar.Indikator,
			del.Tugas,
			"delete",
		)

		if !result.IsSuccess {
			tx.Rollback()
			return "", result.Error
		}

		if err := h.Repo.DeleteDokumenTambahan(ctx, tx, result.Value); err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return cmd.UuidRenstra, nil
}

type MatchKey struct {
	TemplateID uint
	Tahun      string
	Tugas      string
}

func diffMatchKey(a, b []MatchKey) []MatchKey {
	mb := make(map[MatchKey]bool)
	for _, x := range b {
		mb[x] = true
	}

	var res []MatchKey
	for _, x := range a {
		if !mb[x] {
			res = append(res, x)
		}
	}
	return res
}

// debug helper
func printExisting(tag string, existing interface{}) {
	b, _ := json.MarshalIndent(existing, "", "  ")
	fmt.Println(tag, ":")
	fmt.Println(string(b))
}
