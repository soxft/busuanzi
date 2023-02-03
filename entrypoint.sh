#!/bin/sh

set -x

cd /app || exit

if [ -f "install.lock" ];then
  # busuanzi js API address
  if [ -n "$API_SERVER" ];then
    sed -i "s/http:\/\/127.0.0.1:8080\/api/$API_SERVER/g" dist/busuanzi.js
  fi

  mv dist /app/expose/dist
  mv config.yaml /app/expose/config.yaml
fi

touch install.lock

exec /app/busuanzi -c ./expose/config.yaml -d ./expose/dist