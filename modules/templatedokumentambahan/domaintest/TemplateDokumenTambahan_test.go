package domaintest

import (
	"testing"

	"UnpakSiamida/modules/templatedokumentambahan/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// ===============
//  CREATE
// ===============
//

func TestNewTemplateDokumenTambahan_Success(t *testing.T) {
	nama := "Standar Penelitian"

	result := domain.NewTemplateDokumenTambahan(nama)

	require.True(t, result.IsSuccess)

	v := result.Value
	assert.Equal(t, nama, v.Nama)
	assert.NotEqual(t, uuid.Nil, v.UUID)
}

//
// ===================
//  UPDATE (POSITIVE)
// ===================
//

func TestUpdateTemplateDokumenTambahan_Success(t *testing.T) {
	namaAwal := "Standar Lama"
	prev := domain.NewTemplateDokumenTambahan(namaAwal).Value

	newNama := "Standar Baru"

	result := domain.UpdateTemplateDokumenTambahan(prev, prev.UUID, newNama)

	require.True(t, result.IsSuccess)

	require.NotNil(t, result.Value)
	assert.Equal(t, newNama, result.Value.Nama)
}

//
// ====================
//  UPDATE (NEGATIVE)
//  THEORY TABLE TESTS
// ====================
//

func TestUpdateTemplateDokumenTambahan_NegativeCases(t *testing.T) {
	valid := domain.NewTemplateDokumenTambahan("X").Value

	tests := []struct {
		name        string
		prev        *domain.TemplateDokumenTambahan
		uid         uuid.UUID
		expectedErr string
	}{
		{
			name:        "prev is nil → EmptyData",
			prev:        nil,
			uid:         uuid.New(),
			expectedErr: domain.EmptyData().Description,
		},
		{
			name:        "uuid mismatch → InvalidData",
			prev:        valid,
			uid:         uuid.New(), // beda UUID
			expectedErr: domain.InvalidData().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := domain.UpdateTemplateDokumenTambahan(tt.prev, tt.uid, "abc")

			// Failure → IsSuccess == false
			require.False(t, result.IsSuccess)

			// Error.Description dibanding string
			assert.Equal(t, tt.expectedErr, result.Error.Description)
		})
	}
}
