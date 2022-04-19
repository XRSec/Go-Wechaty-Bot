#!/usr/bin/env bash
# set -ex

export WECHATY_PUPPET="wechaty-puppet-padlocal"
export WECHATY_PUPPET_PADLOCAL_TOKEN="puppet_padlocal_xxxxxxxxxxxxxxxxxxxxxxxxx"
# http://pad-local.com/

export WECHATY_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxx"
# export WECHATY_TOKEN="$(curl -s https://www.uuidgenerator.net/api/version4)"

export WECHATY_PUPPET_SERVICE_TOKEN="insecure_xxxxxxxxxxxxxxxxxxxxxxxxx"
# export WECHATY_PUPPET_SERVICE_TOKEN="insecure_$(curl -s https://www.uuidgenerator.net/api/version4)"

export WECHATY_PUPPET_SERVER_PORT="25000"
export WECHATY_LOG="verbose" # silent | error | warn | info | verbose | silly
export WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER="true"

echo "WECHATY_PUPPET_PADLOCAL_TOKEN = ${WECHATY_PUPPET_PADLOCAL_TOKEN}"

if [ ! -n "$1" ]; then
    export tag=latest
else
    export tag=$1
fi

docker run -ti --rm \
    --name wechatBot \
    -e WECHATY_TOKEN \
    -e WECHATY_PUPPET_SERVICE_TOKEN \
    -e WECHATY_PUPPET_PADLOCAL_TOKEN \
    -e WECHATY_LOG \
    -e WECHATY_PUPPET \
    -e WECHATY_PUPPET_SERVER_PORT \
    -e WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER \
    -p "${WECHATY_PUPPET_SERVER_PORT}:${WECHATY_PUPPET_SERVER_PORT}" \
    wechaty/wechaty:"${tag}"