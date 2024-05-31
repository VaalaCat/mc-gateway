package tgbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func TalkHandler(m *tgbotapi.Message, i interface{}) {
	sendFunc := i.(func(string))
	logrus.Infof("id is %d", m.Chat.ID)
	msg := fmt.Sprintf("%v: %v", m.From.UserName, m.CommandArguments())
	sendFunc(msg)
}
