package bootstrap

import (
	"backend/internal/telegram"

	"go.uber.org/fx"
)

var TelegramModule = fx.Module(
	"telegram",
	fx.Provide(telegram.NewTelegramBot),
)
