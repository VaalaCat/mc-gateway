package utils

import (
	"tg-mc/conf"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func SendMsgToGroup(msg string) error {
	msgT := tgbotapi.NewMessage(conf.GetBotSettings().GroupID, msg)
	_, err := conf.Bot.Send(msgT)
	return err
}

func SendMsg(chatID int64, msg string) error {
	msgT := tgbotapi.NewMessage(chatID, msg)
	_, err := conf.Bot.Send(msgT)
	return err
}

func IsAdmin(m *tgbotapi.Message) bool {
	return lo.Contains(conf.GetBotSettings().AdminID, m.From.ID) ||
		lo.Contains(conf.GetBotSettings().AdminID, m.Chat.ID)
}
