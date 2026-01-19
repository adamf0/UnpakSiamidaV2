package applicationtest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	app "UnpakSiamida/modules/renstra/application/GetRenstraDefault"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	infra "UnpakSiamida/modules/renstra/infrastructure"
)

func TestGetRenstraDefaultByUuid_Success(t *testing.T) {
	db, cleanup := setupRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewRenstraRepository(db)
	handler := app.GetRenstraDefaultByUuidQueryHandler{Repo: repo}

	// UUID yang pasti ada di seed data
	existingUuid := "88f3f34e-ee4e-4ca0-8bbc-ddd4b726ca06"
	q := app.GetRenstraDefaultByUuidQuery{Uuid: existingUuid}

	renstra, err := handler.Handle(context.Background(), q)
	assert.NoError(t, err)
	assert.NotNil(t, renstra)
	assert.Equal(t, existingUuid, renstra.UUID.String())
}

func TestGetRenstraDefaultByUuid_Fail_NotFound(t *testing.T) {
	db, cleanup := setupRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewRenstraRepository(db)
	handler := app.GetRenstraDefaultByUuidQueryHandler{Repo: repo}

	// UUID valid tapi tidak ada di DB
	notExistUuid := "11111111-1111-1111-1111-111111111111"
	q := app.GetRenstraDefaultByUuidQuery{Uuid: notExistUuid}

	renstra, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Nil(t, renstra)
	assert.Equal(t, domainrenstra.NotFound(notExistUuid).Error(), err.Error())
}

func TestGetRenstraDefaultByUuid_Edge_InvalidUUID(t *testing.T) {
	db, cleanup := setupRenstraMySQL(t)
	defer cleanup()

	repo := infra.NewRenstraRepository(db)
	handler := app.GetRenstraDefaultByUuidQueryHandler{Repo: repo}

	// UUID tidak valid (parsing gagal)
	invalidUuid := "not-a-uuid"
	q := app.GetRenstraDefaultByUuidQuery{Uuid: invalidUuid}

	renstra, err := handler.Handle(context.Background(), q)
	assert.Error(t, err)
	assert.Nil(t, renstra)
	assert.Equal(t, domainrenstra.NotFound(invalidUuid).Error(), err.Error())
}
