package seed

import (
	"errors"
	"os"

	"github.com/killi1812/wfrp5e-character_sheet/user"

	"go.uber.org/zap"
)

const (
	_PASSWORD_ENV = "SUPERADMIN_PASSWORD"
)

var suadmin *user.User

// createSuperAdmin creates a SuperAdmin user if one doesn't already exist.
func createSuperAdmin() error {
	userCrud := user.NewUserCrudService()

	// Check if SuperAdmin exists
	{
		_, err := userCrud.GetUserByEmail("superadmin@test.hr")
		if err != nil {
			if errors.Is(err, user.ErrRecordNotFound) {
				zap.S().Infof("SuperAdmin not found, err %+v", err)
			} else {
				return err
			}
		} else {
			zap.S().Infoln("SuperAdmin found")
			zap.S().Infoln("Skipping superadmin creation")
			return nil
		}
	}

	zap.S().Infoln("Creating superadmin")
	password := os.Getenv(_PASSWORD_ENV)
	if password == "" {
		return errors.New("env variable is empty")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	dto := user.NewUserDto{
		Username: "superadmin",
		Email:    "superadmin@test.hr",
		Password: password,
		Role:     "admin",
	}
	newUser, err := dto.ToModel()
	if err != nil {
		return err
	}

	usr, err := userCrud.Create(newUser, dto.Password)
	if err != nil {
		return err
	}
	suadmin = usr
	zap.S().Infof("superadmin created, %+v\n", usr)
	return nil
}
