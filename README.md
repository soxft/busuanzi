[![soxft/busuanzi](https://socialify.git.ci/soxft/busuanzi/image?description=1&font=Bitter&forks=1&language=1&logo=https://raw.githubusercontent.com/soxft/busuanzi/main/dist/favicon.svg&name=1&owner=1&pattern=Solid&stargazers=1&theme=Dark)](https://busuanzi.9420.ltd)

## 自建不蒜子

> 一个基于 Golang + Redis 的简易访问量统计系统
>
> A simple visitor statistics system based on Golang + Redis

- 统计站点的 UV, PV
- 统计子页面的 UV, PV
- 使用 Docker 一键部署
- 隐私保障 仅存储 HASH

## 安装

支持多种运行方式: 源码编译运行, Docker 运行. 详见: [Install](https://github.com/soxft/busuanzi/wiki/install)

## 使用方式

支持多种自定义属性, 兼容 pjax 网页, 支持自定义 标签前缀. 详见: [使用文档](https://github.com/soxft/busuanzi/wiki/usage)

## 原理

`Busuanzi` 使用 Redis 进行数据存储于检索。Redis 作为内存数据库拥有极高的读写性能，同时其独特的`RDB`与`AOF`持久化方式，使得 Redis 的数据安全得到保障。

UV 与 PV 数据分别采用以下方式进行存储:

| index  | 数据类型        | key                               |
|--------|-------------|-----------------------------------|
| sitePv | String      | bsz:site_pv:md5(host)             |
| siteUv | HyperLogLog | bsz:site_uv:md5(host)             |
| pagePv | ZSet        | bsz:page_pv:md5(host) / md5(path) |
| pageUv | HyperLogLog | bsz:site_uv:md5(host):md5(path)   |

## 其他

[Yuantuo](https://github.com/yuantuo666) 提供了一个 支持 Web 管理的版本. 可以在 [yuantuo666/busuanzi](https://github.com/yuantuo666/busuanzi) 找到

## 升级建议

- 请务必在升级前备份数据 (dump.rdb)
- 新老版本数据可能并不兼容, 请注意 Release 界面的说明, 谨慎升级