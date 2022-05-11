## Go-Wechaty-Bot

> 仅供学习使用，*请勿用于非法用途*！

[1]: https://img.shields.io/badge/puppet-xp-blue
[2]: https://img.shields.io/badge/puppet-padlocal-blue
[3]: https://img.shields.io/badge/puppet-4u-blue
[5]: /xp#go-wechaty-bot-xp-protocol
[6]: /padlocal#go-wechaty-bot-padlcoal-protocol
[7]: /4u#go-wechaty-bot-4u-protocol

[![puppet-xp][1]][5] 〰️ [![puppet-padlocal][2]][6] 〰️ [![puppet-4u][3]][7] 「 Select Gateway 」

## Architecture

```mermaid
flowchart LR
    Polyglot-->Python
    Polyglot-->Go
    Polyglot -->Rust
    Python-->Grpc
    Go-->Grpc
    Rust-->Grpc
    Grpc-->Gateway{Gateway}-->Puppet{Puppet}
    Grpc-->Docker{Docker}-->Puppet{Puppet}
  Puppet{Puppet}-->xp-->微信
  Puppet{Puppet}-->padlcoal-->微信
    Puppet{Puppet}-->wechat4u-->微信
```

## General

1. Clone Repo

   ```bash
   git clone https://github.com/XRSec/gobot.git wechatbot
   cd wechatbot
   ```

2. Install the Packages

   ```bash
   # node-v16
   npm --registry http://registry.npmmirror.com install -g wechaty
   ```

3. Edit `Server/config.yaml`.

   ```yaml
   bot:
     adminid: wxid_xxxxx
     name: xxxxxxxx
   ding:
     keyword: Wechaty
     token: xxxxxxxxxxxxxxxxxx
     url: https://oapi.dingtalk.com/robot/send?access_token=
   tuling:
     token: xxxxxxxxxxxxxxxx&info=
     url: http://www.tuling123.com/openapi/api?key=
   wechaty:
     wechaty_puppet_endpoint: 127.0.0.1:25001
     wechaty_puppet_service_token: insecure_xxxxxxxxxxxxxxxxxxxxxx
   wxopenai:
     env: online
     token: xxxxxxxxxxxxxxxxxxxxx
   ```

4. Checking the Network Environment

  ```go
  if Gateway near Server {
  IP = NAT_IP
  } else {
  IP = InterNet_IP
  }
  if PORT on {
  continue
  } else {
  os.exit(0)
  }
  ```

## ⚓️ Re

1. [wechat-bot](https://github.com/cixingguangming55555/wechat-bot/blob/master/pic/doc.md)

2. [puppet-services](https://wechaty.js.org/docs/puppet-services/diy/#all-in-one-command)

3. [issues](https://github.com/wechaty/puppet-xp/issues/38)

## ⚠️ Debug

1. Network

   ```bash
   curl cip.cc
   curl -s https://api.chatie.io/v0/hosties/[WECHATY_TOKEN]
   ```

2. Help me improve this project

3. Submit bugs and interesting features

## Update Repo

1. TODO && example.yaml
2. Update.md && Log.md

## Contact me

Reply to group chat with us

![wxid: XRSec_MSG](Image/bot.png)
