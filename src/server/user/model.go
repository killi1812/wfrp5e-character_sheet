package user

import (
	"errors"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRole string

const (
	ROLE_USER        UserRole = "user"
	ROLE_ADMIN       UserRole = "admin"
	ROLE_SUPER_ADMIN UserRole = "superadmin"
)

func StrToUserRole(text string) (UserRole, error) {
	role := UserRole(text)
	switch role {
	case ROLE_ADMIN:
		return ROLE_ADMIN, nil
	case ROLE_USER:
		return ROLE_USER, nil
	case ROLE_SUPER_ADMIN:
		return ROLE_SUPER_ADMIN, nil

	default:
		return "", cerror.ErrUnknownRole
	}
}

var _VALID_USER_ROLES = map[UserRole]bool{
	ROLE_USER:        true,
	ROLE_SUPER_ADMIN: true,
	ROLE_ADMIN:       true,
}

type Session struct {
	UserUuid     uuid.UUID `json:"userUuid" bson:"user_uuid"`
	RefreshToken string    `json:"refreshToken" bson:"refresh_token"`
}

type User struct {
	ID           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Uuid         uuid.UUID     `json:"uuid" bson:"uuid"`
	Username     string        `json:"username" bson:"username"`
	FirstName    string        `json:"firstName" bson:"first_name"`
	LastName     string        `json:"lastName" bson:"last_name"`
	OIB          string        `json:"oib" bson:"oib"`
	Residence    string        `json:"residence" bson:"residence"`
	BirthDate    time.Time     `json:"birthDate" bson:"birth_date"`
	Email        string        `json:"email" bson:"email"`
	PasswordHash string        `json:"passwordHash" bson:"password_hash"`
	Role         UserRole      `json:"role" bson:"role"`
	CreatedAt    time.Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updated_at"`
}

func (u *User) ValidateRole() error {
	if _, ok := _VALID_USER_ROLES[u.Role]; !ok {
		return errors.New("invalid user role")
	}
	return nil
}

func (u *User) Update(user *User) *User {
	u.BirthDate = user.BirthDate
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.OIB = user.OIB
	u.Residence = user.Residence
	u.Email = user.Email
	u.Role = user.Role
	u.UpdatedAt = time.Now()

	return u
}
