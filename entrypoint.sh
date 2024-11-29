#!/bin/sh

set -x

cd /app || exit

if [ ! -f "/app/install.lock" ];then
  # busuanzi js API address
  if [ -n "$API_SERVER" ];then
    sed -i "s|http://127.0.0.1:8080/api|$API_SERVER|g" dist/*.js

    echo "Replace API_SERVER to $API_SERVER"
  fi

  # busuanzi js API address
    if [ -n "$API_JSONP_SERVER" ];then
      sed -i "s|http://127.0.0.1:8080/jsonp|$API_JSONP_SERVER|g" dist/busuanzi.jsonp.js

      echo "Replace API_JSONP_SERVER to $API_JSONP_SERVER"
    fi

  touch /app/install.lock
fi

exec /app/busuanzi -c /app/config.yaml -d /app/dist