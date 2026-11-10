package application

import (
	"context"
	"testing"

	common "UnpakSiamida/common/domain"
	domainaccount "UnpakSiamida/modules/account/domain"

	"gorm.io/gorm"
)

type fakeAccountRepo struct {
	account         *domainaccount.Account
	takenUsernames  map[string]bool
	lastUpdatedUser string
	lastUsername    string
	lastName        string
	lastEmail       *string
	lastPassword    *string
}

func (f *fakeAccountRepo) Auth(ctx context.Context, username string, password string) (*domainaccount.Account, error) {
	return f.account, nil
}

func (f *fakeAccountRepo) Get(ctx context.Context, uuid string) (*domainaccount.Account, error) {
	if f.account == nil || f.account.UUID != uuid {
		return nil, gorm.ErrRecordNotFound
	}
	copyAcc := *f.account
	return &copyAcc, nil
}

func (f *fakeAccountRepo) IsUsernameTakenByOther(ctx context.Context, username string, excludeUserUUID string) (bool, error) {
	return f.takenUsernames[username], nil
}

func (f *fakeAccountRepo) UpdateProfile(ctx context.Context, userUUID string, username string, name string, email *string, password *string) error {
	if f.account == nil || f.account.UUID != userUUID {
		return gorm.ErrRecordNotFound
	}
	f.lastUpdatedUser = userUUID
	f.lastUsername = username
	f.lastName = name
	f.lastEmail = email
	f.lastPassword = password

	f.account.NidnUsername = username
	f.account.Name = name
	if email != nil {
		f.account.Email = *email
	} else {
		f.account.Email = ""
	}
	if password != nil {
		f.account.Password = *password
	}
	return nil
}

const testProfileUUID = "56ce6c95-e23f-463b-bcf6-80fa4bea2a1e"

func TestUpdateProfileCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		cmd     UpdateProfileCommand
		wantErr bool
	}{
		{
			name: "valid update without password or email",
			cmd: UpdateProfileCommand{
				SID:      testProfileUUID,
				Username: "admin",
				Name:     "Administrator",
				Email:    "",
			},
			wantErr: false,
		},
		{
			name: "valid update with email and password change",
			cmd: UpdateProfileCommand{
				SID:             testProfileUUID,
				Username:        "user_prodi",
				Name:            "User Prodi",
				Email:           "prodi@gmail.com",
				CurrentPassword: "123",
				NewPassword:     "secret123",
				ConfirmPassword: "secret123",
			},
			wantErr: false,
		},
		{
			name: "rejects invalid UUID",
			cmd: UpdateProfileCommand{
				SID:      "invalid-uuid",
				Username: "admin",
				Name:     "Admin",
			},
			wantErr: true,
		},
		{
			name: "rejects short username",
			cmd: UpdateProfileCommand{
				SID:      testProfileUUID,
				Username: "ab",
				Name:     "Admin",
			},
			wantErr: true,
		},
		{
			name: "rejects invalid email",
			cmd: UpdateProfileCommand{
				SID:      testProfileUUID,
				Username: "admin",
				Name:     "Admin",
				Email:    "not-an-email",
			},
			wantErr: true,
		},
		{
			name: "rejects password change missing current password",
			cmd: UpdateProfileCommand{
				SID:             testProfileUUID,
				Username:        "admin",
				Name:            "Admin",
				NewPassword:     "newpass",
				ConfirmPassword: "newpass",
			},
			wantErr: true,
		},
		{
			name: "rejects mismatched confirm password",
			cmd: UpdateProfileCommand{
				SID:             testProfileUUID,
				Username:        "admin",
				Name:            "Admin",
				CurrentPassword: "123",
				NewPassword:     "newpass1",
				ConfirmPassword: "newpass2",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := UpdateProfileCommandValidation(tc.cmd)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestUpdateProfileCommandHandler(t *testing.T) {
	t.Run("updates profile and clears email when empty", func(t *testing.T) {
		repo := &fakeAccountRepo{
			account: &domainaccount.Account{
				ID:           2,
				UUID:         testProfileUUID,
				NidnUsername: "admin",
				Password:     "123",
				Level:        "admin",
				Name:         "Admin",
				Email:        "old@unpak.ac.id",
			},
		}
		h := &UpdateProfileCommandHandler{Repo: repo}

		res, err := h.Handle(context.Background(), UpdateProfileCommand{
			SID:      testProfileUUID,
			Username: "admin",
			Name:     "Admin Baru",
			Email:    "",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Name != "Admin Baru" || res.Email != "" {
			t.Fatalf("unexpected updated account: %+v", res)
		}
		if repo.lastPassword != nil {
			t.Fatalf("password should not change when empty")
		}
	})

	t.Run("rejects wrong current password", func(t *testing.T) {
		repo := &fakeAccountRepo{
			account: &domainaccount.Account{
				ID:           2,
				UUID:         testProfileUUID,
				NidnUsername: "admin",
				Password:     "correct123",
				Name:         "Admin",
			},
		}
		h := &UpdateProfileCommandHandler{Repo: repo}

		_, err := h.Handle(context.Background(), UpdateProfileCommand{
			SID:             testProfileUUID,
			Username:        "admin",
			Name:            "Admin",
			CurrentPassword: "wrong",
			NewPassword:     "new123",
			ConfirmPassword: "new123",
		})
		if err == nil {
			t.Fatalf("expected error for wrong current password")
		}
		cerr, ok := err.(common.Error)
		if !ok || cerr.Code != "Account.InvalidCurrentPassword" || cerr.Type != common.Validation {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	t.Run("rejects duplicate username taken by another user", func(t *testing.T) {
		repo := &fakeAccountRepo{
			account: &domainaccount.Account{
				ID:           2,
				UUID:         testProfileUUID,
				NidnUsername: "admin",
				Password:     "123",
				Name:         "Admin",
			},
			takenUsernames: map[string]bool{"other_user": true},
		}
		h := &UpdateProfileCommandHandler{Repo: repo}

		_, err := h.Handle(context.Background(), UpdateProfileCommand{
			SID:      testProfileUUID,
			Username: "other_user",
			Name:     "Admin",
		})
		if err == nil {
			t.Fatalf("expected conflict error")
		}
		cerr, ok := err.(common.Error)
		if !ok || cerr.Code != "Account.UsernameDuplicate" || cerr.Type != common.Conflict {
			t.Fatalf("unexpected error: %#v", err)
		}
	})

	t.Run("changes password when current password matches", func(t *testing.T) {
		repo := &fakeAccountRepo{
			account: &domainaccount.Account{
				ID:           2,
				UUID:         testProfileUUID,
				NidnUsername: "admin",
				Password:     "123",
				Name:         "Admin",
			},
		}
		h := &UpdateProfileCommandHandler{Repo: repo}

		res, err := h.Handle(context.Background(), UpdateProfileCommand{
			SID:             testProfileUUID,
			Username:        "admin",
			Name:            "Admin",
			Email:           "admin@unpak.ac.id",
			CurrentPassword: "123",
			NewPassword:     "newpass456",
			ConfirmPassword: "newpass456",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Password != "newpass456" || res.Email != "admin@unpak.ac.id" {
			t.Fatalf("password or email did not update: %+v", res)
		}
	})
}
