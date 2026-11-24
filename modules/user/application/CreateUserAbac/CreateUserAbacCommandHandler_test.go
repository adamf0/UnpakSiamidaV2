package CreateUserAbac_test

import (
	"context"
	"testing"

	commonDomain "UnpakSiamida/common/domain"
	"UnpakSiamida/modules/user/application/CreateUserAbac"
	domainuser "UnpakSiamida/modules/user/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	lastCheckedUserID uint
	lastCheckedTahun  string
	lastCheckedFuID   *uint
	lastCheckedRole   string
	checkExistsResult bool
	createdAbac       *domainuser.UserAbac
}

func (m *mockUserRepo) GetByUuid(ctx context.Context, uid uuid.UUID) (*domainuser.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetByUuids(ctx context.Context, uids []uuid.UUID) ([]*domainuser.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetAllStrict(ctx context.Context) ([]domainuser.UserOptions, error) {
	return nil, nil
}
func (m *mockUserRepo) GetAll(ctx context.Context, search string, searchFilters []commonDomain.SearchFilter, page, limit *int) ([]domainuser.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Create(ctx context.Context, user *domainuser.User) error {
	return nil
}
func (m *mockUserRepo) Update(ctx context.Context, user *domainuser.User) error {
	return nil
}
func (m *mockUserRepo) Delete(ctx context.Context, uid uuid.UUID) error {
	return nil
}
func (m *mockUserRepo) SetupUuid(ctx context.Context) error {
	return nil
}
func (m *mockUserRepo) MergeUser(ctx context.Context, primaryID uint, secondaryIDs []uint) error {
	return nil
}
func (m *mockUserRepo) UnmergeUser(ctx context.Context, userID uint) error {
	return nil
}
func (m *mockUserRepo) WithTx(ctx context.Context, fn func(txRepo domainuser.IUserRepositoryTx) error) error {
	return nil
}
func (m *mockUserRepo) GetAllUserAbac(ctx context.Context, tahun, status, role, search, userUuid, fuUuid string) ([]domainuser.UserAbac, error) {
	return nil, nil
}
func (m *mockUserRepo) GetUserAbacByUuid(ctx context.Context, uid uuid.UUID) (*domainuser.UserAbac, error) {
	return nil, nil
}
func (m *mockUserRepo) CreateUserAbac(ctx context.Context, abac *domainuser.UserAbac) error {
	m.createdAbac = abac
	return nil
}
func (m *mockUserRepo) UpdateUserAbac(ctx context.Context, id uint, updates map[string]interface{}) error {
	return nil
}
func (m *mockUserRepo) DeleteUserAbac(ctx context.Context, uid uuid.UUID) error {
	return nil
}
func (m *mockUserRepo) CheckUserAbacExists(ctx context.Context, userID uint, tahun string, fuID *uint, role string) (bool, error) {
	m.lastCheckedUserID = userID
	m.lastCheckedTahun = tahun
	m.lastCheckedFuID = fuID
	m.lastCheckedRole = role
	return m.checkExistsResult, nil
}
func (m *mockUserRepo) GetFakultasUnitInfo(ctx context.Context, id *uint, uidStr string) (uint, string, string, error) {
	return 0, "", "", nil
}
func (m *mockUserRepo) InsertHistoryEmails(ctx context.Context, histories []domainuser.HistoryEmail) error {
	return nil
}
func (m *mockUserRepo) GetEmailStatusByOutboxID(ctx context.Context, outboxID string) ([]domainuser.HistoryEmail, error) {
	return nil, nil
}
func (m *mockUserRepo) GetEmailHistories(ctx context.Context, search, status string, page, limit int) ([]domainuser.HistoryEmail, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) GetEmailHistoryByID(ctx context.Context, id uint) (*domainuser.HistoryEmail, error) {
	return nil, nil
}
func (m *mockUserRepo) UpdateEmailHistory(ctx context.Context, id uint, updates map[string]interface{}) error {
	return nil
}

func TestCreateUserAbac_CheckRoleInDuplicateCheck(t *testing.T) {
	ctx := context.Background()
	userID := uint(543)
	tahun := "2027"

	t.Run("Passes role to CheckUserAbacExists and succeeds if not duplicate", func(t *testing.T) {
		mockRepo := &mockUserRepo{
			checkExistsResult: false,
		}
		handler := &CreateUserAbac.CreateUserAbacCommandHandler{
			Repo: mockRepo,
		}

		cmd := CreateUserAbac.CreateUserAbacCommand{
			IdUser: &userID,
			Tahun:  tahun,
			Role:   "anggaran",
			Status: "active",
		}

		res, err := handler.Handle(ctx, cmd)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "anggaran", res.Role)
		assert.Equal(t, userID, mockRepo.lastCheckedUserID)
		assert.Equal(t, tahun, mockRepo.lastCheckedTahun)
		assert.Nil(t, mockRepo.lastCheckedFuID)
		assert.Equal(t, "anggaran", mockRepo.lastCheckedRole)
	})

	t.Run("Returns AbacDuplicateAssignment when same role already exists", func(t *testing.T) {
		mockRepo := &mockUserRepo{
			checkExistsResult: true,
		}
		handler := &CreateUserAbac.CreateUserAbacCommandHandler{
			Repo: mockRepo,
		}

		cmd := CreateUserAbac.CreateUserAbacCommand{
			IdUser: &userID,
			Tahun:  tahun,
			Role:   "anggaran",
			Status: "active",
		}

		res, err := handler.Handle(ctx, cmd)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, domainuser.AbacDuplicateAssignment().Code, err.(commonDomain.Error).Code)
		assert.Equal(t, "anggaran", mockRepo.lastCheckedRole)
	})
}
