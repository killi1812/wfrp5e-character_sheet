package user_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/killi1812/wfrp5e-character_sheet/user"
	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"
	"github.com/stretchr/testify/assert"
)

func TestStrToUserRole(t *testing.T) {
	tests := []struct {
		input   string
		want    user.UserRole
		wantErr error
	}{
		{"admin", user.ROLE_ADMIN, nil},
		{"superadmin", user.ROLE_ADMIN, nil},
		{"user", user.ROLE_USER, nil},
		{"guest", "", cerror.ErrUnknownRole},
		{"", "", cerror.ErrUnknownRole},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := user.StrToUserRole(tt.input)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUser_ValidateRole(t *testing.T) {
	t.Run("Valid roles", func(t *testing.T) {
		uUser := user.User{Role: user.ROLE_USER}
		assert.NoError(t, uUser.ValidateRole())

		uAdmin := user.User{Role: user.ROLE_ADMIN}
		assert.NoError(t, uAdmin.ValidateRole())
	})

	t.Run("Invalid role", func(t *testing.T) {
		uBad := user.User{Role: "invalid_role"}
		assert.Error(t, uBad.ValidateRole())
	})
}

func TestUser_Update(t *testing.T) {
	oldTime := time.Now().Add(-1 * time.Hour)
	u := &user.User{
		Uuid:         uuid.New(),
		Username:     "oldname",
		Email:        "old@example.com",
		Role:         user.ROLE_USER,
		PasswordHash: "secret",
		UpdatedAt:    oldTime,
	}

	newData := &user.User{
		Username: "newname",
		Email:    "new@example.com",
		Role:     user.ROLE_ADMIN,
	}

	res := u.Update(newData)
	assert.Equal(t, "newname", res.Username)
	assert.Equal(t, "new@example.com", res.Email)
	assert.Equal(t, user.ROLE_ADMIN, res.Role)
	assert.Equal(t, "secret", res.PasswordHash) // preserved
	assert.True(t, res.UpdatedAt.After(oldTime))
}
