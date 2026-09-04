package user_test

import (
	"testing"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/user"
	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"
	"github.com/killi1812/wfrp5e-character_sheet/util/format"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUserDto_ToModel(t *testing.T) {
	validDateStr := "1995-03-20"
	validTime, _ := time.Parse(format.DateFormat, validDateStr)
	userUUID := uuid.New()

	tests := []struct {
		name    string
		dto     user.NewUserDto
		want    *user.User
		wantErr error
	}{
		{
			name: "Valid DTO to Model - No UUID in DTO",
			dto: user.NewUserDto{
				FirstName: "Test",
				LastName:  "User",
				OIB:       "11223344556",
				Residence: "Test Residence",
				BirthDate: validDateStr,
				Email:     "test.user@example.com",
				Password:  "password123",
				Role:      string(user.ROLE_USER),
			},
			want: &user.User{
				FirstName: "Test",
				LastName:  "User",
				OIB:       "11223344556",
				Residence: "Test Residence",
				BirthDate: validTime,
				Email:     "test.user@example.com",
				Role:      user.ROLE_USER,
			},
			wantErr: nil,
		},
		{
			name: "Valid DTO to Model - With UUID in DTO",
			dto: user.NewUserDto{
				Uuid:      userUUID.String(),
				FirstName: "Another",
				LastName:  "User",
				OIB:       "66778899001",
				Residence: "Another Residence",
				BirthDate: validDateStr,
				Email:     "another.user@example.com",
				Password:  "securepass",
				Role:      string(user.ROLE_ADMIN),
			},
			want: &user.User{
				FirstName: "Another",
				LastName:  "User",
				OIB:       "66778899001",
				Residence: "Another Residence",
				BirthDate: validTime,
				Email:     "another.user@example.com",
				Role:      user.ROLE_ADMIN,
			},
			wantErr: nil,
		},
		{
			name: "Invalid BirthDate format",
			dto: user.NewUserDto{
				FirstName: "Bad",
				LastName:  "Date",
				BirthDate: "20-03-1995",
				Role:      "hak",
				OIB:       "12345678901",
				Email:     "bad.date@example.com",
				Password:  "password",
				Residence: "Some place",
			},
			want:    nil,
			wantErr: cerror.ErrBadDateFormat,
		},
		{
			name: "Invalid Role",
			dto: user.NewUserDto{
				FirstName: "Bad",
				LastName:  "Role",
				BirthDate: validDateStr,
				Role:      "nonexistent_role",
				OIB:       "12345678901",
				Email:     "bad.role@example.com",
				Password:  "password",
				Residence: "Some place",
			},
			want:    nil,
			wantErr: cerror.ErrUnknownRole,
		},
		{
			name: "Invalid UUID in DTO",
			dto: user.NewUserDto{
				Uuid:      "this-is-not-a-uuid",
				FirstName: "Bad",
				LastName:  "UUID",
				BirthDate: validDateStr,
				Role:      string(user.ROLE_USER),
				OIB:       "12345678901",
				Email:     "bad.uuid@example.com",
				Password:  "password",
				Residence: "Some place",
			},
			want:    nil,
			wantErr: cerror.ErrBadUuid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dto.ToModel()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.NotEqual(t, uuid.Nil, got.Uuid)
				assert.Equal(t, tt.want.FirstName, got.FirstName)
				assert.Equal(t, tt.want.LastName, got.LastName)
				assert.Equal(t, tt.want.OIB, got.OIB)
				assert.Equal(t, tt.want.Residence, got.Residence)
				assert.Equal(t, tt.want.BirthDate, got.BirthDate)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Role, got.Role)
			}
		})
	}
}

func TestNewUserDto_FromModel(t *testing.T) {
	userUUID := uuid.New()
	birthTime, _ := time.Parse(format.DateFormat, "1992-11-05")
	userModel := &user.User{
		Uuid:         userUUID,
		FirstName:    "ModelF",
		LastName:     "ModelL",
		OIB:          "55443322110",
		Residence:    "Model Residence",
		BirthDate:    birthTime,
		Email:        "model.user@example.com",
		PasswordHash: "somehash",
		Role:         user.ROLE_SUPER_ADMIN,
	}

	expectedDto := user.NewUserDto{
		Uuid:      userUUID.String(),
		FirstName: "ModelF",
		LastName:  "ModelL",
		OIB:       "55443322110",
		Residence: "Model Residence",
		BirthDate: "1992-11-05",
		Email:     "model.user@example.com",
		Password:  "",
		Role:      string(user.ROLE_SUPER_ADMIN),
	}

	gotDto := user.NewUserDto{}.FromModel(userModel)

	assert.Equal(t, expectedDto, gotDto)
}

func TestUserDto_ToModel(t *testing.T) {
	validUUID := uuid.New()
	validDateStr := "1990-01-15"
	validTime, _ := time.Parse(format.DateFormat, validDateStr)

	tests := []struct {
		name    string
		dto     user.UserDto
		want    *user.User
		wantErr error
	}{
		{
			name: "Valid DTO to Model",
			dto: user.UserDto{
				Uuid:      validUUID.String(),
				FirstName: "John",
				LastName:  "Doe",
				OIB:       "12345678901",
				Residence: "123 Main St",
				BirthDate: validDateStr,
				Email:     "john.doe@example.com",
				Role:      "user",
			},
			want: &user.User{
				Uuid:      validUUID,
				FirstName: "John",
				LastName:  "Doe",
				OIB:       "12345678901",
				Residence: "123 Main St",
				BirthDate: validTime,
				Email:     "john.doe@example.com",
				Role:      user.ROLE_USER,
			},
			wantErr: nil,
		},
		{
			name: "Invalid UUID",
			dto: user.UserDto{
				Uuid:      "not-a-uuid",
				FirstName: "Jane",
				LastName:  "Doe",
				Role:      "user",
				BirthDate: validDateStr,
			},
			want:    nil,
			wantErr: cerror.ErrBadUuid,
		},
		{
			name: "Invalid BirthDate format",
			dto: user.UserDto{
				Uuid:      validUUID.String(),
				FirstName: "Jim",
				LastName:  "Beam",
				BirthDate: "15-01-1990",
				Role:      "user",
			},
			want:    nil,
			wantErr: cerror.ErrBadDateFormat,
		},
		{
			name: "Invalid Role",
			dto: user.UserDto{
				Uuid:      validUUID.String(),
				FirstName: "Jack",
				LastName:  "Daniels",
				BirthDate: validDateStr,
				Role:      "invalid_role",
			},
			want:    nil,
			wantErr: cerror.ErrUnknownRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dto.ToModel()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Uuid, got.Uuid)
				assert.Equal(t, tt.want.FirstName, got.FirstName)
				assert.Equal(t, tt.want.LastName, got.LastName)
				assert.Equal(t, tt.want.OIB, got.OIB)
				assert.Equal(t, tt.want.Residence, got.Residence)
				assert.Equal(t, tt.want.BirthDate, got.BirthDate)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Role, got.Role)
			}
		})
	}
}

func TestUserDto_FromModel(t *testing.T) {
	userUUID := uuid.New()
	birthTime, _ := time.Parse(format.DateFormat, "1985-07-20")
	userModel := &user.User{
		Uuid:      userUUID,
		FirstName: "Alice",
		LastName:  "Smith",
		OIB:       "09876543210",
		Residence: "456 Oak Ave",
		BirthDate: birthTime,
		Email:     "alice.smith@example.com",
		Role:      user.ROLE_USER,
	}

	expectedDto := user.UserDto{
		Uuid:      userUUID.String(),
		FirstName: "Alice",
		LastName:  "Smith",
		OIB:       "09876543210",
		Residence: "456 Oak Ave",
		BirthDate: "1985-07-20",
		Email:     "alice.smith@example.com",
		Role:      string(user.ROLE_USER),
	}

	var gotDto user.UserDto
	gotDto = gotDto.FromModel(userModel)

	assert.Equal(t, expectedDto, gotDto)
}
