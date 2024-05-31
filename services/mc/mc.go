package mc

import (
	"fmt"
	"log"
	"strings"
	"tg-mc/conf"
	"tg-mc/services/utils"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/bot/playerlist"
	"github.com/Tnze/go-mc/bot/screen"
	"github.com/Tnze/go-mc/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func StartBridgeClient() {
	for {
		if err := Run(); err != nil {
			utils.SendMsgToGroup("致命错误：" + err.Error())
		}
		time.Sleep(time.Second * 5)
	}
}

func Run() error {
	conf.Client = bot.NewClient()
	conf.Client.Auth.Name = conf.GetBotSettings().MCBotName
	client := conf.Client

	conf.Player = basic.NewPlayer(client, basic.DefaultSettings, basic.EventsListener{
		GameStart:    onGameStart,
		Disconnect:   onDisconnect,
		HealthChange: onHealthChange,
		Death:        onDeath,
	})
	conf.PlayerList = playerlist.New(client)
	conf.ChatHandler = msg.New(client, conf.Player, conf.PlayerList, msg.EventsHandler{
		SystemChat:        onSystemMsg,
		PlayerChatMessage: onPlayerMsg,
		DisguisedChat:     onDisguisedMsg,
	})
	conf.ScreenManager = screen.NewManager(client, screen.EventsListener{
		Open:  nil,
		Close: nil,
	})

	err := client.JoinServer(
		conf.GetBotSettings().MCServer,
	)

	if err != nil {
		log.Printf("joinserver error: %v", err)
		return err
	}

	log.Println("Login success")

	return client.HandleGame()
}

func SendMsg(msg string) {
	go func() {
		err := conf.ChatHandler.SendMessage(msg)
		if err != nil {
			logrus.Error("send msg error: ", err)
		}
	}()
}

func SendMsgToPlayer(msg string, playerName string) {
	go func() {
		err := SendCommand(fmt.Sprintf("tell %s %s", playerName, msg))
		if err != nil {
			logrus.Error("send msg to player error: ", err)
		}
	}()
}

func onSystemMsg(msg chat.Message, overlay bool) error {
	go func() {
		log.Printf("System: %v", msg.String())
		switch msg.Translate {
		case EventPlayerJoined:
			userName, err := GetJoinedPlayer(msg)
			if err != nil {
				logrus.Error("user join error ", err)
				break
			}
			go HandleJoinGame(userName, true, true)
		case EventPlayerLeft:
			userName, err := GetLeftPlayer(msg)
			if err != nil {
				logrus.Error("user left error ", err)
				break
			}
			go HandleLeftGame(userName)
		default:
			break
		}
		// m := tgbotapi.NewMessage(conf.GetBotSettings().GroupID, fmt.Sprintf("%v", msg))
		// conf.Bot.Send(m)
	}()
	return nil
}

func onPlayerMsg(msg chat.Message, validated bool) error {
	go func() {
		if msg.Translate == "commands.message.display.outgoing" {
			return
		}
		log.Printf("Player: %s", msg)
		s := strings.Split(msg.String(), " ")
		if len(s) > 1 && !isBotMsg(msg) {
			if s[0] != fmt.Sprintf("<%v>", conf.GetBotSettings().MCBotName) {
				m := tgbotapi.NewMessage(conf.GetBotSettings().GroupID, fmt.Sprintf("%v", msg))
				_, err := conf.Bot.Send(m)
				if err != nil {
					logrus.Error(err)
				}
			}
		}
	}()
	return nil
}

func onDisguisedMsg(msg chat.Message) error {
	go func() {
		log.Printf("Disguised: %v", msg)
		m := tgbotapi.NewMessage(conf.GetBotSettings().GroupID, fmt.Sprintf("%v", msg))
		conf.Bot.Send(m)
	}()
	return nil
}

func onDeath() error {
	log.Println("Died and Respawned")
	// If we exclude Respawn(...) then the player won't press the "Respawn" button upon death
	go func() {
		time.Sleep(time.Second * 5)
		err := conf.Player.Respawn()
		if err != nil {
			log.Print(err)
		}
	}()
	return nil
}

func onGameStart() error {
	go func() {
		log.Println("Game start")
		// SendMsgToPlayer("Hello", "test")
		go CronKick()
	}()
	return nil // if err isn't nil, HandleGame() will return it.
}

func onHealthChange(health float32, foodLevel int32, foodSaturation float32) error {
	log.Printf("Health: %.2f, FoodLevel: %d, FoodSaturation: %.2f", health, foodLevel, foodSaturation)
	return nil
}

type DisconnectErr struct {
	Reason chat.Message
}

func (d DisconnectErr) Error() string {
	return "disconnect: " + d.Reason.String()
}

func onDisconnect(reason chat.Message) error {
	// return an error value so that we can stop main loop
	logrus.Error("Disconnected: ", reason)
	return DisconnectErr{Reason: reason}
}
