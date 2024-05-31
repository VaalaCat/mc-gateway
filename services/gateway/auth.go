package gateway

import (
	"fmt"
	"sync"
	"tg-mc/conf"
	"tg-mc/defs"
	"tg-mc/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Auth interface {
	IsAuthed(u models.User) bool
	RequestAuth(u models.User, req *LoginRequest)
	Reject(u models.User)
	SetAuth(u models.User)
}

type Authcator struct {
	UserMap *sync.Map
}

var authcator *Authcator

func GetAuthcator() Auth {
	if authcator == nil {
		authcator = &Authcator{
			UserMap: &sync.Map{},
		}
	}
	return authcator
}

func (a *Authcator) getUserLoginReq(u models.User) *LoginRequest {
	chAny, ok := a.UserMap.Load(u.MCName)
	if !ok {
		return nil
	}

	loginReq, ok := chAny.(*LoginRequest)
	if !ok {
		return nil
	}

	return loginReq
}

func (a *Authcator) IsAuthed(u models.User) bool {
	loginReq := a.getUserLoginReq(u)
	if loginReq == nil {
		return false
	}

	return <-loginReq.Resolve
}

func (a *Authcator) RequestAuth(u models.User, req *LoginRequest) {
	a.UserMap.Store(u.MCName, req)
	m := tgbotapi.NewMessage(u.TGID, fmt.Sprintf("MC用户：%v 尝试登录，请选择操作", u.MCName))

	m.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("批准", defs.NewApproveCommand(u.MCName).ToJSON()),
			tgbotapi.NewInlineKeyboardButtonData("拒绝", defs.NewRejectCommand(u.MCName).ToJSON())),
	)
	conf.Bot.Send(m)
}

func (a *Authcator) Reject(u models.User) {
	loginReq := a.getUserLoginReq(u)
	if loginReq == nil {
		return
	}

	loginReq.Resolve <- false
}

func (a *Authcator) SetAuth(u models.User) {
	loginReq := a.getUserLoginReq(u)
	if loginReq == nil {
		return
	}

	loginReq.Resolve <- true
}
