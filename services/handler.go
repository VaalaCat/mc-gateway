package services

import (
	"tg-mc/conf"
	"tg-mc/services/gateway"
	"tg-mc/services/mc"
	"tg-mc/services/tgbot"
)

func Run() {
	settings := conf.GetBotSettings()

	if settings.EnableGateway {
		go gateway.StartGateway()
	}

	if settings.EnableBridge {
		go mc.StartBridgeClient()
	}

	if settings.EnableBot {
		if settings.EnableBridge {
			tgbot.Run(mc.SendMsg, mc.SendCommand)
		} else {
			tgbot.Run(mc.SendMsg, func(s string) error { return nil })
		}
	}
}
