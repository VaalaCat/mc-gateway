package mc

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"tg-mc/conf"
	su "tg-mc/services/utils"
	"tg-mc/utils"
	"time"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

func GetJoinedPlayer(m chat.Message) (userName string, err error) {
	if m.Translate != "multiplayer.player.joined" && len(m.With) == 0 {
		return "", errors.New(ErrNotJoined)
	}
	userName = m.With[0].Text
	return
}

func GetLeftPlayer(m chat.Message) (userName string, err error) {
	if m.Translate != "multiplayer.player.left" && len(m.With) == 0 {
		return "", errors.New(ErrNotJoined)
	}
	userName = m.With[0].Text
	return
}

func HandleJoinGame(userName string, mention bool, expireMode bool) {
	// return

	// u, err := models.GetUserByMCName(userName)
	// if err != nil {
	// 	logrus.Error("get user name error: ", err)
	// }

	// switch u.Status {
	// case StatusNormal:
	// if !GetAuthcator().IsAuthed(u, expireMode) {
	// 	m := tgbotapi.NewMessage(u.TGID, fmt.Sprintf("MC用户：%v 尝试登录，请手动允许，每次授权持续30秒", userName))
	// 	m.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
	// 		tgbotapi.NewInlineKeyboardRow(
	// 			tgbotapi.NewInlineKeyboardButtonData("批准", defs.NewApproveCommand(u.MCName).ToJSON()),
	// 			tgbotapi.NewInlineKeyboardButtonData("拒绝", defs.NewRejectCommand(u.MCName).ToJSON())),
	// 	)
	// 	conf.Bot.Send(m)
	// 	KickPlayer(userName)
	// 	return
	// }
	if mention {
		SendMsgToPlayer("欢迎回来！", userName)
	}
	// case StatusPending:
	// 	SendMsgToPlayer("你还没有绑定 Telegram 哦, 5秒后你将会被踢出。请在群组中发送 /bind <你的 MC 用户名> 进行绑定。", userName)
	// 	time.Sleep(5 * time.Second)
	// 	KickPlayer(userName)
	// 	m := tgbotapi.NewMessage(conf.GetBotSettings().GroupID, fmt.Sprintf("用户：%v 没有绑定，尝试登录已被T出", userName))
	// 	conf.Bot.Send(m)
	// case StatusBanned:
	// 	SendMsgToPlayer("你已被封禁，如有疑问请联系管理员。", userName)
	// 	m := tgbotapi.NewMessage(conf.GetBotSettings().GroupID, fmt.Sprintf("用户：%v 已被封禁，尝试登录已被T出", userName))
	// 	conf.Bot.Send(m)
	// default:
	// 	SendMsgToPlayer("未知错误，请联系管理员，你将被踢出", userName)
	// 	time.Sleep(3 * time.Second)
	// 	KickPlayer(userName)
	// 	m := tgbotapi.NewMessage(conf.GetBotSettings().GroupID, fmt.Sprintf("用户：%v 登录失败，错误未知，已被T出", userName))
	// 	conf.Bot.Send(m)
	// }
}

func SendCommand(cmd string) error {
	var salt int64
	if err := binary.Read(rand.Reader, binary.BigEndian, &salt); err != nil {
		return err
	}

	err := conf.Client.Conn.WritePacket(packet.Marshal(
		packetid.ServerboundChatCommand,
		packet.String(cmd),
		packet.Long(time.Now().UnixMilli()),
		packet.Long(salt),
		packet.VarInt(0), // signature
		packet.VarInt(0),
		packet.NewFixedBitSet(20),
	))
	return err
}

func KickPlayer(userName string) error {
	err := SendCommand("kick " + userName)
	return err
}

func CronKick() {
	utils.CronStart(func() {
		users := su.GetAlivePlayerList()
		for _, u := range users {
			HandleJoinGame(u, false, false)
		}
	})
}

func isBotMsg(msg chat.Message) bool {
	return msg.Translate == "commands.message.display.outgoing"
}

func HandleLeftGame(userName string) {
	// u, err := models.GetUserByMCName(userName)
	// if err != nil {
	// 	logrus.Error("get user name error: ", err)
	// }
	// GetAuthcator().Reject(u)
}
