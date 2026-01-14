package domain

import (
	"errors"
	"strconv"
	"strings"
	"time"

	// "github.com/goforj/godump"

	common "UnpakSiamida/common/domain"
	event "UnpakSiamida/modules/kts/event"

	"github.com/google/uuid"
)

type Kts struct {
	common.Entity

	ID   uint      `gorm:"primaryKey;autoIncrement"`
	UUID uuid.UUID `gorm:"type:char(36);uniqueIndex"`

	RenstraNilai    *uint  `gorm:"column:id_renstra_nilai"`
	DokumenTambahan *uint  `gorm:"column:id_dokumen_tambahan"`
	Status          string `gorm:"column:status"`

	// auditor
	NomorLaporan     *string    `gorm:"column:nomor_laporan"`
	TanggalLaporan   *time.Time `gorm:"column:tanggal_laporan;type:date"`
	Auditor          *string    `gorm:"column:auditor"`
	KetidaksesuaianP *string    `gorm:"column:uraian_ketidaksesuaian_p"`
	KetidaksesuaianL *string    `gorm:"column:uraian_ketidaksesuaian_l"`
	KetidaksesuaianO *string    `gorm:"column:uraian_ketidaksesuaian_o"`
	KetidaksesuaianR *string    `gorm:"column:uraian_ketidaksesuaian_r"`
	AkarMasalah      *string    `gorm:"column:akar_masalah"`
	TindakanKoreksi  *string    `gorm:"column:tindakan_koreksi"`
	AccAuditor       *uint      `gorm:"column:acc_auditor"`

	// auditee
	StatusAccAuditee  *uint   `gorm:"column:status_acc_auditee"`
	AccAuditee        *uint   `gorm:"column:acc_auditee"`
	KeteranganTolak   *string `gorm:"column:keterangan_tolak_auditee"`
	TindakanPerbaikan *string `gorm:"column:tindakan_perbaikan"`

	// auditor
	TanggalPenyelesaian *time.Time `gorm:"column:tanggal_penyelesaian;type:date"`

	// auditee
	TinjauanTindakanPerbaikan *string    `gorm:"column:tinjauan_tindakan_perbaikan"`
	TanggalClosing            *time.Time `gorm:"column:tanggal_closing_auditee;type:date"`
	AccFinal                  *uint      `gorm:"column:acc_auditor_final"`

	// auditor
	TanggalClosingFinal *time.Time `gorm:"column:tanggal_closing;type:date"`
	WmmUpmfUpmps        *string    `gorm:"column:wmm_upmf_upmps"`
	ClosingBy           *uint      `gorm:"column:closingBy"`
}

func (Kts) TableName() string {
	return "kts_renstra"
}

// === CREATE ===
// [pr] belum dipakai
func NewKtsRenstra(
	auditor *string,
	renstraNilai *uint,
	templateRenstra *uint,
	standar *string,
	indikator *string,
	tahun string,
	idTarget uint,
	target string,
	isDataExist bool,
) common.ResultValue[*Kts] {
	if isDataExist {
		return common.FailureValue[*Kts](ExistData())
	}

	kts := &Kts{
		UUID:         uuid.New(),
		RenstraNilai: renstraNilai,
		Auditor:      auditor,
		Status:       "draf",
	}

	kts.Raise(event.KtsCreatedEvent{
		EventID:         uuid.New(),
		OccurredOn:      time.Now().UTC(),
		KtsUUID:         kts.UUID,
		TemplateRenstra: templateRenstra,
		Standar:         standar,
		Indikator:       indikator,
		Tahun:           tahun,
		IdTarget:        idTarget,
		Target:          target,
	})

	return common.SuccessValue(kts)
}

// [pr] belum dipakai
func NewKtsDokumen(
	auditor *string,
	dokumenTambahan *uint,
	templateDokumen *uint,
	pertanyaan *string,
	jenisFile *string,
	tahun string,
	idTarget uint,
	target string,
	isDataExist bool,
) common.ResultValue[*Kts] {
	if isDataExist {
		return common.FailureValue[*Kts](ExistData())
	}

	kts := &Kts{
		UUID:            uuid.New(),
		DokumenTambahan: dokumenTambahan,
		Auditor:         auditor,
		Status:          "draf",
	}

	kts.Raise(event.KtsCreatedEvent{
		EventID:         uuid.New(),
		OccurredOn:      time.Now().UTC(),
		KtsUUID:         kts.UUID,
		TemplateDokumen: templateDokumen,
		Pertanyaan:      pertanyaan,
		JenisFile:       jenisFile,
		Tahun:           tahun,
		IdTarget:        idTarget,
		Target:          target,
		Status:          "draf",
	})

	return common.SuccessValue(kts)
}

//
// ======================= UPDATE STEP 1 =======================
//

func UpdateKtsStep1(
	prev *Kts,
	prevKts *KtsDefault,
	uid uuid.UUID,
	nomorLaporan string,
	tanggalLaporan string,
	ketidaksesuaianP string,
	ketidaksesuaianL string,
	ketidaksesuaianO string,
	ketidaksesuaianR string,
	akarMasalah string,
	tindakanKoreksi string,
	accAuditor uint,
	tahun string,
) common.ResultValue[*Kts] {

	if prev == nil || prevKts == nil {
		return common.FailureValue[*Kts](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*Kts](InvalidData())
	}
	if prevKts.Tahun == nil || *prevKts.Tahun != tahun {
		return common.FailureValue[*Kts](InvalidTahun())
	}
	if accAuditor == 0 {
		return common.FailureValue[*Kts](InvalidAuditor())
	}
	if strings.TrimSpace(nomorLaporan) == "" {
		return common.FailureValue[*Kts](RequiredNomorLaporan())
	}

	tgl, err := ParseDatePtr(tanggalLaporan)
	if err != nil {
		return common.FailureValue[*Kts](InvalidTanggal())
	}

	prev.NomorLaporan = StringPtr(nomorLaporan)
	prev.TanggalLaporan = tgl
	prev.KetidaksesuaianP = StringPtr(ketidaksesuaianP)
	prev.KetidaksesuaianL = StringPtr(ketidaksesuaianL)
	prev.KetidaksesuaianO = StringPtr(ketidaksesuaianO)
	prev.KetidaksesuaianR = StringPtr(ketidaksesuaianR)
	prev.AkarMasalah = StringPtr(akarMasalah)
	prev.TindakanKoreksi = StringPtr(tindakanKoreksi)
	prev.AccAuditor = UintPtr(accAuditor)
	prev.Status = "menunggu_verif_auditee"

	if prev.DokumenTambahan != nil {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prev.UUID,
			TemplateDokumen: prevKts.DokumenTambahan,
			Pertanyaan:      prevKts.Pertanyaan,
			JenisFileId:     prevKts.JenisFileId,
			JenisFile:       prevKts.JenisFile,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     StringPtr(nomorLaporan),
			TanggalLaporan:   TimeToStringPtr(prev.TanggalLaporan),
			KetidaksesuaianP: StringPtr(ketidaksesuaianP),
			KetidaksesuaianL: StringPtr(ketidaksesuaianL),
			KetidaksesuaianO: StringPtr(ketidaksesuaianO),
			KetidaksesuaianR: StringPtr(ketidaksesuaianR),
			AkarMasalah:      StringPtr(akarMasalah),
			TindakanKoreksi:  StringPtr(tindakanKoreksi),
			Status:           StringPtr("menunggu_verif_auditee"),
		})
	} else {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prevKts.UUID,
			TemplateRenstra: prevKts.TemplateRenstra,
			StandarId:       prevKts.StandarId,
			Standar:         prevKts.Standar,
			IndikatorId:     prevKts.IndikatorId,
			Indikator:       prevKts.Indikator,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     StringPtr(nomorLaporan),
			TanggalLaporan:   TimeToStringPtr(prev.TanggalLaporan),
			KetidaksesuaianP: StringPtr(ketidaksesuaianP),
			KetidaksesuaianL: StringPtr(ketidaksesuaianL),
			KetidaksesuaianO: StringPtr(ketidaksesuaianO),
			KetidaksesuaianR: StringPtr(ketidaksesuaianR),
			AkarMasalah:      StringPtr(akarMasalah),
			TindakanKoreksi:  StringPtr(tindakanKoreksi),
			Status:           StringPtr("menunggu_verif_auditee"),
		})
	}

	return common.SuccessValue(prev)
}

//
// ======================= UPDATE STEP 2 =======================
//

func UpdateKtsStep2(
	prev *Kts,
	prevKts *KtsDefault,
	uid uuid.UUID,
	statusAccAuditee uint,
	accAuditee uint,
	keteranganTolak *string,
	tindakanPerbaikan *string,
	tahun string,
) common.ResultValue[*Kts] {

	if prev == nil || prevKts == nil {
		return common.FailureValue[*Kts](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*Kts](InvalidData())
	}
	if prevKts.Tahun == nil || *prevKts.Tahun != tahun {
		return common.FailureValue[*Kts](InvalidTahun())
	}
	if accAuditee == 0 {
		return common.FailureValue[*Kts](InvalidAuditee())
	}
	if statusAccAuditee > 1 {
		return common.FailureValue[*Kts](InvalidStatusAcc())
	}
	if statusAccAuditee == 0 && (keteranganTolak == nil || strings.TrimSpace(*keteranganTolak) == "") {
		return common.FailureValue[*Kts](RequiredKeteranganTolak())
	}

	prev.StatusAccAuditee = UintPtr(statusAccAuditee)
	prev.AccAuditee = UintPtr(accAuditee)
	prev.KeteranganTolak = keteranganTolak
	prev.TindakanPerbaikan = tindakanPerbaikan
	if statusAccAuditee == 1 {
		prev.Status = "terima_auditee"
	} else {
		prev.Status = "tolak_auditee"
	}

	if prev.DokumenTambahan != nil {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prev.UUID,
			TemplateDokumen: prevKts.DokumenTambahan,
			Pertanyaan:      prevKts.Pertanyaan,
			JenisFile:       prevKts.JenisFile,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: prev.TindakanPerbaikan,

			Status: StringPtr(prev.Status),
		})
	} else {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prevKts.UUID,
			TemplateRenstra: prevKts.TemplateRenstra,
			StandarId:       prevKts.StandarId,
			Standar:         prevKts.Standar,
			IndikatorId:     prevKts.IndikatorId,
			Indikator:       prevKts.Indikator,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: prev.TindakanPerbaikan,

			Status: StringPtr(prev.Status),
		})
	}

	return common.SuccessValue(prev)
}

func UpdateKtsTindakan(
	prev *Kts,
	prevKts *KtsDefault,
	uid uuid.UUID,
	tindakanPerbaikan string,
	tahun string,
) common.ResultValue[*Kts] {

	if prev == nil || prevKts == nil {
		return common.FailureValue[*Kts](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*Kts](InvalidData())
	}
	if prevKts.Tahun == nil || *prevKts.Tahun != tahun {
		return common.FailureValue[*Kts](InvalidTahun())
	}

	prev.TindakanPerbaikan = StringPtr(tindakanPerbaikan)

	if prev.DokumenTambahan != nil {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prev.UUID,
			TemplateDokumen: prevKts.DokumenTambahan,
			Pertanyaan:      prevKts.Pertanyaan,
			JenisFile:       prevKts.JenisFile,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: StringPtr(tindakanPerbaikan),
			Status:            StringPtr(prevKts.Status),
		})
	} else {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prevKts.UUID,
			TemplateRenstra: prevKts.TemplateRenstra,
			StandarId:       prevKts.StandarId,
			Standar:         prevKts.Standar,
			IndikatorId:     prevKts.IndikatorId,
			Indikator:       prevKts.Indikator,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: StringPtr(tindakanPerbaikan),
			Status:            StringPtr(prevKts.Status),
		})
	}

	return common.SuccessValue(prev)
}

//
// ======================= UPDATE STEP 3 =======================
//

func UpdateKtsStep3(
	prev *Kts,
	prevKts *KtsDefault,
	uid uuid.UUID,
	accAuditor uint,
	tanggalPenyelesaian string,
	tahun string,
) common.ResultValue[*Kts] {

	if prev == nil || prevKts == nil {
		return common.FailureValue[*Kts](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*Kts](InvalidData())
	}
	if prevKts.Tahun == nil || *prevKts.Tahun != tahun {
		return common.FailureValue[*Kts](InvalidTahun())
	}
	if accAuditor == 0 || prevKts.Auditor == nil || *prevKts.Auditor != strconv.FormatUint(uint64(accAuditor), 10) {
		return common.FailureValue[*Kts](InvalidAuditor())
	}

	tgl, err := ParseDatePtr(tanggalPenyelesaian)
	if err != nil {
		return common.FailureValue[*Kts](InvalidTanggal())
	}

	prev.TanggalPenyelesaian = tgl
	prev.Status = "menunggu_closing_auditor"

	if prev.DokumenTambahan != nil {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prev.UUID,
			TemplateDokumen: prevKts.DokumenTambahan,
			Pertanyaan:      prevKts.Pertanyaan,
			JenisFile:       prevKts.JenisFile,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: prev.TindakanPerbaikan,

			TanggalPenyelesaian: TimeToStringPtr(prev.TanggalPenyelesaian),
			Status:              StringPtr("menunggu_closing_auditor"),
		})
	} else {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prevKts.UUID,
			TemplateRenstra: prevKts.TemplateRenstra,
			StandarId:       prevKts.StandarId,
			Standar:         prevKts.Standar,
			IndikatorId:     prevKts.IndikatorId,
			Indikator:       prevKts.Indikator,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: prev.TindakanPerbaikan,

			TanggalPenyelesaian: TimeToStringPtr(prev.TanggalPenyelesaian),
			Status:              StringPtr("menunggu_closing_auditor"),
		})
	}

	return common.SuccessValue(prev)
}

//
// ======================= UPDATE STEP 4 =======================
//

func UpdateKtsStep4(
	prev *Kts,
	prevKts *KtsDefault,
	uid uuid.UUID,
	tinjauan string,
	tanggalClosing string,
	accFinal uint,
	tahun string,
) common.ResultValue[*Kts] {

	if prev == nil || prevKts == nil {
		return common.FailureValue[*Kts](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*Kts](InvalidData())
	}
	if prevKts.Tahun == nil || *prevKts.Tahun != tahun {
		return common.FailureValue[*Kts](InvalidTahun())
	}
	if accFinal == 0 || prevKts.Auditee == nil || *prevKts.Auditee != strconv.FormatUint(uint64(accFinal), 10) {
		return common.FailureValue[*Kts](InvalidAuditee())
	}

	tgl, err := ParseDatePtr(tanggalClosing)
	if err != nil {
		return common.FailureValue[*Kts](InvalidTanggal())
	}

	prev.TinjauanTindakanPerbaikan = StringPtr(tinjauan)
	prev.TanggalClosing = tgl
	prev.AccFinal = UintPtr(accFinal)
	prev.Status = "closing_auditor"

	if prev.DokumenTambahan != nil {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prev.UUID,
			TemplateDokumen: prevKts.DokumenTambahan,
			Pertanyaan:      prevKts.Pertanyaan,
			JenisFile:       prevKts.JenisFile,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: prev.TindakanPerbaikan,

			TanggalPenyelesaian: TimeToStringPtr(prev.TanggalPenyelesaian),

			TanggalClosing: TimeToStringPtr(prev.TanggalClosing),
			Status:         StringPtr("closing_auditor"),
		})
	} else {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prevKts.UUID,
			TemplateRenstra: prevKts.TemplateRenstra,
			StandarId:       prevKts.StandarId,
			Standar:         prevKts.Standar,
			IndikatorId:     prevKts.IndikatorId,
			Indikator:       prevKts.Indikator,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: prev.TindakanPerbaikan,

			TanggalPenyelesaian: TimeToStringPtr(prev.TanggalPenyelesaian),

			TanggalClosing: TimeToStringPtr(prev.TanggalClosing),
			Status:         StringPtr("closing_auditor"),
		})
	}

	return common.SuccessValue(prev)
}

//
// ======================= UPDATE STEP 5 =======================
//

func UpdateKtsStep5(
	prev *Kts,
	prevKts *KtsDefault,
	uid uuid.UUID,
	tanggalClosingFinal string,
	wmmUpmfUpmps string,
	closingBy uint,
	tahun string,
) common.ResultValue[*Kts] {
	if prev == nil || prevKts == nil {
		return common.FailureValue[*Kts](EmptyData())
	}
	if prev.UUID != uid {
		return common.FailureValue[*Kts](InvalidData())
	}
	if prevKts.Tahun == nil || *prevKts.Tahun != tahun {
		return common.FailureValue[*Kts](InvalidTahun())
	}
	if closingBy == 0 || prevKts.Auditor == nil || *prevKts.Auditor != strconv.FormatUint(uint64(closingBy), 10) {
		return common.FailureValue[*Kts](InvalidAuditor())
	}

	tgl, err := ParseDatePtr(tanggalClosingFinal)
	if err != nil {
		return common.FailureValue[*Kts](InvalidTanggal())
	}

	prev.TanggalClosingFinal = tgl
	prev.WmmUpmfUpmps = StringPtr(wmmUpmfUpmps)
	prev.ClosingBy = UintPtr(closingBy)
	prev.Status = "tutup_kts"

	if prev.DokumenTambahan != nil {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prev.UUID,
			TemplateDokumen: prevKts.DokumenTambahan,
			Pertanyaan:      prevKts.Pertanyaan,
			JenisFile:       prevKts.JenisFile,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: prev.TindakanPerbaikan,

			TanggalPenyelesaian: TimeToStringPtr(prev.TanggalPenyelesaian),

			TanggalClosing: TimeToStringPtr(prev.TanggalClosing),

			TanggalClosingFinal: TimeToStringPtr(prev.TanggalClosingFinal),
			WmmUpmfUpmps:        StringPtr(wmmUpmfUpmps),
			Status:              StringPtr("tutup_kts"),
		})
	} else {
		prev.Raise(event.KtsUpdatedEvent{
			EventID:         uuid.New(),
			OccurredOn:      time.Now().UTC(),
			KtsUUID:         prevKts.UUID,
			TemplateRenstra: prevKts.TemplateRenstra,
			StandarId:       prevKts.StandarId,
			Standar:         prevKts.Standar,
			IndikatorId:     prevKts.IndikatorId,
			Indikator:       prevKts.Indikator,
			Tahun:           prevKts.Tahun,
			IdTarget:        prevKts.IdTarget,
			Target:          prevKts.Target,

			NomorLaporan:     prevKts.NomorLaporan,
			TanggalLaporan:   prevKts.TanggalLaporan,
			KetidaksesuaianP: prevKts.KetidaksesuaianP,
			KetidaksesuaianL: prevKts.KetidaksesuaianL,
			KetidaksesuaianO: prevKts.KetidaksesuaianO,
			KetidaksesuaianR: prevKts.KetidaksesuaianR,
			AkarMasalah:      prevKts.AkarMasalah,
			TindakanKoreksi:  prevKts.TindakanKoreksi,

			StatusAccAuditee:  prev.StatusAccAuditee,
			KeteranganTolak:   prev.KeteranganTolak,
			TindakanPerbaikan: prev.TindakanPerbaikan,

			TanggalPenyelesaian: TimeToStringPtr(prev.TanggalPenyelesaian),

			TanggalClosing: TimeToStringPtr(prev.TanggalClosing),

			TanggalClosingFinal: TimeToStringPtr(prev.TanggalClosingFinal),
			WmmUpmfUpmps:        StringPtr(wmmUpmfUpmps),
			Status:              StringPtr("tutup_kts"),
		})
	}

	return common.SuccessValue(prev)
}

//
// ======================= HELPERS =======================
//

func StringPtr(v string) *string { return &v }
func UintPtr(v uint) *uint       { return &v }

func ParseDatePtr(input string) (*time.Time, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil
	}
	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, input); err == nil {
			tt := t.Truncate(24 * time.Hour)
			return &tt, nil
		}
	}
	return nil, errors.New("invalid date format")
}

func TimeToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}
