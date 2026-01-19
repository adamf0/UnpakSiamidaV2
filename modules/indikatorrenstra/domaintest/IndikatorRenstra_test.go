package domaintest

import (
	domain "UnpakSiamida/modules/indikatorrenstra/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewIndikatorRenstra_Success(t *testing.T) {
	standar := uint(1)
	parent := uint(10)
	operator := "SUM"

	res := domain.NewIndikatorRenstra(
		"Jumlah mahasiswa",
		&standar,
		&parent,
		"2025",
		"Peningkatan",
		&operator,
		true,
	)

	require.True(t, res.IsSuccess)
	ir := res.Value
	require.NotNil(t, ir)

	assert.NotEqual(t, uuid.Nil, ir.UUID)
	assert.Equal(t, "Jumlah mahasiswa", ir.Indikator)
	assert.Equal(t, "2025", ir.Tahun)
	assert.Equal(t, "Peningkatan", ir.TipeTarget)
	assert.Equal(t, &operator, ir.Operator)
	assert.Equal(t, &standar, ir.StandarRenstra)
	assert.Equal(t, &parent, ir.Parent)
}

func NewIndikatorRenstra_Fail(t *testing.T) {
	standar := uint(1)

	tests := []struct {
		name            string
		standar         *uint
		isUnique        bool
		expectedErrDesc string
	}{
		{
			name:            "StandarNil → InvalidStandar",
			standar:         nil,
			isUnique:        true,
			expectedErrDesc: domain.InvalidStandar().Description,
		},
		{
			name:            "NotUnique → NotUniqueIndikator",
			standar:         &standar,
			isUnique:        false,
			expectedErrDesc: domain.NotUniqueIndikator().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewIndikatorRenstra(
				"Test Indikator",
				tt.standar,
				nil,
				"2025",
				"A",
				nil,
				tt.isUnique,
			)
			require.False(t, res.IsSuccess)
			require.NotNil(t, res.Error)
			assert.Equal(t, tt.expectedErrDesc, res.Error.Description)
		})
	}
}

func UpdateIndikatorRenstra_Success(t *testing.T) {
	standar := uint(1)
	parent := uint(2)
	operator := "AVG"

	prevRes := domain.NewIndikatorRenstra(
		"Mahasiswa baru",
		&standar,
		&parent,
		"2024",
		"Target",
		&operator,
		true,
	)
	require.True(t, prevRes.IsSuccess)

	prev := prevRes.Value
	prevUUID := prev.UUID

	newStandar := uint(88)
	newParent := uint(99)
	newOperator := "SUM"

	res := domain.UpdateIndikatorRenstra(
		prev,
		prevUUID,
		"Indikator Update",
		&newStandar,
		&newParent,
		"2026",
		"Perubahan",
		&newOperator,
	)

	require.True(t, res.IsSuccess)
	ir := res.Value
	assert.Equal(t, "Indikator Update", ir.Indikator)
	assert.Equal(t, "2026", ir.Tahun)
	assert.Equal(t, "Perubahan", ir.TipeTarget)
	assert.Equal(t, &newOperator, ir.Operator)
	assert.Equal(t, &newStandar, ir.StandarRenstra)
	assert.Equal(t, &newParent, ir.Parent)
}

func UpdateIndikatorRenstra_Fail(t *testing.T) {
	standar := uint(1)

	prevRes := domain.NewIndikatorRenstra(
		"Awal",
		&standar,
		nil,
		"2024",
		"A",
		nil,
		true,
	)
	require.True(t, prevRes.IsSuccess)
	validPrev := prevRes.Value

	tests := []struct {
		name        string
		prev        *domain.IndikatorRenstra
		uid         uuid.UUID
		standar     *uint
		expectedErr string
	}{
		{
			name:        "StandarNil → InvalidStandar",
			prev:        validPrev,
			uid:         validPrev.UUID,
			standar:     nil,
			expectedErr: domain.InvalidStandar().Description,
		},
		{
			name:        "PrevNil → EmptyData",
			prev:        nil,
			uid:         uuid.New(),
			standar:     &standar,
			expectedErr: domain.EmptyData().Description,
		},
		{
			name:        "UUIDMismatch → InvalidData",
			prev:        validPrev,
			uid:         uuid.New(),
			standar:     &standar,
			expectedErr: domain.InvalidData().Description,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateIndikatorRenstra(
				tt.prev,
				tt.uid,
				"Test Update",
				tt.standar,
				nil,
				"2025",
				"X",
				nil,
			)

			require.False(t, res.IsSuccess)
			require.NotNil(t, res.Error)
			assert.Equal(t, tt.expectedErr, res.Error.Description)
		})
	}
}

func IndikatorRenstra_EdgeCases(t *testing.T) {
	standar := uint(1)

	// Create dengan parent dan operator nil
	res := domain.NewIndikatorRenstra(
		"Edge Indikator",
		&standar,
		nil,
		"",
		"",
		nil,
		true,
	)
	require.True(t, res.IsSuccess)
	ir := res.Value
	assert.Nil(t, ir.Parent)
	assert.Nil(t, ir.Operator)
	assert.Equal(t, "", ir.Tahun)
	assert.Equal(t, "", ir.TipeTarget)

	// Update tanpa mengubah operator dan parent
	newTahun := "2030"
	newTipe := "Edge"
	resUpdate := domain.UpdateIndikatorRenstra(
		ir,
		ir.UUID,
		"Edge Update",
		&standar,
		nil,
		newTahun,
		newTipe,
		nil,
	)
	require.True(t, resUpdate.IsSuccess)
	assert.Nil(t, resUpdate.Value.Operator)
	assert.Nil(t, resUpdate.Value.Parent)
	assert.Equal(t, newTahun, resUpdate.Value.Tahun)
	assert.Equal(t, newTipe, resUpdate.Value.TipeTarget)
	assert.Equal(t, "Edge Update", resUpdate.Value.Indikator)
}
