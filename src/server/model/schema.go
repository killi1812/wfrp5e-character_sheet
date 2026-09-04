package model

import "github.com/killi1812/wfrp5e-character_sheet/user"

// NOTE: Here register all models to be used in migration

// GetAllModels returns an array of all models
func GetAllModels() []any {
	return []any{
		&user.User{},
		&user.Session{},
	}
}
