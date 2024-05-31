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

func ApproveHandler(update tgbotapi.Update, cmd defs.Command) {
	u, err := models.GetUserByMCName(cmd.Argstr)
	if err != nil {
		return
	}
	gateway.GetAuthcator().SetAuth(u)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "已授权")
	if _, err := conf.Bot.Request(callback); err != nil {
		logrus.Panic(err)
	}
	conf.Bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID))
	conf.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("已授权☑️: %s 登录MC", u.MCName)))
}

// func ApproveHandler(update tgbotapi.Update, cmd defs.Command) {
// 	u, err := models.GetUserByTGID(update.CallbackQuery.From.ID)
// 	if err != nil {
// 		return
// 	}
// 	mc.GetAuthcator().Auth(u)
// 	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "已授权")
// 	if _, err := conf.Bot.Request(callback); err != nil {
// 		logrus.Panic(err)
// 	}
// 	conf.Bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
// 		update.CallbackQuery.Message.MessageID))
// 	conf.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
// 		fmt.Sprintf("已授权☑️: %s 登录MC", u.MCName)))
// }
