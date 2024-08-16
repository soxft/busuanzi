#!/bin/sh

set -x

cd /app || exit

if [ ! -f "/app/install.lock" ];then
  # busuanzi js API address
  if [ -n "$API_SERVER" ];then
    sed -i "s/http:\/\/127.0.0.1:8080\/api/$API_SERVER/g" dist/busuanzi.js
  fi

  mv dist /app/dist
  mv config.yaml /app/config.yaml
  touch /app/install.lock
fi

exec /app/busuanzi -c /app/config.yaml -d /app/dist