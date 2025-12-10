package domaintest

import (
	"testing"

	domain "UnpakSiamida/modules/user/domain"
	
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Test NewUser - Full Scenarios
// -----------------------------------------------------------------------------
func TestNewUser(t *testing.T) {
	validEmail := "user@unpak.ac.id"
	invalidEmail := "invalid-email"

	tests := []struct {
		name         string
		username     string
		password     string
		nameField    string
		email        string
		level        string
		fakultasUnit *int
		wantSuccess  bool
		wantErrCode  string
	}{
		{
			name:        "ValidUser",
			username:    "user1",
			password:    "secret",
			nameField:   "John Doe",
			email:       validEmail,
			level:       "admin",
			fakultasUnit: nil,
			wantSuccess: true,
		},
		{
			name:        "InvalidEmail",
			username:    "user2",
			password:    "secret",
			nameField:   "Jane Doe",
			email:       invalidEmail,
			level:       "user",
			fakultasUnit: nil,
			wantSuccess: false,
			wantErrCode: domain.InvalidEmail().Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.NewUser(tt.username, tt.password, tt.nameField, tt.email, tt.level, tt.fakultasUnit)

			if tt.wantSuccess {
				require.True(t, res.IsSuccess)
				assert.Equal(t, tt.username, res.Value.Username)
				assert.Equal(t, tt.nameField, res.Value.Name)
				assert.Equal(t, tt.email, res.Value.Email)
				assert.Equal(t, tt.level, res.Value.Level)
			} else {
				require.False(t, res.IsSuccess)
				assert.Equal(t, tt.wantErrCode, res.Error.Code)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// Test UpdateUser - Full Scenarios
// -----------------------------------------------------------------------------
func TestUpdateUser(t *testing.T) {
	baseUUID := uuid.New()
	validEmail := "update@unpak.ac.id"
	invalidEmail := "invalid-email"

	prev := &domain.User{
		UUID:     baseUUID,
		Username: "existing",
		Password: "oldpass",
		Name:     "Old Name",
		Email:    "old@unpak.ac.id",
		Level:    "user",
	}

	newName := "New Name"
	newPassword := "newpass"
	newLevel := "admin"
	fakultasUnit := 10

	tests := []struct {
		name         string
		prev         *domain.User
		uid          uuid.UUID
		username     string
		password     *string
		nameField    string
		email        string
		level        string
		fakultasUnit *int
		wantSuccess  bool
		wantErrCode  string
	}{
		{
			name:        "PrevNil",
			prev:        nil,
			uid:         baseUUID,
			username:    "u1",
			password:    nil,
			nameField:   "X",
			email:       validEmail,
			level:       "user",
			wantSuccess: false,
			wantErrCode: domain.EmptyData().Code,
		},
		{
			name:        "UUIDMismatch",
			prev:        prev,
			uid:         uuid.New(),
			username:    "u1",
			password:    nil,
			nameField:   "X",
			email:       validEmail,
			level:       "user",
			wantSuccess: false,
			wantErrCode: domain.InvalidData().Code,
		},
		{
			name:        "InvalidEmail",
			prev:        prev,
			uid:         baseUUID,
			username:    "u1",
			password:    nil,
			nameField:   "X",
			email:       invalidEmail,
			level:       "user",
			wantSuccess: false,
			wantErrCode: domain.InvalidEmail().Code,
		},
		{
			name:        "SuccessWithPassword",
			prev:        prev,
			uid:         baseUUID,
			username:    "u1",
			password:    &newPassword,
			nameField:   newName,
			email:       validEmail,
			level:       newLevel,
			fakultasUnit: &fakultasUnit,
			wantSuccess: true,
		},
		{
			name:        "SuccessWithoutPassword",
			prev:        prev,
			uid:         baseUUID,
			username:    "u1",
			password:    nil,
			nameField:   newName,
			email:       validEmail,
			level:       newLevel,
			fakultasUnit: &fakultasUnit,
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.UpdateUser(tt.prev, tt.uid, tt.username, tt.password, tt.nameField, tt.email, tt.level, tt.fakultasUnit)

			if tt.wantSuccess {
				require.True(t, res.IsSuccess)
				assert.Equal(t, tt.nameField, res.Value.Name)
				assert.Equal(t, tt.email, res.Value.Email)
				assert.Equal(t, tt.level, res.Value.Level)
				if tt.password != nil {
					assert.Equal(t, *tt.password, res.Value.Password)
				}
				if tt.fakultasUnit != nil {
					assert.Equal(t, tt.fakultasUnit, res.Value.FakultasUnit)
				}
			} else {
				require.False(t, res.IsSuccess)
				assert.Equal(t, tt.wantErrCode, res.Error.Code)
			}
		})
	}
}
