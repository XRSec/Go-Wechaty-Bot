#!/usr/bin/env bash
# set -ex

# http://pad-local.com/
export WECHATY_PUPPET="wechaty-puppet-padlocal"
export WECHATY_PUPPET_PADLOCAL_TOKEN="puppet_padlocal_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
# http://pad-local.com/

export WECHATY_TOKEN="insecure_xxxxxxxxxxxxxxxxxxxxxxxxx"
# export WECHATY_PUPPET_SERVICE_TOKEN="insecure_$(curl -s https://www.uuidgenerator.net/api/version4)"

export WECHATY_PUPPET_SERVER_PORT="25000"
export WECHATY_LOG="verbose" # silent | error | warn | info | verbose | silly
export WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER="true"

echo "WECHATY_TOKEN = ${WECHATY_TOKEN}"

# *********************PLAN 1***********************
if [ ! -n "$1" ]; then
    export tag="gateway"
else
    export tag=$1
fi

#docker images wechaty/wechaty | grep gateway | wc -l == 0

if [ "${tag}" == "gateway" ] && [ "$(docker images wechaty/wechaty | grep -c gateway)" == "0" ]; then
    docker build -t wechaty/wechaty:gateway .
fi

docker run -itd \
   --name Go-Wechaty-Bot \
   --restart=always \
   -e WECHATY_TOKEN \
   -e TZ=Asia/Shanghai \
   -e WECHATY_PUPPET_PADLOCAL_TOKEN \
   -e WECHATY_LOG \
   -e WECHATY_PUPPET \
   -e WECHATY_PUPPET_SERVER_PORT \
   -e WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER \
   -p "${WECHATY_PUPPET_SERVER_PORT}:${WECHATY_PUPPET_SERVER_PORT}" \
   wechaty/wechaty:"${tag}"
# *********************PLAN 1***********************


# *********************PLAN 2***********************
#npm run serve
# *********************PLAN 2***********************
