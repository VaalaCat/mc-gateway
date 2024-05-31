package conf

import (
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/bot/playerlist"
	"github.com/Tnze/go-mc/bot/screen"
)

var (
	Client        *bot.Client
	Player        *basic.Player
	ChatHandler   *msg.Manager
	PlayerList    *playerlist.PlayerList
	ScreenManager *screen.Manager
)
