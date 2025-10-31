[![soxft/busuanzi](https://socialify.cmds.run/soxft/busuanzi/image?description=1&font=Raleway&forks=1&language=1&logo=https%3A%2F%2Fraw.githubusercontent.com%2Fsoxft%2Fbusuanzi%2Fmain%2Fdist%2Ffavicon.png&name=1&owner=1&pattern=Circuit%20Board&stargazers=1&theme=Dark&cache=43200)](https://busuanzi.9420.ltd)

## 自建不蒜子

> 一个基于 Golang + Redis 的简易访问量统计系统

- 统计站点的 UV, PV
- 统计子页面的 UV, PV
- 使用 Docker 一键部署
- 隐私保障 仅存储 HASH
- 兼容 Pjax 技术的网页
- 支持从原版不蒜子迁移数据

## 安装

支持多种运行方式: 源码编译运行, Docker 运行. 详见: [安装方法](https://gitee.com/soxft/busuanzi/wikis/install)

## 使用方式

支持多种自定义属性, 兼容 pjax 网页, 支持自定义 标签前缀. 详见: [使用文档](https://gitee.com/soxft/busuanzi/wikis/usage)

## 原理

- `Busuanzi` 使用 Redis 进行数据存储与检索。Redis 作为内存数据库拥有极高的读写性能，同时其独特的`RDB`与`AOF`持久化方式，使得 Redis 的数据安全得到保障。

- UV 与 PV 数据分别采用以下方式进行存储:

| index  | 数据类型        | key                               |
|--------|-------------|-----------------------------------|
| sitePv | String      | bsz:site_pv:md5(host)             |
| siteUv | HyperLogLog | bsz:site_uv:md5(host)             |
| pagePv | ZSet        | bsz:page_pv:md5(host) / md5(path) |
| pageUv | HyperLogLog | bsz:site_uv:md5(host):md5(path)   |

## SQLite 持久化

Busuanzi 可以按照指定的时间间隔，将 Redis 中的 PV、UV 与 HyperLogLog 数据自动快照到
SQLite 数据库文件中，便于长期保存与快速恢复。所有配置项均支持使用环境变量覆盖：

| 配置项 | 类型 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `Persistence.Enable` | bool | `false` | 是否开启后台快照任务 |
| `Persistence.DBPath` | string | `data/busuanzi.db` | SQLite 数据库文件路径 |
| `Persistence.Interval` | number | `300` | 快照间隔（秒） |
| `Persistence.RestoreOnStart` | bool | `false` | 启动时是否自动恢复最近一次快照 |

快照会保存每个统计键当前的数值、序列化的 HyperLogLog 数据以及剩余 TTL。恢复快照时
会用存储的数据覆盖 Redis 中对应的键。

### 快照数据映射

不同的 Redis 数据类型会按下表映射到 SQLite `snapshots` 表中：

| Redis 键 | Redis 类型 | SQLite 字段 | 说明 |
| --- | --- | --- | --- |
| `site_pv` | 字符串计数器 | `count` | 直接保存站点 PV 的整数值。 |
| `site_uv` | HyperLogLog | `count`, `payload` | `payload` 中保存 `DUMP` 返回的原始二进制，并以十六进制编码。 |
| `page_pv` | ZSet | `count`, `payload` | `payload` 是 `{path_unique,count}` 组成的 JSON 数组，用于恢复排名明细。 |
| `page_uv` | HyperLogLog | `count`, `payload` | 与 `site_uv` 相同的序列化方式。 |

`ttl_ms` 字段会记录剩余的过期时间（毫秒，永不过期时为 `0`），以便在恢复时继续保持 Redis 的 TTL 语义。


## 其他

Logo 由 ChatGPT 设计

## 数据迁移

- 可使用 [busuanzi-sync](https://gitee.com/soxft/busuanzi-sync) 工具迁移[原版不蒜子](http://busuanzi.ibruce.info)的数据

## 升级建议

- 请务必在升级前备份数据 (dump.rdb)
- 新老版本数据可能并不兼容, 请注意 Release 界面的说明, 谨慎升级
- 2.5.x - 2.7.x 可以使用 [bsz-transfer](https://gitee.com/soxft/busuanzi-transfer) 工具进行数据迁移至 2.8.x