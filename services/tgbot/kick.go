package tgbot

import (
	"fmt"
	"tg-mc/conf"
	"tg-mc/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func KickHandler(m *tgbotapi.Message, i interface{}) {
	f, ok := i.(func(string) error)
	if !ok {
		return
	}

	u, err := models.GetUserByTGID(m.From.ID)
	if err != nil {
		conf.Bot.Send(tgbotapi.NewMessage(m.Chat.ID, "您还没有绑定账号，请先绑定"))
		return
	}

	err = f(fmt.Sprintf("kick %s", u.MCName))
	if err != nil {
		conf.Bot.Send(tgbotapi.NewMessage(m.Chat.ID, err.Error()))
		return
	}

	conf.Bot.Send(tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("已踢出用户 %s", u.MCName)))
}
