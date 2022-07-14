#!/bin/sh

set -x

cd /app || exit

if [ -f "config.yaml" ]; then
  # busuanzi js API address
  if [ -n "$API_SERVER" ];then
    sed -i "s/http:\/\/127.0.0.1:8080\/api/$API_SERVER/g" dist/busuanzi.js
  fi

  # redis地址
  if [ -n "$REDIS_ADDR" ];then
    sed -i "s/Address: 127.0.0.1:6379/Address: $REDIS_ADDR/g" config.yaml
  else
    sed -i "s/Address: 127.0.0.1:6379/Address: redis:6379/g" config.yaml
  fi

  # redis 密码
  if [ -n "$REDIS_PWD" ];then
    sed -i "s/Password:/Password: $REDIS_PWD/g" config.yaml
  fi

  # 是否开启日志
  if [ -n "$LOG_ENABLE" ];then
    sed -i "s/Log: true/Log: $LOG_ENABLE/g" config.yaml
  fi

  mv config.yaml /app/expose/config.yaml
  mv dist /app/expose/dist
fi

exec /app/busuanzi -c ./expose/config.yaml -d ./expose/dist