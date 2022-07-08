#!/usr/bin/env bash
# set -ex

export WECHATY_PUPPET="wechaty-puppet-wechat"
export WECHATY_PUPPET_PADLOCAL_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxx"
# export WECHATY_PUPPET_PADLOCAL_TOKEN="$(curl -s https://www.uuidgenerator.net/api/version4)"

export WECHATY_TOKEN="insecure_xxxxxxxxxxxxxxxxxxxxxxxxx"
# export WECHATY_PUPPET_SERVICE_TOKEN="insecure_$(curl -s https://www.uuidgenerator.net/api/version4)"
export WECHATY_PUPPET_SERVER_PORT="25000"
export WECHATY_LOG="verbose" # silent | error | warn | info | verbose | silly
export WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER="true"

if [ ! -n "$1" ]; then
    export tag=0.56
else
    export tag=$1
fi

touch "$(pwd)/${WECHATY_TOKEN}.memory-card.json"

docker run -ti --rm \
    --name wechatBot \
    -e WECHATY_TOKEN \
    -e TZ=Asia/Shanghai \
    -e WECHATY_PUPPET_SERVICE_TOKEN \
    -e WECHATY_LOG \
    -e WECHATY_PUPPET \
    -e WECHATY_PUPPET_SERVER_PORT \
    -e WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER \
    -p "${WECHATY_PUPPET_SERVER_PORT}:${WECHATY_PUPPET_SERVER_PORT}" \
    -v "$(pwd)/${WECHATY_TOKEN}.memory-card.json:/wechaty/${WECHATY_TOKEN}.memory-card.json" \
    wechaty/wechaty:"${tag}"