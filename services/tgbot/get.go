package tgbot

import (
	"fmt"
	"strconv"
	"strings"
	"tg-mc/conf"
	"tg-mc/models"
	"tg-mc/services/utils"
	commonUtils "tg-mc/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func GetHandler(msg *tgbotapi.Message, i interface{}) {
	if !utils.IsAdmin(msg) &&
		len(msg.CommandArguments()) != 0 {
		tm := tgbotapi.NewMessage(msg.Chat.ID, "您不是管理员，没有该权限")
		conf.Bot.Send(tm)
		return
	} else if utils.IsAdmin(msg) &&
		len(msg.CommandArguments()) != 0 {
		a := commonUtils.GetArgs(msg.CommandArguments())
		if len(a) != 2 {
			tm := tgbotapi.NewMessage(msg.Chat.ID, "参数错误，样例：\n```\n/get <tgid|username> <value>\n```")
			tm.ParseMode = "Markdown"
			conf.Bot.Send(tm)
			return
		}
		if a[0] == "tgid" {
			tgid, err := strconv.ParseInt(a[1], 10, 64)
			if err != nil {
				conf.Bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "ID错误，应该为int64"))
				return
			}
			u, err := models.GetUsersByTGID(tgid)
			if err != nil {
				tm := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("查询出错，err：\n```\n%+v\n```", err))
				tm.ParseMode = "Markdown"
				conf.Bot.Send(tm)
				return
			}
			tm := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("用户信息：\n```\n%+v\n```", u))
			tm.ParseMode = "Markdown"
			conf.Bot.Send(tm)
		}
		if a[0] == "username" {
			u, err := models.GetUserByMCName(a[1])
			if err != nil {
				tm := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("查询出错，err：\n```\n%+v\n```", err))
				tm.ParseMode = "Markdown"
				conf.Bot.Send(tm)
				return
			}
			tm := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("用户信息：\n```\n%+v\n```", u))
			tm.ParseMode = "Markdown"
			conf.Bot.Send(tm)
		}
		return
	}
	logrus.Infof("id is %d", msg.Chat.ID)
	u, err := models.GetUsersByTGID(msg.From.ID)
	if err != nil {
		m := tgbotapi.NewMessage(msg.Chat.ID, "你还没有绑定")
		conf.Bot.Send(m)
		return
	}
	ans := strings.Join(lo.Map(u, func(u models.User, _ int) string { return u.MCName }), "\n")
	m := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("你的MCID绑定有\n%v", ans))
	conf.Bot.Send(m)
	return
}
