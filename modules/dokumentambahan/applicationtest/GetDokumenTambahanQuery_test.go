package applicationtest

import (
	"context"
	"fmt"
	"testing"

	common "UnpakSiamida/common/domain"
	app "UnpakSiamida/modules/dokumentambahan/application/GetDokumenTambahan"
	infra "UnpakSiamida/modules/dokumentambahan/infrastructure"

	"github.com/stretchr/testify/assert"
)

func TestGetDokumenTambahanByUuid_Success(t *testing.T) {
	db, cleanup := setupDokumenTambahanMySQL(t)
	defer cleanup()

	repo := infra.NewDokumenTambahanRepository(db)
	handler := app.GetDokumenTambahanByUuidQueryHandler{Repo: repo}

	// UUID fix yang kamu tentukan
	fixedUUID := "c836800f-8c09-4e04-ba16-e0ca027ca571"

	q := app.GetDokumenTambahanByUuidQuery{Uuid: fixedUUID}

	res, err := handler.Handle(context.Background(), q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UUID.String() != fixedUUID {
		t.Fatalf("expected UUID %s, got %s", fixedUUID, res.UUID)
	}
}

func TestGetDokumenTambahanByUuid_Fail(t *testing.T) {
	db, cleanup := setupDokumenTambahanMySQL(t)
	defer cleanup()

	repo := infra.NewDokumenTambahanRepository(db)
	handler := app.GetDokumenTambahanByUuidQueryHandler{Repo: repo}

	uuid := "11111111-1111-1111-1111-111111111111"
	q := app.GetDokumenTambahanByUuidQuery{Uuid: uuid}

	_, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)

	commonErr, ok := err.(common.Error)
	assert.True(t, ok)
	assert.Equal(t, "DokumenTambahan.NotFound", commonErr.Code)
	assert.Contains(t, fmt.Sprintf("DokumenTambahan with identifier %s not found", uuid), commonErr.Description)
}
