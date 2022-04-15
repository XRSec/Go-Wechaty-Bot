# Go-Wechaty-Bot XP PROTOCOL

> 仅供学习使用，*请勿用于非法用途*！

## Info

1. This is Use puppet-xp ,Not web4u
2. Glossary
   ```ini
   Gateway : puppet-xp
   Server: go-wechaty
   ```

## ⇲ Use

### Init

1. Clone Repo
   ```bash
   git clone https://github.com/XRSec/Go-wechaty-Bot.git wechatbot
   cd wechatbot\Gateway
   ```
   
2. Generate Token

   ```bash
   # Generate Token
   WECHATY_TOKEN：curl -s https://www.uuidgenerator.net/api/version4
   WECHATY_PUPPET_SERVICE_TOKEN："insecure_" + WECHATY_TOKEN
   # WECHATY_TOKEN WECHATY_PUPPET_SERVICE_TOKEN 可同可不同
   ```

3. Modifying a Configuration File
   ```bash
   # wechatyGateway.bat
   @set WECHATY_TOKEN=5f3029c0-0f46-4436-bdc6-02efcbad3309
   @set WECHATY_PUPPET_SERVICE_TOKEN=insecure_34bf8353-0874-4b29-851d-e8a2502fc747
   @set WECHATY_PUPPET_SERVER_PORT=25000
   ```
4. Install the Software ([**WeChat.exe Check Download**](https://github.com/wechaty/wechaty-puppet-xp/releases/download/v0.5/WeChatSetup-v3.3.0.115.exe))
   ```bash
   > node-v16.exe WeChatSetup-v3.3.0.115.exe
   > @cnpm install -g windows-build-tools
   > @cnpm install -g wechaty wechaty-puppet-xp
   ```
5. Optional operation
   ```bash
   # Set Environment
   @chdir
   # Google: How to set the path and environment variables in Windows
   ```


### Install

1. Edit `Server/config.yaml`
   ```yaml
   bot:
     adminid: wxid_n4t6b # 这个从 终端读取,你用管理员发一条消息,就会显示你的ID,或者 fmt.Println(message.From().ID())
     chat: false # 开启聊天功能
     robots: false # 开启机器人功能
   chat: # 群聊 聊天功能管理
     27714710426@chatroom: # 群聊的ID
       chat: "off" # 关闭聊天回复功能
       tuling: "off" # 关闭图灵机器人功能
   ding:
     keyword: Wechaty # Dingding 的关键词
     token: e2e18da4d2deaed25edad74bc4c91c96f48a9aa3edf937bda6a76c6a1305177c # Dingding 的 Token
     url: https://oapi.dingtalk.com/robot/send?access_token= # Dingding 的 URL
   faild: "[\e[01;31m✗\e[0m] " # 这是系统内置的输出颜色的
   info: "[\e[01;33m➜\e[0m] " # 这是系统内置的输出颜色的
   keyword: # 自定义关键词回复
     你好!: 你好!
     在: 当前主人不在,有事请留言!
     在?: 当前主人不在,有事请留言!
     在吗: 当前主人不在,有事请留言!
     在吗?: 当前主人不在,有事请留言!
   rootpath: /Users/xr/IDEA/Go-Wechaty-Bot/Server # 系统当前路劲,自动修改
   success: "[\e[01;32m✓\e[0m] " # 这是系统内置的输出颜色的
   tuling: # 图灵机器人
     token: 90984738f742454a98a0e92343371ec2&info= # 图灵机器人的 Token
     url: http://www.tuling123.com/openapi/api?key= # 图灵机器人的 URL
   wechaty: # 微信机器人
     url: 192.168.0.1:30000 # 微信机器人的 URL
     wechaty_puppet_service_token: insecure_46b80f25-12b3-4eb7-8c79-322398e413b9 # 微信机器人的 Token
   ```
   
2. Checking the Network Environment
   ```go
   if Gateway near Server{
   	IP=NATIP
   } else {
   	IP=InterNetIP
   }
   if PORT on {
   	continue
   } else {
   	os.exit(0)
   }
   ```
3. Start Server
   ```bash
   .\wechaty.bat # Start puppet-xp Gateway
   make server # Start Server
   ```


## ⚓️ ReCommended

1. [wechat-bot](https://github.com/cixingguangming55555/wechat-bot/blob/master/pic/doc.md)
2. [puppet-services](https://wechaty.js.org/docs/puppet-services/diy/#all-in-one-command)
3. [issues](https://github.com/wechaty/puppet-xp/issues/38)

## ⚠️ Notice

1. Network

   ```bash
   curl cip.cc
   curl -s https://api.chatie.io/v0/hosties/[WECHATY_TOKEN]
   ```

2. Help me improve this project

3. Submit bugs and interesting features

## Update Repo

1. TODO || example.yaml
2. Update.log

## Contact me

Reply to group chat with us

![](Image/bot.png)