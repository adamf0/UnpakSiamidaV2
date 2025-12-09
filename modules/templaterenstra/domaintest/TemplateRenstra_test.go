package domaintest

import (
	"testing"

	"UnpakSiamida/modules/templaterenstra/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// ===============
//  CREATE
// ===============
//

func TestNewTemplateRenstra_Success(t *testing.T) {
	nama := "Standar Penelitian"

	result := domain.NewTemplateRenstra(nama)

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

func TestUpdateTemplateRenstra_Success(t *testing.T) {
	namaAwal := "Standar Lama"
	prev := domain.NewTemplateRenstra(namaAwal).Value

	newNama := "Standar Baru"

	result := domain.UpdateTemplateRenstra(prev, prev.UUID, newNama)

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

func TestUpdateTemplateRenstra_NegativeCases(t *testing.T) {
	valid := domain.NewTemplateRenstra("X").Value

	tests := []struct {
		name        string
		prev        *domain.TemplateRenstra
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

			result := domain.UpdateTemplateRenstra(tt.prev, tt.uid, "abc")

			// Failure → IsSuccess == false
			require.False(t, result.IsSuccess)

			// Error.Description dibanding string
			assert.Equal(t, tt.expectedErr, result.Error.Description)
		})
	}
}
