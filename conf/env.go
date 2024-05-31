package conf

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type botSettings struct {
	HTTPProxy       string          `env:"HTTP_PROXY"`
	BotToken        string          `env:"BOT_TOKEN"`
	MCServer        string          `env:"MC_SERVER"`
	MCBotName       string          `env:"MC_BOT_NAME"`
	GroupID         int64           `env:"GROUP_ID"`
	DBPath          string          `env:"DB_PATH" env-default:"tg-mc.db"`
	BotAPI          string          `env:"TG_BOT_API"`
	AdminID         []int64         `env:"ADMIN_ID"`
	GatewaySettings GatewaySettings `env-prefix:"GATEWAY_"`
	EnableGateway   bool            `env:"ENABLE_GATEWAY" env-default:"true"`
	EnableBridge    bool            `env:"ENABLE_BRIDGE" env-default:"true"`
	EnableBot       bool            `env:"ENABLE_BOT" env-default:"true"`
}

type GatewaySettings struct {
	ServerHost string `json:"server_host" env:"MC_SERVER_HOST" env-default:"127.0.0.1"`
	ServerPort int    `json:"server_port" env:"MC_SERVER_PORT" env-default:"25566"`
	ProxyHost  string `json:"proxy_host" env:"PROXY_HOST" env-default:"127.0.0.1"`
	ProxyPort  int    `json:"proxy_port" env:"PROXY_PORT" env-default:"25565"`
}

var (
	botSettingsInstance botSettings
)

func init() {
	godotenv.Load()
	cleanenv.ReadEnv(&botSettingsInstance)
	log.Printf("Bot settings: %+v", botSettingsInstance)
}

func GetBotSettings() *botSettings {
	return &botSettingsInstance
}
