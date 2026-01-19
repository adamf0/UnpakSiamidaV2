package domaintest

import (
	"testing"

	domain "UnpakSiamida/modules/kts/domain"

	"github.com/google/uuid"
)

func NewKtsRenstra(t *testing.T) {
	tests := []struct {
		name        string
		isDataExist bool
		expectFail  bool
	}{
		{"Success case", false, false},
		{"Fail case: data exists", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewKtsRenstra(nil, nil, nil, nil, nil, "2026", 1, "target", tt.isDataExist)
			if tt.expectFail {
				if res.IsSuccess {
					t.Errorf("expected failure but got success")
				}
			} else {
				if !res.IsSuccess {
					t.Errorf("expected success but got failure")
				}
			}
		})
	}
}

func NewKtsDokumen(t *testing.T) {
	tests := []struct {
		name        string
		isDataExist bool
		expectFail  bool
	}{
		{"Success case", false, false},
		{"Fail case: data exists", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewKtsDokumen(nil, nil, nil, nil, nil, "2026", 1, "target", tt.isDataExist)
			if tt.expectFail {
				if res.IsSuccess {
					t.Errorf("expected failure but got success")
				}
			} else {
				if !res.IsSuccess {
					t.Errorf("expected success but got failure")
				}
			}
		})
	}
}

// -------------------- Update Step1 --------------------
func UpdateKtsStep1(t *testing.T) {
	validUUID := uuid.New()
	prev := &domain.Kts{UUID: validUUID}
	prevKts := &domain.KtsDefault{Tahun: domain.StringPtr("2026")}

	tests := []struct {
		name           string
		prev           *domain.Kts
		prevKts        *domain.KtsDefault
		uid            uuid.UUID
		nomorLaporan   string
		tanggalLaporan string
		accAuditor     uint
		expectFail     bool
	}{
		{"Fail: prev nil", nil, prevKts, validUUID, "001", "2026-01-01", 1, true},
		{"Fail: prevKts nil", prev, nil, validUUID, "001", "2026-01-01", 1, true},
		{"Fail: UUID mismatch", prev, prevKts, uuid.New(), "001", "2026-01-01", 1, true},
		{"Fail: tahun mismatch", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2025")}, validUUID, "001", "2026-01-01", 1, true},
		{"Fail: accAuditor 0", prev, prevKts, validUUID, "001", "2026-01-01", 0, true},
		{"Fail: nomorLaporan empty", prev, prevKts, validUUID, "   ", "2026-01-01", 1, true},
		{"Fail: tanggal invalid", prev, prevKts, validUUID, "001", "invalid", 1, true},
		{"Success case", prev, prevKts, validUUID, "001", "2026-01-01", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateKtsStep1(tt.prev, tt.prevKts, tt.uid, tt.nomorLaporan, tt.tanggalLaporan,
				"P", "L", "O", "R", "akar", "koreksi", tt.accAuditor, "2026")
			if tt.expectFail {
				if res.IsSuccess {
					t.Errorf("expected failure but got success")
				}
			} else {
				if !res.IsSuccess {
					t.Errorf("expected success but got failure")
				}
			}
		})
	}
}

// -------------------- Update Step2 --------------------
func UpdateKtsStep2(t *testing.T) {
	validUUID := uuid.New()
	prev := &domain.Kts{UUID: validUUID}
	prevKts := &domain.KtsDefault{Tahun: domain.StringPtr("2026")}

	keterangan := "tolak alasan"

	tests := []struct {
		name             string
		prev             *domain.Kts
		prevKts          *domain.KtsDefault
		uid              uuid.UUID
		statusAccAuditee uint
		accAuditee       uint
		keteranganTolak  *string
		expectFail       bool
	}{
		{"Fail: prev nil", nil, prevKts, validUUID, 1, 1, &keterangan, true},
		{"Fail: prevKts nil", prev, nil, validUUID, 1, 1, &keterangan, true},
		{"Fail: UUID mismatch", prev, prevKts, uuid.New(), 1, 1, &keterangan, true},
		{"Fail: tahun mismatch", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2025")}, validUUID, 1, 1, &keterangan, true},
		{"Fail: accAuditee 0", prev, prevKts, validUUID, 1, 0, &keterangan, true},
		{"Fail: statusAccAuditee >1", prev, prevKts, validUUID, 2, 1, &keterangan, true},
		{"Fail: statusAccAuditee 0 but keterangan nil", prev, prevKts, validUUID, 0, 1, nil, true},
		{"Fail: statusAccAuditee 0 but keterangan empty", prev, prevKts, validUUID, 0, 1, domain.StringPtr(" "), true},
		{"Success case acc 1", prev, prevKts, validUUID, 1, 1, nil, false},
		{"Success case acc 0 with keterangan", prev, prevKts, validUUID, 0, 1, &keterangan, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateKtsStep2(tt.prev, tt.prevKts, tt.uid, tt.statusAccAuditee,
				tt.accAuditee, tt.keteranganTolak, nil, "2026")
			if tt.expectFail {
				if res.IsSuccess {
					t.Errorf("expected failure but got success")
				}
			} else {
				if !res.IsSuccess {
					t.Errorf("expected success but got failure")
				}
			}
		})
	}
}

// -------------------- UpdateKtsTindakan --------------------
func UpdateKtsTindakan(t *testing.T) {
	validUUID := uuid.New()
	prev := &domain.Kts{UUID: validUUID}
	prevKts := &domain.KtsDefault{Tahun: domain.StringPtr("2026")}

	tests := []struct {
		name       string
		prev       *domain.Kts
		prevKts    *domain.KtsDefault
		uid        uuid.UUID
		tindakan   string
		expectFail bool
	}{
		{"Fail: prev nil", nil, prevKts, validUUID, "tindakan", true},
		{"Fail: prevKts nil", prev, nil, validUUID, "tindakan", true},
		{"Fail: UUID mismatch", prev, prevKts, uuid.New(), "tindakan", true},
		{"Fail: tahun mismatch", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2025")}, validUUID, "tindakan", true},
		{"Success case", prev, prevKts, validUUID, "tindakan", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateKtsTindakan(tt.prev, tt.prevKts, tt.uid, tt.tindakan, "2026")
			if tt.expectFail {
				if res.IsSuccess {
					t.Errorf("expected failure but got success")
				}
			} else {
				if !res.IsSuccess {
					t.Errorf("expected success but got failure")
				}
			}
		})
	}
}

// -------------------- Update Step3 --------------------
func UpdateKtsStep3(t *testing.T) {
	validUUID := uuid.New()
	prev := &domain.Kts{UUID: validUUID}
	prevKts := &domain.KtsDefault{Tahun: domain.StringPtr("2026"), Auditor: domain.StringPtr("1")}

	tests := []struct {
		name                string
		prev                *domain.Kts
		prevKts             *domain.KtsDefault
		uid                 uuid.UUID
		accAuditor          uint
		tanggalPenyelesaian string
		expectFail          bool
	}{
		{"Fail: prev nil", nil, prevKts, validUUID, 1, "2026-01-01", true},
		{"Fail: prevKts nil", prev, nil, validUUID, 1, "2026-01-01", true},
		{"Fail: UUID mismatch", prev, prevKts, uuid.New(), 1, "2026-01-01", true},
		{"Fail: tahun mismatch", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2025"), Auditor: domain.StringPtr("1")}, validUUID, 1, "2026-01-01", true},
		{"Fail: accAuditor 0", prev, prevKts, validUUID, 0, "2026-01-01", true},
		{"Fail: accAuditor mismatch Auditor field", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2026"), Auditor: domain.StringPtr("2")}, validUUID, 1, "2026-01-01", true},
		{"Fail: tanggal invalid", prev, prevKts, validUUID, 1, "invalid", true},
		{"Success case", prev, prevKts, validUUID, 1, "2026-01-01", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateKtsStep3(tt.prev, tt.prevKts, tt.uid, tt.accAuditor, tt.tanggalPenyelesaian, "2026")
			if tt.expectFail {
				if res.IsSuccess {
					t.Errorf("expected failure but got success")
				}
			} else {
				if !res.IsSuccess {
					t.Errorf("expected success but got failure")
				}
			}
		})
	}
}

// -------------------- Update Step4 --------------------
func UpdateKtsStep4(t *testing.T) {
	validUUID := uuid.New()
	prev := &domain.Kts{UUID: validUUID}
	prevKts := &domain.KtsDefault{Tahun: domain.StringPtr("2026"), Auditee: domain.StringPtr("1")}

	tests := []struct {
		name           string
		prev           *domain.Kts
		prevKts        *domain.KtsDefault
		uid            uuid.UUID
		tinjauan       string
		tanggalClosing string
		accFinal       uint
		expectFail     bool
	}{
		{"Fail: prev nil", nil, prevKts, validUUID, "tinjauan", "2026-01-01", 1, true},
		{"Fail: prevKts nil", prev, nil, validUUID, "tinjauan", "2026-01-01", 1, true},
		{"Fail: UUID mismatch", prev, prevKts, uuid.New(), "tinjauan", "2026-01-01", 1, true},
		{"Fail: tahun mismatch", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2025"), Auditee: domain.StringPtr("1")}, validUUID, "tinjauan", "2026-01-01", 1, true},
		{"Fail: accFinal 0", prev, prevKts, validUUID, "tinjauan", "2026-01-01", 0, true},
		{"Fail: accFinal mismatch Auditee field", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2026"), Auditee: domain.StringPtr("2")}, validUUID, "tinjauan", "2026-01-01", 1, true},
		{"Fail: tanggal invalid", prev, prevKts, validUUID, "tinjauan", "invalid", 1, true},
		{"Success case", prev, prevKts, validUUID, "tinjauan", "2026-01-01", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateKtsStep4(tt.prev, tt.prevKts, tt.uid, tt.tinjauan, tt.tanggalClosing, tt.accFinal, "2026")
			if tt.expectFail {
				if res.IsSuccess {
					t.Errorf("expected failure but got success")
				}
			} else {
				if !res.IsSuccess {
					t.Errorf("expected success but got failure")
				}
			}
		})
	}
}

// -------------------- Update Step5 --------------------
func UpdateKtsStep5(t *testing.T) {
	validUUID := uuid.New()
	prev := &domain.Kts{UUID: validUUID}
	prevKts := &domain.KtsDefault{Tahun: domain.StringPtr("2026"), Auditor: domain.StringPtr("1")}

	tests := []struct {
		name                string
		prev                *domain.Kts
		prevKts             *domain.KtsDefault
		uid                 uuid.UUID
		tanggalClosingFinal string
		wmmUpmfUpmps        string
		closingBy           uint
		expectFail          bool
	}{
		{"Fail: prev nil", nil, prevKts, validUUID, "2026-01-01", "wmm", 1, true},
		{"Fail: prevKts nil", prev, nil, validUUID, "2026-01-01", "wmm", 1, true},
		{"Fail: UUID mismatch", prev, prevKts, uuid.New(), "2026-01-01", "wmm", 1, true},
		{"Fail: tahun mismatch", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2025"), Auditor: domain.StringPtr("1")}, validUUID, "2026-01-01", "wmm", 1, true},
		{"Fail: closingBy 0", prev, prevKts, validUUID, "2026-01-01", "wmm", 0, true},
		{"Fail: closingBy mismatch Auditor field", prev, &domain.KtsDefault{Tahun: domain.StringPtr("2026"), Auditor: domain.StringPtr("2")}, validUUID, "2026-01-01", "wmm", 1, true},
		{"Fail: tanggal invalid", prev, prevKts, validUUID, "invalid", "wmm", 1, true},
		{"Success case", prev, prevKts, validUUID, "2026-01-01", "wmm", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateKtsStep5(tt.prev, tt.prevKts, tt.uid, tt.tanggalClosingFinal, tt.wmmUpmfUpmps, tt.closingBy, "2026")
			if tt.expectFail {
				if res.IsSuccess {
					t.Errorf("expected failure but got success")
				}
			} else {
				if !res.IsSuccess {
					t.Errorf("expected success but got failure")
				}
			}
		})
	}
}
