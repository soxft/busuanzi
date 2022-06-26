# 不蒜子

> 自建不蒜子API
> 
> 基于 Golang + Redis 的简易访问量统计系统

  - 统计站点的 UV, PV
  - 统计文章页的 UV, PV

# 安装

> 详见 [WIKI](https://github.com/soxft/busuanzi/wiki/Install)

## Docker 部署

部署 Redis

```shell
docker run --name redis -d redis:latest
```

部署 busuanzi

```shell
docker run -p 8080:8080 --link redis -d qlucky/busuanzi:latest
```

可选参数：

| 环境变量     | 参数说明                                                     |
| ------------ | ------------------------------------------------------------ |
| API_SERVER   | busuanzi.js API地址，默认：http://127.0.0.1:8080/api，修改需要添加转符：https:\/\/busuanzi.xxx.com\/api |
| WEB_ADDR     | 默认 0.0.0.0:8080                                            |
| REDIS_ADDR   | redis 连接地址，默认 redis:6379                              |
| REDIS_PWD    | redis 密码，默认为空                                         |
| DEBUG_ENABLE | 是否开启debug模式，默认 true                                 |
| LOG_ENABLE   | 是否开启日志，默认 true                                      |
| EXPIRE_TIME  | 统计数据过期时间 单位秒, 请输入整数 (无任何访问, 超过这个时间后, 统计数据将被清空, 0为不过期)，默认 `2592000` |

修改默认参数示例：

```shell
docker run \
-e API_SERVER="https:\/\/busuanzi.xxx.com\/api" \
-e REDIS_ADDR=192.168.6.6:6379 \
-e REDIS_PWD=123456 \
-e DEBUG_ENABLE=true \
-e LOG_ENABLE=true \
-e EXPIRE_TIME=0 \
-p 8080:8080 \
--link redis \
-d qlucky/busuanzi:latest
```

挂载配置文件示例：

```shell
docker run \
-v ~/config.yaml:/app/config.yaml \
-p 8080:8080 \
--link redis \
-d qlucky/busuanzi:latest
```

## Docker 源码构建

只需安装docker环境即可

```shell
docker build -t busuanzi .
```