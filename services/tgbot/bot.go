package tgbot

import (
	"log"
	"net/http"
	"net/url"
	"tg-mc/conf"
	"tg-mc/defs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

var funcHandlers = map[string]func(*tgbotapi.Message, interface{}){
	"talk":   TalkHandler,
	"list":   ListHandler,
	"bind":   BindHandler,
	"unbind": UnbindHandler,
	"get":    GetHandler,
	"set":    SetHandler,
	"kick":   KickHandler,
	"ban":    BanHandler,
}

var callBackHandlers = map[string]func(tgbotapi.Update, defs.Command){
	defs.CMD_APPROVE: ApproveHandler,
	defs.CMD_REJECT:  RejectHandler,
}

func init() {
	var err error

	api := conf.GetBotSettings().BotAPI
	if len(api) == 0 {
		api = tgbotapi.APIEndpoint
	}

	HttpProxy := conf.GetBotSettings().HTTPProxy
	proxyUrl, err := url.Parse(HttpProxy)
	if err != nil {
		log.Panic(err, "HTTP_PROXY environment variable is not set correctly")
	}

	if len(HttpProxy) != 0 {
		client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		conf.Bot, err = tgbotapi.NewBotAPIWithClient(
			conf.GetBotSettings().BotToken,
			api,
			client)
	} else {
		conf.Bot, err = tgbotapi.NewBotAPIWithAPIEndpoint(conf.GetBotSettings().BotToken, api)
	}

	if err != nil {
		log.Panic(err)
	}
	conf.Bot.Debug = false
	log.Printf("Authorized on account %s", conf.Bot.Self.UserName)
}

func Run(sendFunc func(string), cmdFunc func(string) error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := conf.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			go func(update tgbotapi.Update) {
				logrus.Infof("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)
				cmd, err := defs.NewCommandFromJSON(update.CallbackQuery.Data)
				if err != nil {
					logrus.Error(err)
					return
				}
				if handler, ok := callBackHandlers[cmd.Command]; ok {
					handler(update, *cmd)
				}
			}(update)
		} else if update.Message != nil {
			go func(m *tgbotapi.Message) {
				logrus.Infof("[%s] %s", m.From.UserName, m.Text)
				if handler, ok := funcHandlers[m.Command()]; ok {
					handler(m, sendFunc)
				}
			}(update.Message)
		}
	}
}
