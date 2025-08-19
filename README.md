[![soxft/busuanzi](https://socialify.cmds.run/soxft/busuanzi/image?description=1&font=Raleway&forks=1&language=1&logo=https%3A%2F%2Fraw.githubusercontent.com%2Fsoxft%2Fbusuanzi%2Fmain%2Fdist%2Ffavicon.png&name=1&owner=1&pattern=Circuit%20Board&stargazers=1&theme=Dark&cache=43200)](https://busuanzi.9420.ltd)

- [简体中文](README.zh_CN.md)

## self-hosted busuanzi

> A simple visitor statistics system based on Golang + Redis

- Calculate the UV and PV of the website
- Calculate the UV and PV of the subpage
- One-click deployment using Docker
- Privacy protection only stores HASH
- Pjax compatible webpage
- Support migration from the original busuanzi

## Installation

Support multiple running methods: compile and run from source code, run with Docker. See [Install](https://github.com/soxft/busuanzi/wiki/install) for details

### Quick Start with Docker

1. Edit the `docker-compose.yaml` file with your own configuration.
2. Run `docker-compose up -d` to start the service. 
3. Visit `http://localhost:8080` to view the data.

## Usage

Supports multiple custom attributes, compatible with pjax web pages, supports custom tag prefixes. See: [Usage documentation](https://github.com/soxft/busuanzi/wiki/usage)

## Principle

- `Busuanzi` supports both Redis and SQLite for data storage and retrieval. 
- Redis, as an in-memory database, has extremely high read and write performance. At the same time, its unique RDB and AOF persistence mechanisms ensure the security of Redis data.
- SQLite provides a lightweight, serverless database option that's perfect for smaller deployments or when you don't want to manage a separate Redis server.

### Database Configuration

You can choose between Redis and SQLite by setting the `Database.Type` configuration:

```yaml
Database:
  Type: redis     # Use 'redis' for Redis, 'sqlite' for SQLite
  Prefix: bsz     # Database prefix

# Redis configuration (when Type: redis)
Redis:
  Address: redis:6379
  Password: 
  Database: 0
  # ... other Redis settings

# SQLite configuration (when Type: sqlite)  
SQLite:
  Path: ./busuanzi.db  # SQLite database file path
```

### Data Storage

When using Redis, UV and PV data are stored in the following keys:

| index  | Types       | key                               |
| ------ | ----------- | --------------------------------- |
| sitePv | String      | bsz:site_pv:md5(host)             |
| siteUv | HyperLogLog | bsz:site_uv:md5(host)             |
| pagePv | ZSet        | bsz:page_pv:md5(host) / md5(path) |
| pageUv | HyperLogLog | bsz:site_uv:md5(host):md5(path)   |

When using SQLite, data is stored in the following tables:

| Table   | Purpose                    | Columns                              |
| ------- | -------------------------- | ------------------------------------ |
| site_pv | Site page views           | site_key, count                     |
| page_pv | Page views per path       | site_key, path_key, count           |
| site_uv | Site unique visitors      | site_key, user_hash                 |
| page_uv | Page unique visitors      | site_key, path_key, user_hash       |

## Data Migration

- You can use the [busuanzi-sync](https://github.com/soxft/busuanzi-sync) tool to sync data from the [original busuanzi](http://busuanzi.ibruce.info) to the self-hosted busuanzi.

## Upgrade Suggestions

- Please be sure to back up your data (dump.rdb) before upgrading.
- New and old version data may not be compatible, please pay attention to the instructions on the Release interface, upgrade cautiously
- 2.5.x - 2.7.x can use the [bsz-transfer](https://github.com/soxft/busuanzi-transfer) tool to migrate data to 2.8.x.

## Other

Logo created by ChatGPT

## Thanks

- CDN acceleration and security protection for this project are sponsored by [Tencent EdgeOne](https://edgeone.ai/?from=github).
- Thanks to [JetBrains](https://www.jetbrains.com/?from=busuanzi) for providing free student licenses for this project.

<p align="center">
    <a href="https://edgeone.ai/?from=github" style="margin-right: 24px; display: inline-block;">
        <img src="https://raw.githubusercontent.com/soxft/busuanzi/refs/heads/main/static/edgeone.png" alt="Tencent EdgeOne" width="200" style="vertical-align: middle; margin-right: 24px;"/>
    </a>
    <img src="https://resources.jetbrains.com.cn/storage/products/company/brand/logos/jetbrains.png" alt="JetBrains Logo" width="200" style="vertical-align: middle; margin-right: 24px;"/>
    <img src="https://resources.jetbrains.com.cn/storage/products/company/brand/logos/GoLand_icon.png" alt="GoLand Logo" width="50" style="vertical-align: middle;"/>
</p>