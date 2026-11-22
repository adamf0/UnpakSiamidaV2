package application

import (
	commonDomain "UnpakSiamida/common/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stubFakultasUnitRepo struct {
	items        map[uuid.UUID]*domainfakultasunit.FakultasUnit
	lastNama     *string
	lastIsActive *bool
}

func (s *stubFakultasUnitRepo) GetDefaultByUuid(_ context.Context, uid uuid.UUID) (*domainfakultasunit.FakultasUnit, error) {
	item, ok := s.items[uid]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *item
	return &cp, nil
}

func (s *stubFakultasUnitRepo) GetAll(_ context.Context, _ string, _ []commonDomain.SearchFilter, _, _ *int, _ bool) ([]domainfakultasunit.FakultasUnit, int64, error) {
	return nil, 0, nil
}

func (s *stubFakultasUnitRepo) UpdateTarget(_ context.Context, uid uuid.UUID, nama *string, isActive *bool) error {
	s.lastNama = nama
	s.lastIsActive = isActive
	item, ok := s.items[uid]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	if nama != nil {
		item.NamaRaw = *nama
		item.Nama = *nama
	}
	if isActive != nil {
		item.IsActive = *isActive
	}
	return nil
}

func (s *stubFakultasUnitRepo) SetupUuid(_ context.Context) error    { return nil }
func (s *stubFakultasUnitRepo) EnsureSchema(_ context.Context) error { return nil }

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func TestUpdateFakultasUnitCommandValidation(t *testing.T) {
	validUuid := uuid.NewString()

	if err := UpdateFakultasUnitCommandValidation(UpdateFakultasUnitCommand{
		Uuid:     validUuid,
		IsActive: boolPtr(false),
	}); err != nil {
		t.Fatalf("expected valid visibility-only command, got %v", err)
	}

	if err := UpdateFakultasUnitCommandValidation(UpdateFakultasUnitCommand{
		Uuid: validUuid,
		Nama: strPtr("BPSI UNPAK"),
	}); err != nil {
		t.Fatalf("expected valid rename-only command, got %v", err)
	}

	if err := UpdateFakultasUnitCommandValidation(UpdateFakultasUnitCommand{
		Uuid: validUuid,
	}); err == nil {
		t.Fatal("expected error when both nama and is_active are nil")
	}

	if err := UpdateFakultasUnitCommandValidation(UpdateFakultasUnitCommand{
		Uuid: validUuid,
		Nama: strPtr("   "),
	}); err == nil {
		t.Fatal("expected error when nama is blank")
	}

	if err := UpdateFakultasUnitCommandValidation(UpdateFakultasUnitCommand{
		Uuid: validUuid,
		Nama: strPtr(strings.Repeat("a", 101)),
	}); err == nil {
		t.Fatal("expected error when nama exceeds 100 chars")
	}
}

func TestUpdateFakultasUnitCommandHandler(t *testing.T) {
	uid := uuid.New()
	repo := &stubFakultasUnitRepo{
		items: map[uuid.UUID]*domainfakultasunit.FakultasUnit{
			uid: {
				ID:       1,
				UUID:     uid,
				Nama:     "ILMU HUKUM (S1)",
				NamaRaw:  "ILMU HUKUM",
				Type:     "prodi",
				Jenjang:  "s1",
				IsActive: true,
			},
		},
	}
	handler := &UpdateFakultasUnitCommandHandler{Repo: repo}

	// Hide only must not touch Nama
	updated, err := handler.Handle(context.Background(), UpdateFakultasUnitCommand{
		Uuid:     uid.String(),
		IsActive: boolPtr(false),
	})
	if err != nil {
		t.Fatalf("hide failed: %v", err)
	}
	if updated.IsActive {
		t.Fatal("expected IsActive=false")
	}
	if repo.lastNama != nil {
		t.Fatalf("expected nama not updated during hide, got %v", *repo.lastNama)
	}

	// Rename only must trim and preserve IsActive=false
	updated, err = handler.Handle(context.Background(), UpdateFakultasUnitCommand{
		Uuid: uid.String(),
		Nama: strPtr("  ILMU HUKUM BARU  "),
	})
	if err != nil {
		t.Fatalf("rename failed: %v", err)
	}
	if updated.NamaRaw != "ILMU HUKUM BARU" {
		t.Fatalf("expected trimmed nama, got %q", updated.NamaRaw)
	}
	if updated.IsActive {
		t.Fatal("expected IsActive to stay false after rename-only update")
	}

	// Unknown UUID returns NotFound
	_, err = handler.Handle(context.Background(), UpdateFakultasUnitCommand{
		Uuid:     uuid.NewString(),
		IsActive: boolPtr(true),
	})
	if err == nil {
		t.Fatal("expected error for unknown UUID")
	}
}
