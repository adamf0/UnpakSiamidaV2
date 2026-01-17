package domaintest

import (
	"testing"

	domain "UnpakSiamida/modules/dokumentambahan/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateDokumenTambahan_Success(t *testing.T) {
	baseUUID := uuid.New()
	renstraUUID := uuid.New()
	link := "file.pdf"
	capaian := "capaian"
	catatan := "catatan"

	prev := &domain.DokumenTambahan{
		UUID:                    baseUUID,
		Renstra:                 1,
		TemplateDokumenTambahan: 1,
		Link:                    nil,
	}

	renstra := &domainrenstra.Renstra{
		UUID:  renstraUUID,
		Tahun: "2026",
	}

	tests := []struct {
		name            string
		mode            string
		link            *string
		capaian         *string
		catatan         *string
		wantSuccess     bool
		wantCheckFields func(t *testing.T, res *domain.DokumenTambahan)
	}{
		{
			name:        "Success_Auditee",
			mode:        "auditee",
			link:        &link,
			capaian:     nil,
			catatan:     nil,
			wantSuccess: true,
			wantCheckFields: func(t *testing.T, res *domain.DokumenTambahan) {
				assert.Equal(t, &link, res.Link)
				assert.Nil(t, res.CapaianAuditor)
				assert.Nil(t, res.CatatanAuditor)
			},
		},
		{
			name:        "Success_Auditor2",
			mode:        "auditor2",
			link:        nil,
			capaian:     &capaian,
			catatan:     &catatan,
			wantSuccess: true,
			wantCheckFields: func(t *testing.T, res *domain.DokumenTambahan) {
				assert.Nil(t, res.Link)
				assert.Equal(t, &capaian, res.CapaianAuditor)
				assert.Equal(t, &catatan, res.CatatanAuditor)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateDokumenTambahan(prev, renstra, baseUUID, renstraUUID, "2026", tt.mode, "2026#"+tt.mode, tt.link, tt.capaian, tt.catatan)
			require.True(t, res.IsSuccess)
			tt.wantCheckFields(t, res.Value)
		})
	}
}

func TestUpdateDokumenTambahan_FailAndEdgeCases(t *testing.T) {
	baseUUID := uuid.New()
	renstraUUID := uuid.New()
	prev := &domain.DokumenTambahan{
		UUID:                    baseUUID,
		Renstra:                 1,
		TemplateDokumenTambahan: 1,
	}

	validRenstra := &domainrenstra.Renstra{
		UUID:  renstraUUID,
		Tahun: "2026",
	}

	tests := []struct {
		name        string
		prev        *domain.DokumenTambahan
		renstra     *domainrenstra.Renstra
		UUID        uuid.UUID
		UUIDRenstra uuid.UUID
		tahun       string
		mode        string
		granted     string
		wantErrCode string
	}{
		{
			name:        "Fail_PrevNil",
			prev:        nil,
			renstra:     validRenstra,
			UUID:        baseUUID,
			UUIDRenstra: renstraUUID,
			tahun:       "2026",
			mode:        "auditee",
			granted:     "2026#auditee",
			wantErrCode: domain.EmptyData().Code,
		},
		{
			name:        "Fail_RenstraNil",
			prev:        prev,
			renstra:     nil,
			UUID:        baseUUID,
			UUIDRenstra: renstraUUID,
			tahun:       "2026",
			mode:        "auditee",
			granted:     "2026#auditee",
			wantErrCode: domain.InvalidRenstra().Code,
		},
		{
			name:        "Fail_UUIDMismatch",
			prev:        prev,
			renstra:     validRenstra,
			UUID:        uuid.New(),
			UUIDRenstra: renstraUUID,
			tahun:       "2026",
			mode:        "auditee",
			granted:     "2026#auditee",
			wantErrCode: domain.InvalidData().Code,
		},
		{
			name:        "Fail_UUIDRenstraMismatch",
			prev:        prev,
			renstra:     validRenstra,
			UUID:        baseUUID,
			UUIDRenstra: uuid.New(),
			tahun:       "2026",
			mode:        "auditee",
			granted:     "2026#auditee",
			wantErrCode: domain.InvalidData().Code,
		},
		{
			name:        "Fail_TahunMismatch",
			prev:        prev,
			renstra:     validRenstra,
			UUID:        baseUUID,
			UUIDRenstra: renstraUUID,
			tahun:       "2025",
			mode:        "auditee",
			granted:     "2026#auditee",
			wantErrCode: domain.InvalidData().Code,
		},
		{
			name:        "Fail_ModeInvalid",
			prev:        prev,
			renstra:     validRenstra,
			UUID:        baseUUID,
			UUIDRenstra: renstraUUID,
			tahun:       "2026",
			mode:        "invalid_mode",
			granted:     "2026#invalid_mode",
			wantErrCode: domain.RejectAction().Code,
		},
		{
			name:        "Fail_NotGranted",
			prev:        prev,
			renstra:     validRenstra,
			UUID:        baseUUID,
			UUIDRenstra: renstraUUID,
			tahun:       "2026",
			mode:        "auditee",
			granted:     "someotherkey",
			wantErrCode: domain.NotGranted().Code,
		},
		{
			name:        "Edge_GrantedWithSpaces",
			prev:        prev,
			renstra:     validRenstra,
			UUID:        baseUUID,
			UUIDRenstra: renstraUUID,
			tahun:       "2026",
			mode:        "auditee",
			granted:     " 2026#auditee,other#key",
			wantErrCode: "", // success
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateDokumenTambahan(tt.prev, tt.renstra, tt.UUID, tt.UUIDRenstra, tt.tahun, tt.mode, tt.granted, nil, nil, nil)
			if tt.wantErrCode != "" {
				require.False(t, res.IsSuccess)
				assert.Equal(t, tt.wantErrCode, res.Error.Code)
			} else {
				require.True(t, res.IsSuccess)
			}
		})
	}
}
