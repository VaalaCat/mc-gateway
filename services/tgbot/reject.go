package tgbot

import (
	"fmt"
	"tg-mc/conf"
	"tg-mc/defs"
	"tg-mc/models"
	"tg-mc/services/gateway"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func RejectHandler(update tgbotapi.Update, cmd defs.Command) {
	u, err := models.GetUserByMCName(cmd.Argstr)
	if err != nil {
		return
	}
	gateway.GetAuthcator().Reject(u)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "已拒绝")
	if _, err := conf.Bot.Request(callback); err != nil {
		logrus.Panic(err)
	}
	conf.Bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	conf.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("已拒绝❌: %s 登录MC", u.MCName)))
}
