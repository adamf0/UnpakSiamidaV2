package domaintest

import (
	"testing"

	domain "UnpakSiamida/modules/generaterenstra/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerateRenstra_Success(t *testing.T) {
	validTahun := "2026"
	validFakultasUnit := uint(1)
	validRenstraId := uint(10)
	validTemplate := uint(20)
	validTemplateUUID := "uuid-template"
	indikator := "indikator1"
	validTugas := "Tugas 1"

	res := domain.NewGenerateRenstra(
		validTahun, validTahun, // tahun template & renstra
		validFakultasUnit, validFakultasUnit, // fakultasUnit template & renstra
		validRenstraId,
		validTemplate,
		validTemplateUUID,
		indikator,
		validTugas,
		"insert",
	)

	require.True(t, res.IsSuccess)
	assert.Equal(t, validRenstraId, res.Value.RenstraId)
	assert.Equal(t, validTemplate, res.Value.TemplateRenstra)
	assert.Equal(t, validTugas, res.Value.Tugas)
	assert.NotEqual(t, domain.GenerateRenstra{}.UUID, res.Value.UUID) // UUID harus terisi
}

func TestNewGenerateRenstra_FailAndEdgeCases(t *testing.T) {
	validTahun := "2026"
	invalidTahun := "2025"
	validFakultasUnit := uint(1)
	invalidFakultasUnit := uint(2)
	validRenstraId := uint(10)
	invalidRenstraId := uint(0)
	validTemplate := uint(20)
	invalidTemplate := uint(0)
	validTemplateUUID := "uuid-template"
	indikator := "indikator1"
	validTugas := "Tugas 1"
	invalidTugas := "" // dianggap tidak valid

	tests := []struct {
		name                string
		tahun               string
		renstraTahun        string
		fakultasUnit        uint
		renstraFakultasUnit uint
		renstraId           uint
		template            uint
		tugas               string
		operation           string
		wantErrCode         string
	}{
		{
			name:                "InvalidTahun",
			tahun:               validTahun,
			renstraTahun:        invalidTahun,
			fakultasUnit:        validFakultasUnit,
			renstraFakultasUnit: validFakultasUnit,
			renstraId:           validRenstraId,
			template:            validTemplate,
			tugas:               validTugas,
			operation:           "insert",
			wantErrCode:         domain.InvalidTahunRenstra(validTemplateUUID, indikator, validTahun, invalidTahun, "insert").Code,
		},
		{
			name:                "InvalidFakultasUnit",
			tahun:               validTahun,
			renstraTahun:        validTahun,
			fakultasUnit:        validFakultasUnit,
			renstraFakultasUnit: invalidFakultasUnit,
			renstraId:           validRenstraId,
			template:            validTemplate,
			tugas:               validTugas,
			operation:           "insert",
			wantErrCode:         domain.InvalidFakultasUnit().Code,
		},
		{
			name:                "InvalidTemplate",
			tahun:               validTahun,
			renstraTahun:        validTahun,
			fakultasUnit:        validFakultasUnit,
			renstraFakultasUnit: validFakultasUnit,
			renstraId:           validRenstraId,
			template:            invalidTemplate,
			tugas:               validTugas,
			operation:           "insert",
			wantErrCode:         domain.InvalidTemplate().Code,
		},
		{
			name:                "InvalidRenstraId",
			tahun:               validTahun,
			renstraTahun:        validTahun,
			fakultasUnit:        validFakultasUnit,
			renstraFakultasUnit: validFakultasUnit,
			renstraId:           invalidRenstraId,
			template:            validTemplate,
			tugas:               validTugas,
			operation:           "insert",
			wantErrCode:         domain.InvalidRenstra().Code,
		},
		{
			name:                "InvalidTugas",
			tahun:               validTahun,
			renstraTahun:        validTahun,
			fakultasUnit:        validFakultasUnit,
			renstraFakultasUnit: validFakultasUnit,
			renstraId:           validRenstraId,
			template:            validTemplate,
			tugas:               invalidTugas,
			operation:           "insert",
			wantErrCode:         domain.InvalidTugas().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewGenerateRenstra(
				tt.tahun,
				tt.renstraTahun,
				tt.fakultasUnit,
				tt.renstraFakultasUnit,
				tt.renstraId,
				tt.template,
				validTemplateUUID,
				indikator,
				tt.tugas,
				tt.operation,
			)

			require.False(t, res.IsSuccess)
			assert.Equal(t, tt.wantErrCode, res.Error.Code)
		})
	}
}
