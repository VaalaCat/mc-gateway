> 我的世界 Minecraft 非侵入式非插件认证网关

## 背景

Vaala在V服部署了很久的我的世界服务器，前写日子一直被熊玩家登录毁坏服务器设施诟病。为了提升游戏体验，同时保障服务器安全，我开发了这款非侵入式认证网关。它能够在不开启插件的情况下保证玩家登录的安全，不影响游戏本身的乐趣。

同时因为没有开启插件，服主可以轻松的升级游戏版本，无需担忧插件版本兼容性的问题，网关独立于游戏部署，也不会影响游戏的表现。

为了实现这一目标，网关采用了Golang编写，保证了高效性和稳定性。网关的部署极大提升了V服的整体安全性，受到了玩家的一致好评。

## 部署

- [https://github.com/vaala/mc-gateway](https://github.com/vaala/mc-gateway)

我们建议使用docker

```bash
docker run -d -v `PWD`/data:/data \
	-p 25565:25565 \
	-e GATEWAY_MC_SERVER_HOST=<MC服务器的地址> \
	-e GATEWAY_MC_SERVER_PORT=<MC服务器的端口> \
	-e ENABLE_BRIDGE=false \
	-e BOT_TOKEN=<TGBot的token，例如123123123:xxxxxxxxxxxxxx> \
	-e HTTP_PROXY=<HTTP代理地址，例如http://x.x.x.x:8080> \
	vaalacat/mc-gateway
```