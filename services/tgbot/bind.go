package tgbot

import (
	"fmt"
	"tg-mc/conf"
	"tg-mc/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func BindHandler(m *tgbotapi.Message, i interface{}) {
	logrus.Infof("id is %d", m.Chat.ID)
	err := models.CreateUser(&models.User{
		TGID:   m.From.ID,
		MCName: m.CommandArguments(),
		Status: 1,
	})
	if err != nil {
		m := tgbotapi.NewMessage(m.Chat.ID, "绑定失败, err: "+err.Error())
		conf.Bot.Send(m)
		return
	}

	msg := tgbotapi.NewMessage(m.Chat.ID,
		fmt.Sprintf("绑定成功，你的MCID是%v", m.CommandArguments()))
	conf.Bot.Send(msg)
	return
}
