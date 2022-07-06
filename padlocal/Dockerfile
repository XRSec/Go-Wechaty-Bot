FROM  node:16
LABEL maintainer="xrsec"
LABEL mail="troy@zygd.site"
LABEL Github="https://github.com/XRSec/Go-Wechaty-Bot"
LABEL org.opencontainers.image.source="https://github.com/XRSec/Go-Wechaty-Bot"
LABEL org.opencontainers.image.title="Go-Wechaty-Bot"

WORKDIR /opt/go-wechaty-bot
COPY . .

RUN npm install --registry=https://registry.npmmirror.com

WORKDIR /opt/go-wechaty-bot/node_modules/.bin/

EXPOSE 25000
ENV NODE_ENV production
ENV WECHATY_PUPPET wechaty-puppet-padlocal
ENV WECHATY_PUPPET_SERVER_PORT 25000
ENV WECHATY_TOKEN insecure_xxxxxxxxxxxxxxxxxxxxxxxxx
ENV WECHATY_PUPPET_PADLOCAL_TOKEN puppet_padlocal_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
ENV LANG 'zh_CN.UTF-8'
STOPSIGNAL SIGQUIT

CMD ["npm", "run", "serve"]