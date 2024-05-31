package tgbot

import (
	"tg-mc/conf"
	"tg-mc/services/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func ListHandler(m *tgbotapi.Message, i interface{}) {
	logrus.Infof("id is %d", m.Chat.ID)
	msg := tgbotapi.NewMessage(m.Chat.ID, utils.GetAlivePlayer())
	conf.Bot.Send(msg)
}
