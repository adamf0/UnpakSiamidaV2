package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ReverseRenstraOperator membalikkan operator untuk range (misal >= menjadi <=)
func ReverseRenstraOperator(op string) string {
	op = strings.TrimSpace(op)
	switch op {
	case ">=":
		return "<="
	case "<=":
		return ">="
	case ">":
		return "<"
	case "<":
		return ">"
	case "≥":
		return "≤"
	case "≤":
		return "≥"
	case "=":
		return "="
	default:
		return op
	}
}

// FormatRenstraTargetUnit memformat target unit dari template renstra x master indikator renstra
// Sesuai aturan:
// - Numerik: master_indikator_renstra.operator + template_renstra.target + template_renstra.satuan
// - Kategori: template_renstra.target
// - Range: template_renstra.target_min + master_indikator_renstra.operator + master_indikator_renstra.operator (reverse) + template_renstra.target_max
func FormatRenstraTargetUnit(tipeTarget, operator, target, targetMin, targetMax, satuan string) string {
	tipeTarget = strings.ToLower(strings.TrimSpace(tipeTarget))
	operator = strings.TrimSpace(operator)
	target = strings.TrimSpace(target)
	targetMin = strings.TrimSpace(targetMin)
	targetMax = strings.TrimSpace(targetMax)
	satuan = strings.TrimSpace(satuan)

	switch tipeTarget {
	case "numerik":
		parts := make([]string, 0, 3)
		if operator != "" {
			parts = append(parts, operator)
		}
		if target != "" {
			parts = append(parts, target)
		}
		if satuan != "" {
			parts = append(parts, satuan)
		}
		return strings.Join(parts, " ")

	case "kategori":
		return target

	case "range":
		rev := ReverseRenstraOperator(operator)
		if targetMin != "" && targetMax != "" {
			if operator != "" {
				res := fmt.Sprintf("%s %s - %s %s", operator, targetMin, rev, targetMax)
				if satuan != "" {
					res += " " + satuan
				}
				return strings.TrimSpace(res)
			}
			res := fmt.Sprintf("%s - %s", targetMin, targetMax)
			if satuan != "" {
				res += " " + satuan
			}
			return strings.TrimSpace(res)
		} else if targetMin != "" {
			res := targetMin
			if operator != "" {
				res = operator + " " + res
			}
			if satuan != "" {
				res += " " + satuan
			}
			return strings.TrimSpace(res)
		} else if targetMax != "" {
			res := targetMax
			if rev != "" {
				res = rev + " " + res
			}
			if satuan != "" {
				res += " " + satuan
			}
			return strings.TrimSpace(res)
		}
		return ""

	default:
		// Fallback jika tipe target tidak eksplisit
		if target != "" {
			if operator != "" {
				return strings.TrimSpace(operator + " " + target + " " + satuan)
			}
			return strings.TrimSpace(target + " " + satuan)
		}
		if targetMin != "" && targetMax != "" {
			return strings.TrimSpace(targetMin + " - " + targetMax + " " + satuan)
		}
		return ""
	}
}

// ComputeProkerVerificationStatus menghitung status verifikator dan status anggaran proker
// berdasarkan agregasi aktivitasv2 dan status_anggaran matrix_proker.
func ComputeProkerVerificationStatus(totalAct, terimaVerif, tolakVerif, diajukanVerif, terimaAngg, tolakAngg int64) (statusVerifikator, statusAnggaran string) {
	if totalAct == 0 {
		return "draf evaluator", "menunggu_verifikator"
	}

	// Status Verifikator (Evaluator)
	if terimaVerif == totalAct {
		statusVerifikator = "diterima evaluator"
	} else if terimaVerif > 0 {
		statusVerifikator = fmt.Sprintf("diterima evaluator (parsial %d/%d)", terimaVerif, totalAct)
	} else if tolakVerif == totalAct {
		statusVerifikator = "ditolak evaluator"
	} else if tolakVerif > 0 {
		statusVerifikator = "revisi evaluator"
	} else if diajukanVerif > 0 {
		statusVerifikator = "diajukan evaluator"
	} else {
		statusVerifikator = "draf evaluator"
	}

	// Status Anggaran
	if totalAct > 0 && terimaAngg == totalAct && terimaVerif == totalAct {
		statusAnggaran = "terima anggaran"
	} else if terimaAngg > 0 {
		statusAnggaran = fmt.Sprintf("terima anggaran (parsial %d/%d)", terimaAngg, totalAct)
	} else if terimaVerif == totalAct {
		statusAnggaran = "menunggu anggaran"
	} else {
		statusAnggaran = "menunggu_verifikator"
	}

	return statusVerifikator, statusAnggaran
}

// SyncMatrixProkerStatus menghitung ulang dan mengupdate matrix_proker.status dan status_anggaran
// berdasarkan agregasi aktivitasv2 yang berelasi.
// Verifikasi anggaran dilakukan pada program kerja (matrix_proker), bukan per baris aktivitas.
func SyncMatrixProkerStatus(ctx context.Context, db *gorm.DB, idMatrixProker uint) (newStatus, newStatusAnggaran string, err error) {
	if idMatrixProker == 0 || db == nil {
		return "", "", nil
	}

	type actStat struct {
		ID                uint       `gorm:"column:id"`
		StatusVerifikator string     `gorm:"column:status_verifikator"`
		StatusEvaluator   string     `gorm:"column:status_evaluator"`
		TerimaAnggaranAt  *time.Time `gorm:"column:terima_anggaran_at"`
		IsBatal           bool       `gorm:"column:is_batal"`
	}

	var acts []actStat
	err = db.WithContext(ctx).Table("aktivitasv2").
		Select("id, status_verifikator, status_evaluator, terima_anggaran_at, is_batal").
		Where("id_matrik_proker = ? AND COALESCE(is_batal, 0) = 0", idMatrixProker).
		Find(&acts).Error
	if err != nil {
		return "", "", err
	}

	totalAct := len(acts)
	var currentMp struct {
		Status         string `gorm:"column:status"`
		StatusAnggaran string `gorm:"column:status_anggaran"`
	}
	_ = db.WithContext(ctx).Table("matrix_proker").
		Select("status, status_anggaran").
		Where("id = ?", idMatrixProker).
		Take(&currentMp).Error

	currentStatus := strings.ToLower(strings.TrimSpace(currentMp.Status))
	currentStatusAnggaran := strings.ToLower(strings.TrimSpace(currentMp.StatusAnggaran))

	if totalAct == 0 {
		if currentStatus == "" || currentStatus == "draft" {
			newStatus = "draft"
		} else {
			newStatus = currentStatus
		}
		newStatusAnggaran = "menunggu_verifikator"
	} else {
		var terimaEval, tolakEval, draftCount, diajukanCount, terimaAngg int

		for _, a := range acts {
			sEval := strings.ToLower(strings.TrimSpace(a.StatusEvaluator))
			sVerif := strings.ToLower(strings.TrimSpace(a.StatusVerifikator))

			isApprovedEval := sEval == "disetujui" || sEval == "diterima" || sEval == "terima_verifikator" || sVerif == "diterima" || sVerif == "disetujui"
			isRejectedEval := sEval == "ditolak" || sVerif == "ditolak" || sEval == "revisi" || sVerif == "revisi"

			if isApprovedEval {
				terimaEval++
			} else if isRejectedEval {
				tolakEval++
			} else if sEval == "draft" || sEval == "" || sVerif == "draft" {
				draftCount++
			} else {
				diajukanCount++
			}

			if a.TerimaAnggaranAt != nil {
				terimaAngg++
			}
		}

		// Kalkulasi Status Evaluator Proker (matrix_proker.status)
		if terimaEval == totalAct {
			newStatus = "diterima"
		} else if terimaEval > 0 {
			newStatus = "diterima_parsial"
		} else if tolakEval == totalAct {
			newStatus = "ditolak"
		} else if currentStatus == "draft" || currentStatus == "" || draftCount == totalAct {
			// Jika belum pernah diajukan (draft) dan belum ada yang diterima, tetap draft!
			newStatus = "draft"
		} else {
			newStatus = "diajukan"
		}

		// Kalkulasi Status Anggaran Proker (matrix_proker.status_anggaran)
		// Verifikasi anggaran dilakukan pada tingkat program kerja (matrix_proker) dan ditandai oleh terima_anggaran_at di aktivitasv2.
		if totalAct > 0 && terimaAngg == totalAct && (newStatus == "diterima" || currentStatus == "diterima" || currentStatus == "terima_anggaran") {
			newStatusAnggaran = "diterima"
		} else if terimaAngg > 0 {
			newStatusAnggaran = "diterima_parsial"
		} else {
			if currentStatusAnggaran == "ditolak" {
				newStatusAnggaran = "ditolak"
			} else if currentStatusAnggaran == "revisi" {
				newStatusAnggaran = "revisi"
			} else if newStatus == "diterima" {
				newStatusAnggaran = "menunggu_anggaran"
			} else {
				newStatusAnggaran = "menunggu_verifikator"
			}
		}
	}

	// Update langsung ke database matrix_proker jika berubah
	if currentStatus != newStatus || currentStatusAnggaran != newStatusAnggaran {
		err = db.WithContext(ctx).Table("matrix_proker").
			Where("id = ?", idMatrixProker).
			Updates(map[string]interface{}{
				"status":          newStatus,
				"status_anggaran": newStatusAnggaran,
				"updated_at":      time.Now(),
			}).Error
	}

	return newStatus, newStatusAnggaran, err
}
