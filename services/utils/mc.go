package utils

import "tg-mc/conf"

func GetAlivePlayer() string {
	ans := "当前在线玩家：\n"
	for _, v := range conf.PlayerList.PlayerInfos {
		ans += v.Name + "\n"
	}
	return ans
}

func GetAlivePlayerList() []string {
	ans := []string{}
	for _, v := range conf.PlayerList.PlayerInfos {
		if v.Name == conf.GetBotSettings().MCBotName {
			continue
		}
		ans = append(ans, v.Name)
	}
	return ans
}
