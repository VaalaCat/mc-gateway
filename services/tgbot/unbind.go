package tgbot

import (
	"tg-mc/conf"
	"tg-mc/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func UnbindHandler(msg *tgbotapi.Message, i interface{}) {
	logrus.Infof("id is %d", msg.Chat.ID)
	if len(msg.CommandArguments()) == 0 {
		m := tgbotapi.NewMessage(msg.Chat.ID, "请输入正确的参数")
		conf.Bot.Send(m)
		return
	}

	us, err := models.GetUsersByTGID(msg.From.ID)
	if err != nil {
		m := tgbotapi.NewMessage(msg.Chat.ID, "你还没有绑定")
		conf.Bot.Send(m)
		return
	}

	lo.Map(us, func(u models.User, _ int) bool {
		if u.MCName == msg.CommandArguments() {
			u.Status = 0
			err = u.Delete(msg.From.ID)
			if err != nil {
				m := tgbotapi.NewMessage(msg.Chat.ID, "解绑失败")
				conf.Bot.Send(m)
				return false
			}
			m := tgbotapi.NewMessage(msg.Chat.ID, "解绑成功")
			conf.Bot.Send(m)
			return false
		}
		return true
	})
}
