package user

import (
	"fmt"

	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NewUserDto struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

// ToModel create a model from a dto
func (dto NewUserDto) ToModel() (*User, error) {
	role, err := StrToUserRole(dto.Role)
	if err != nil {
		zap.S().Errorf("Failed to parse role = %+v, err = %+v", dto.Role, err)
		return nil, cerror.ErrUnknownRole
	}
	if dto.Uuid != "" {
		_, err := uuid.Parse(dto.Uuid)
		if err != nil {
			zap.S().Errorf("Failed to parse uuid = %s, err = %+v", dto.Uuid, err)
			return nil, cerror.ErrBadUuid
		}
	}

	return &User{
		Uuid:     uuid.New(),
		Username: dto.Username,
		Email:    dto.Email,
		Role:     role,
	}, nil
}

// FromModel returns a dto from model struct
func (NewUserDto) FromModel(m *User) NewUserDto {
	dto := NewUserDto{
		Uuid:     m.Uuid.String(),
		Username: m.Username,
		Email:    m.Email,
		Role:     string(m.Role),
	}
	return dto
}

type UserDto struct {
	Uuid        string `json:"uuid"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	PoliceToken string `json:"policeToken,omitempty"`
}

func (dto UserDto) ToModel() (*User, error) {
	parsedUuid, err := uuid.Parse(dto.Uuid)
	if err != nil {
		zap.S().Errorf("Failed to parse uuid = %s, err = %+v", dto.Uuid, err)
		return nil, cerror.ErrBadUuid
	}

	role, err := StrToUserRole(dto.Role)
	if err != nil {
		zap.S().Errorf("Failed to parse role = %+v, err = %+v", dto.Role, err)
		return nil, cerror.ErrUnknownRole
	}

	return &User{
		Uuid:     parsedUuid,
		Username: dto.Username,
		Email:    dto.Email,
		Role:     role,
	}, nil
}

// FromModel returns a dto from model struct
func (UserDto) FromModel(m *User) UserDto {
	dto := &UserDto{
		Uuid:     m.Uuid.String(),
		Username: m.Username,
		Email:    m.Email,
		Role:     fmt.Sprint(m.Role),
	}
	return *dto
}
