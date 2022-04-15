## Go-Wechaty-Bot PADLCOAL PROTOCOL

> 仅供学习使用，*请勿用于非法用途*！

[1]: https://img.shields.io/badge/puppet-xp-blue
[2]: https://img.shields.io/badge/puppet-padlocal-blue
[3]: https://img.shields.io/badge/puppet-4u-blue
[5]: https://github.com/XRSec/gobot/tree/xp
[6]: https://github.com/XRSec/gobot/tree/padlocal
[7]: https://github.com/XRSec/gobot/tree/4u

[![puppet-xp][1]][5] 〰️ [![puppet-padlocal][2]][6] 〰️ [![puppet-4u][3]][7] 「 Select Gateway 」

## Info

### Glossary

```ini
Gateway : puppet-padlocal
Server: go-wechaty
```

### Architecture

```mermaid
flowchart LR
    Polyglot-->Python
    Polyglot-->Go
    Polyglot -->Rust
    Python-->Grpc
    Go-->Grpc
    Rust-->Grpc
    Grpc-->Docker{Docker}-->padlocal
    padlocal-->微信
```

## ⇲ Use

### Init (depend main.General)

1. Checkout branch

   ```bash
   cd padlocal
   ```

2. Get & Generate Token

   ```bash
   WECHATY_PUPPET_PADLOCAL_TOKEN: http://pad-local.com/#/tokens
   WECHATY_TOKEN：curl -s https://www.uuidgenerator.net/api/version4
   WECHATY_PUPPET_SERVICE_TOKEN："insecure_" + WECHATY_TOKEN
   # WECHATY_TOKEN WECHATY_PUPPET_SERVICE_TOKEN 可同可不同
   ```

3. Modifying a Configuration File

   ```bash
   # wechatyGateway.sh
   export WECHATY_TOKEN="xxxxxxxxxxxxxxxx"
   export WECHATY_PUPPET_SERVICE_TOKEN="insecure_xxxxxxxxxxxxxxxx"
   export WECHATY_PUPPET_PADLOCAL_TOKEN="puppet_padlocal_xxxxxxxxxxxxxxxx"
   export WECHATY_PUPPET_SERVER_PORT="25000"
   ```

<!-- 4. Install the Packages

   ```bash
   # node-v16
   npm install -g wechaty-puppet-padlocal
   # npm --registry https://registry.npm.taobao.org install -g wechaty-puppet-padlocal
   ``` -->

### Start Server

```bash
cd padlocal && bash ./wechatyGateway.sh # Start puppet-xp Gateway
cd server && go run main.go # Start Server
```
