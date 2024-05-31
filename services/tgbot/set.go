package tgbot

import (
	"fmt"
	"math/rand"
	"strconv"
	"tg-mc/conf"
	"tg-mc/models"
	"tg-mc/services/utils"
	commonUtils "tg-mc/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func SetHandler(m *tgbotapi.Message, i interface{}) {
	var a []string
	if t, ok := i.([]string); ok {
		a = t
	} else {
		a = commonUtils.GetArgs(m.CommandArguments())
	}

	if !utils.IsAdmin(m) {
		tm := tgbotapi.NewMessage(m.Chat.ID, "您不是管理员，没有该权限")
		conf.Bot.Send(tm)
		return
	}
	if len(a) != 3 {
		tm := tgbotapi.NewMessage(m.Chat.ID, "参数错误，样例：\n```\n/set username tgid status\n```")
		tm.ParseMode = "Markdown"
		conf.Bot.Send(tm)
		return
	}
	tgid, err := strconv.ParseInt(a[1], 10, 64)
	if err != nil {
		conf.Bot.Send(tgbotapi.NewMessage(m.Chat.ID, "ID错误，应该为int64"))
		return
	}
	status, err := strconv.ParseInt(a[2], 10, 64)
	if err != nil || !lo.Contains([]int64{0, 1, 2}, status) {
		conf.Bot.Send(tgbotapi.NewMessage(m.Chat.ID, "Status错误，应该为0(Pending),1(Normal),2(Banned)"))
		return
	}
	if err := models.CreateUser(&models.User{
		TGID:   tgid,
		MCName: a[0],
		Status: int(status),
	}); err != nil {
		tm := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("创建用户错误，err：\n```\n%+v\n```", err))
		tm.ParseMode = "Markdown"
		conf.Bot.Send(tm)
		return
	}
	conf.Bot.Send(tgbotapi.NewMessage(m.Chat.ID, "设置用户成功"))
}

func BanHandler(m *tgbotapi.Message, i interface{}) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(11)
	for {
		_, err := models.GetUserByTGID(int64(randomNumber))
		if err != nil {
			break
		} else {
			randomNumber = rand.Intn(11)
		}
	}
	SetHandler(m, []string{m.CommandArguments(), fmt.Sprintf("-%d", randomNumber), "2"})
}
