#!/bin/sh

set -x

cd /app || exit

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

# 是否开启 debug 模式
if [ -n "$DEBUG_ENABLE" ];then
  sed -i "s/Debug: true/Debug: $DEBUG_ENABLE/g" config.yaml
fi 

# 是否开启日志
if [ -n "$LOG_ENABLE" ];then
  sed -i "s/Log: true/Log: $LOG_ENABLE/g" config.yaml
fi 

# 统计数据过期时间 单位秒, 请输入整数 (无任何访问, 超过这个时间后, 统计数据将被清空, 0为不过期)
if [ -n "$EXPIRE_TIME" ];then
  sed -i "s/Expire: 2592000/Expire: $EXPIRE_TIME/g" config.yaml
fi

exec /app/busuanzi