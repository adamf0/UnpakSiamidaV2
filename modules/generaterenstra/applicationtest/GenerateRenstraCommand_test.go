package applicationtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	commondomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/fakultasunit/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	application "UnpakSiamida/modules/generaterenstra/application/GenerateRenstra"
	domaingenerate "UnpakSiamida/modules/generaterenstra/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	domaintemplatedokumentambahan "UnpakSiamida/modules/templatedokumentambahan/domain"
	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
)

type MockGenerateRenstraRepo struct {
	mock.Mock
}

func (m *MockGenerateRenstraRepo) BeginTx(ctx context.Context) (*gorm.DB, error) {
	args := m.Called(ctx)
	return nil, args.Error(1)
}

func (m *MockGenerateRenstraRepo) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockGenerateRenstraRepo) Rollback() {
	m.Called()
}

func (m *MockGenerateRenstraRepo) GetAllRenstraNilaiByTahunFakUnitDefault(ctx context.Context, tahun string, fakultasUnit uint) ([]domaingenerate.GenerateRenstraDefault, error) {
	args := m.Called(ctx, tahun, fakultasUnit)
	return args.Get(0).([]domaingenerate.GenerateRenstraDefault), args.Error(1)
}

func (m *MockGenerateRenstraRepo) GetAllDokumenTambahanByTahunFakUnitDefault(ctx context.Context, tahun string, fakultasUnit uint) ([]domaingenerate.GenerateDokumenTambahanDefault, error) {
	args := m.Called(ctx, tahun, fakultasUnit)
	return args.Get(0).([]domaingenerate.GenerateDokumenTambahanDefault), args.Error(1)
}

func (m *MockGenerateRenstraRepo) CreateRenstraNilai(ctx context.Context, tx *gorm.DB, generateRenstra *domaingenerate.GenerateRenstra) error {
	args := m.Called(ctx, tx, generateRenstra)
	return args.Error(0)
}

func (m *MockGenerateRenstraRepo) DeleteRenstraNilai(ctx context.Context, tx *gorm.DB, generateRenstra *domaingenerate.GenerateRenstra) error {
	args := m.Called(ctx, tx, generateRenstra)
	return args.Error(0)
}

func (m *MockGenerateRenstraRepo) ForceDeleteRenstraNilai(ctx context.Context, uid uuid.UUID, renstra uint) error {
	args := m.Called(ctx, uid, renstra)
	return args.Error(0)
}

func (m *MockGenerateRenstraRepo) CreateDokumenTambahan(ctx context.Context, tx *gorm.DB, generateDokumenTambahan *domaingenerate.GenerateDokumenTambahan) error {
	args := m.Called(ctx, tx, generateDokumenTambahan)
	return args.Error(0)
}

func (m *MockGenerateRenstraRepo) DeleteDokumenTambahan(ctx context.Context, tx *gorm.DB, generateDokumenTambahan *domaingenerate.GenerateDokumenTambahan) error {
	args := m.Called(ctx, tx, generateDokumenTambahan)
	return args.Error(0)
}

func (m *MockGenerateRenstraRepo) ForceDeleteDokumenTambahan(ctx context.Context, uid uuid.UUID, renstra uint) error {
	args := m.Called(ctx, uid, renstra)
	return args.Error(0)
}

// ==================== Mock: IRenstraRepository ====================
type MockRenstraRepo struct {
	mock.Mock
}

func (m *MockRenstraRepo) IsUnique(ctx context.Context, fakultas_unit uint, tahun string) (bool, error) {
	args := m.Called(ctx, fakultas_unit, tahun)
	return args.Bool(0), args.Error(1)
}

func (m *MockRenstraRepo) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainrenstra.Renstra, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*domainrenstra.Renstra), args.Error(1)
}

func (m *MockRenstraRepo) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domainrenstra.RenstraDefault, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*domainrenstra.RenstraDefault), args.Error(1)
}

func (m *MockRenstraRepo) GetAll(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int, scope string) ([]domainrenstra.RenstraDefault, int64, error) {
	args := m.Called(ctx, search, searchFilters, page, limit, scope)
	return args.Get(0).([]domainrenstra.RenstraDefault), args.Get(1).(int64), args.Error(2)
}

func (m *MockRenstraRepo) Create(ctx context.Context, renstra *domainrenstra.Renstra) error {
	args := m.Called(ctx, renstra)
	return args.Error(0)
}

func (m *MockRenstraRepo) Update(ctx context.Context, renstra *domainrenstra.Renstra) error {
	args := m.Called(ctx, renstra)
	return args.Error(0)
}

func (m *MockRenstraRepo) Delete(ctx context.Context, uid uuid.UUID) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockRenstraRepo) SetupUuid(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ==================== Mock: IFakultasUnitRepository ====================
type MockFakultasUnitRepo struct {
	mock.Mock
}

func (m *MockFakultasUnitRepo) GetDefaultByUuid(ctx context.Context, uid uuid.UUID) (*domainfakultasunit.FakultasUnit, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*domainfakultasunit.FakultasUnit), args.Error(1)
}

func (m *MockFakultasUnitRepo) GetAll(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int) ([]domainfakultasunit.FakultasUnit, int64, error) {
	args := m.Called(ctx, search, searchFilters, page, limit)
	return args.Get(0).([]domainfakultasunit.FakultasUnit), args.Get(1).(int64), args.Error(2)
}

func (m *MockFakultasUnitRepo) SetupUuid(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ==================== Mock: ITemplateRenstraRepository ====================
type MockTemplateRenstraRepo struct {
	mock.Mock
}

func (m *MockTemplateRenstraRepo) GetByUuid(ctx context.Context, uid uuid.UUID) (*domaintemplaterenstra.TemplateRenstra, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*domaintemplaterenstra.TemplateRenstra), args.Error(1)
}

func (m *MockTemplateRenstraRepo) GetAll(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int) ([]domaintemplaterenstra.TemplateRenstraDefault, int64, error) {
	args := m.Called(ctx, search, searchFilters, page, limit)
	return args.Get(0).([]domaintemplaterenstra.TemplateRenstraDefault), args.Get(1).(int64), args.Error(2)
}

func (m *MockTemplateRenstraRepo) GetAllByTahunFakUnit(ctx context.Context, tahun string, fakultasUnit uint) ([]domaintemplaterenstra.TemplateRenstra, error) {
	args := m.Called(ctx, tahun, fakultasUnit)
	return args.Get(0).([]domaintemplaterenstra.TemplateRenstra), args.Error(1)
}

func (m *MockTemplateRenstraRepo) GetAllByTahunFakUnitDefault(ctx context.Context, tahun string, fakultasUnit uint) ([]domaintemplaterenstra.TemplateRenstraDefault, error) {
	args := m.Called(ctx, tahun, fakultasUnit)
	return args.Get(0).([]domaintemplaterenstra.TemplateRenstraDefault), args.Error(1)
}

func (m *MockTemplateRenstraRepo) Create(ctx context.Context, templaterenstra *domaintemplaterenstra.TemplateRenstra) error {
	args := m.Called(ctx, templaterenstra)
	return args.Error(0)
}

func (m *MockTemplateRenstraRepo) Update(ctx context.Context, templaterenstra *domaintemplaterenstra.TemplateRenstra) error {
	args := m.Called(ctx, templaterenstra)
	return args.Error(0)
}

func (m *MockTemplateRenstraRepo) Delete(ctx context.Context, uid uuid.UUID) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockTemplateRenstraRepo) SetupUuid(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ==================== Mock: ITemplateDokumenTambahanRepository ====================
type MockTemplateDokTambahanRepo struct {
	mock.Mock
}

func (m *MockTemplateDokTambahanRepo) GetByUuid(ctx context.Context, uid uuid.UUID) (*domaintemplatedokumentambahan.TemplateDokumenTambahan, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*domaintemplatedokumentambahan.TemplateDokumenTambahan), args.Error(1)
}

func (m *MockTemplateDokTambahanRepo) GetAll(ctx context.Context, search string, searchFilters []commondomain.SearchFilter, page, limit *int) ([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault, int64, error) {
	args := m.Called(ctx, search, searchFilters, page, limit)
	return args.Get(0).([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault), args.Get(1).(int64), args.Error(2)
}

func (m *MockTemplateDokTambahanRepo) GetAllByTahunFakUnitDefault(ctx context.Context, tahun string, fakKey string) ([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault, error) {
	args := m.Called(ctx, tahun, fakKey)
	return args.Get(0).([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault), args.Error(1)
}

func (m *MockTemplateDokTambahanRepo) Create(ctx context.Context, doc *domaintemplatedokumentambahan.TemplateDokumenTambahan) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockTemplateDokTambahanRepo) Update(ctx context.Context, doc *domaintemplatedokumentambahan.TemplateDokumenTambahan) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockTemplateDokTambahanRepo) Delete(ctx context.Context, uid uuid.UUID) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockTemplateDokTambahanRepo) SetupUuid(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ==================== Test Validation ====================
func TestGenerateRenstraCommandValidation(t *testing.T) {
	validCmd := application.GenerateRenstraCommand{
		Tahun:            "2026",
		UuidRenstra:      uuid.New().String(),
		UuidFakultasUnit: uuid.New().String(),
	}
	err := application.GenerateRenstraCommandValidation(validCmd)
	assert.NoError(t, err)

	invalidCmd := application.GenerateRenstraCommand{
		Tahun:            "",
		UuidRenstra:      "invalid",
		UuidFakultasUnit: "1234",
	}
	err = application.GenerateRenstraCommandValidation(invalidCmd)
	assert.Error(t, err)
}

// ==================== Test Command Handler Success ====================
func TestGenerateRenstraCommand_Success(t *testing.T) {
	ctx := context.Background()
	renstraUUID := uuid.New()
	fakultasUUID := uuid.New()

	mockRepo := new(MockGenerateRenstraRepo)
	mockFak := new(MockFakultasUnitRepo)
	mockRenstra := new(MockRenstraRepo)
	mockTemplateRenstra := new(MockTemplateRenstraRepo)
	mockTemplateDok := new(MockTemplateDokTambahanRepo)

	// Mock return values
	mockRenstra.On("GetByUuid", mock.Anything, renstraUUID).Return(&domainrenstra.Renstra{
		ID:           1,
		UUID:         renstraUUID,
		Tahun:        "2026",
		FakultasUnit: 1,
	}, nil)

	mockFak.On("GetDefaultByUuid", mock.Anything, fakultasUUID).Return(&domain.FakultasUnit{
		ID:   1,
		UUID: fakultasUUID,
		Nama: "Fakultas A",
		Type: "F",
	}, nil)

	mockTemplateRenstra.On("GetAllByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).Return([]domaintemplaterenstra.TemplateRenstraDefault{
		{
			ID:    1,
			UUID:  uuid.New(),
			Tugas: "tugas1",
		},
	}, nil)

	mockTemplateDok.On("GetAllByTahunFakUnitDefault", mock.Anything, "2026", "F#all").Return([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault{
		{
			ID:    1,
			UUID:  uuid.New(),
			Tugas: "tugas1",
		},
	}, nil)

	mockRepo.On("GetAllRenstraNilaiByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).
		Return([]domaingenerate.GenerateRenstraDefault{}, nil)
	mockRepo.On("GetAllDokumenTambahanByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).
		Return([]domaingenerate.GenerateDokumenTambahanDefault{}, nil)

	handler := &application.GenerateRenstraCommandHandler{
		Repo:                        mockRepo,
		RepoRenstra:                 mockRenstra,
		RepoFakultasUnit:            mockFak,
		RepoTemplateRenstra:         mockTemplateRenstra,
		RepoTemplateDokumenTambahan: mockTemplateDok,
	}

	cmd := application.GenerateRenstraCommand{
		Tahun:            "2026",
		UuidRenstra:      renstraUUID.String(),
		UuidFakultasUnit: fakultasUUID.String(),
	}

	result, err := handler.Handle(ctx, cmd)
	assert.NoError(t, err)
	assert.Equal(t, renstraUUID.String(), result)
}

// ==================== Test Command Handler Fail NotFound ====================
func TestGenerateRenstraCommand_NotFound(t *testing.T) {
	ctx := context.Background()
	renstraUUID := uuid.New()
	fakultasUUID := uuid.New()

	mockRenstra := new(MockRenstraRepo)
	mockRenstra.On("GetByUuid", mock.Anything, renstraUUID).Return(nil, gorm.ErrRecordNotFound)

	handler := &application.GenerateRenstraCommandHandler{
		RepoRenstra: mockRenstra,
	}

	cmd := application.GenerateRenstraCommand{
		Tahun:            "2026",
		UuidRenstra:      renstraUUID.String(),
		UuidFakultasUnit: fakultasUUID.String(),
	}

	result, err := handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "Renstra not found")
}

func TestGenerateRenstraCommand_AllRules(t *testing.T) {
	ctx := context.Background()
	renstraUUID := uuid.New()
	fakultasUUID := uuid.New()

	// -------------------
	// Mock semua repo
	// -------------------
	mockRenstra := new(MockRenstraRepo)
	mockFak := new(MockFakultasUnitRepo)
	mockTemplateRenstra := new(MockTemplateRenstraRepo)
	mockTemplateDok := new(MockTemplateDokTambahanRepo)
	mockRepo := new(MockGenerateRenstraRepo)

	// handler
	handler := &application.GenerateRenstraCommandHandler{
		RepoRenstra:                 mockRenstra,
		RepoFakultasUnit:            mockFak,
		RepoTemplateRenstra:         mockTemplateRenstra,
		RepoTemplateDokumenTambahan: mockTemplateDok,
		Repo:                        mockRepo,
	}

	// common mock: fakultasUnit valid
	mockFak.On("GetDefaultByUuid", mock.Anything, fakultasUUID).Return(&domainfakultasunit.FakultasUnit{
		ID: 1, UUID: fakultasUUID, Nama: "Fakultas A", Type: "F",
	}, nil)

	// common mock: template renstra valid
	mockTemplateRenstra.On("GetAllByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).
		Return([]domaintemplaterenstra.TemplateRenstraDefault{{ID: 1, UUID: uuid.New(), Tugas: "auditor1"}}, nil)
	mockTemplateDok.On("GetAllByTahunFakUnitDefault", mock.Anything, "2026", "F#all").
		Return([]domaintemplatedokumentambahan.TemplateDokumenTambahanDefault{{ID: 1, UUID: uuid.New(), Tugas: "auditor1"}}, nil)

	// -------------------
	// CASE 1: Tahun berbeda → InvalidTahunRenstra
	// -------------------
	cmd1 := application.GenerateRenstraCommand{
		Tahun:            "2025", // berbeda dari renstra tahun mock
		UuidRenstra:      renstraUUID.String(),
		UuidFakultasUnit: fakultasUUID.String(),
	}
	mockRenstra.On("GetByUuid", mock.Anything, renstraUUID).Return(&domainrenstra.Renstra{
		ID: 1, UUID: renstraUUID, Tahun: "2026", FakultasUnit: 1,
	}, nil)

	_, err := handler.Handle(ctx, cmd1)
	assert.Error(t, err)
	dErr1, ok1 := err.(commondomain.Error)
	assert.True(t, ok1)
	assert.Equal(t, "GenerateRenstra.InvalidTahunRenstra", dErr1.Code)
	assert.Contains(t, dErr1.Description, "cannot be generated")

	// -------------------
	// CASE 2: Fakultas unit berbeda → InvalidFakultasUnit
	// -------------------
	cmd2 := application.GenerateRenstraCommand{
		Tahun:            "2026",
		UuidRenstra:      renstraUUID.String(),
		UuidFakultasUnit: uuid.New().String(), // berbeda
	}
	_, err2 := handler.Handle(ctx, cmd2)
	assert.Error(t, err2)
	dErr2, ok2 := err2.(commondomain.Error)
	assert.True(t, ok2)
	assert.Equal(t, "GenerateRenstra.InvalidFakultasUnit", dErr2.Code)
	assert.Contains(t, dErr2.Description, "fakultas unit is invalid")

	// -------------------
	// CASE 3: Template <= 0 → InvalidTemplate
	// -------------------
	mockRenstra.On("GetByUuid", mock.Anything, renstraUUID).Return(&domainrenstra.Renstra{
		ID: 1, UUID: renstraUUID, Tahun: "2026", FakultasUnit: 1,
	}, nil)
	mockTemplateRenstra.On("GetAllByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).
		Return([]domaintemplaterenstra.TemplateRenstraDefault{{ID: 0, UUID: uuid.New(), Tugas: "auditor1"}}, nil)

	cmd3 := application.GenerateRenstraCommand{
		Tahun:            "2026",
		UuidRenstra:      renstraUUID.String(),
		UuidFakultasUnit: fakultasUUID.String(),
	}
	_, err3 := handler.Handle(ctx, cmd3)
	assert.Error(t, err3)
	dErr3, ok3 := err3.(commondomain.Error)
	assert.True(t, ok3)
	assert.Equal(t, "GenerateRenstra.InvalidTemplate", dErr3.Code)

	// -------------------
	// CASE 4: Tugas invalid → InvalidTugas
	// -------------------
	mockTemplateRenstra.On("GetAllByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).
		Return([]domaintemplaterenstra.TemplateRenstraDefault{{ID: 1, UUID: uuid.New(), Tugas: "invalid_tugas"}}, nil)

	cmd4 := application.GenerateRenstraCommand{
		Tahun:            "2026",
		UuidRenstra:      renstraUUID.String(),
		UuidFakultasUnit: fakultasUUID.String(),
	}
	_, err4 := handler.Handle(ctx, cmd4)
	assert.Error(t, err4)
	dErr4, ok4 := err4.(commondomain.Error)
	assert.True(t, ok4)
	assert.Equal(t, "GenerateRenstra.InvalidTugas", dErr4.Code)

	// -------------------
	// CASE 5: Sukses → valid GenerateRenstra dibuat
	// -------------------
	mockTemplateRenstra.On("GetAllByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).
		Return([]domaintemplaterenstra.TemplateRenstraDefault{{ID: 1, UUID: uuid.New(), Tugas: "auditor1"}}, nil)

	mockRepo.On("GetAllRenstraNilaiByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).
		Return([]domaingenerate.GenerateRenstraDefault{}, nil)
	mockRepo.On("GetAllDokumenTambahanByTahunFakUnitDefault", mock.Anything, "2026", uint(1)).
		Return([]domaingenerate.GenerateDokumenTambahanDefault{}, nil)

	cmd5 := application.GenerateRenstraCommand{
		Tahun:            "2026",
		UuidRenstra:      renstraUUID.String(),
		UuidFakultasUnit: fakultasUUID.String(),
	}
	result, err5 := handler.Handle(ctx, cmd5)
	assert.NoError(t, err5)
	assert.NotEmpty(t, result)
}

// ==================== Test Command Handler Fail Timeout ====================
func TestGenerateRenstraCommand_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	renstraUUID := uuid.New()
	fakultasUUID := uuid.New()

	mockRenstra := new(MockRenstraRepo)

	handler := &application.GenerateRenstraCommandHandler{
		RepoRenstra: mockRenstra,
	}

	cmd := application.GenerateRenstraCommand{
		Tahun:            "2026",
		UuidRenstra:      renstraUUID.String(),
		UuidFakultasUnit: fakultasUUID.String(),
	}

	result, err := handler.Handle(ctx, cmd)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
