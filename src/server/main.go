package main

import (
	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/auth"
	"github.com/killi1812/wfrp5e-character_sheet/charactersheet"
	"github.com/killi1812/wfrp5e-character_sheet/info"
	"github.com/killi1812/wfrp5e-character_sheet/service"
	"github.com/killi1812/wfrp5e-character_sheet/user"
	"github.com/killi1812/wfrp5e-character_sheet/util/seed"

	"go.uber.org/zap"
)

//	@securitydefinitions.bearerauth	BearerAuth

func init() {
	app.Setup()
}

func main() {
	// Provide logger
	app.Provide(zap.S)

	app.Provide(service.NewDiscordService)
	app.Provide(user.NewUserCrudService)
	app.Provide(auth.NewAuthService)
	app.Provide(charactersheet.NewCharacterSheetCrudService)

	app.RegisterApi(info.NewInfoCnt)
	app.RegisterApi(user.NewUserCtn)
	app.RegisterApi(auth.NewAuthCtn)
	app.RegisterApi(charactersheet.NewCharacterSheetCtn)

	seed.Insert()

	app.Start()
}
