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

- `Busuanzi` uses Redis for data storage and retrieval. Redis, as an in-memory database, has extremely high read and write performance. At the same time, its unique RDB and AOF persistence mechanisms ensure the security of Redis data.

UV and PV data are stored in the following keys:

| index  | Types       | key                               |
| ------ | ----------- | --------------------------------- |
| sitePv | String      | bsz:site_pv:md5(host)             |
| siteUv | HyperLogLog | bsz:site_uv:md5(host)             |
| pagePv | ZSet        | bsz:page_pv:md5(host) / md5(path) |
| pageUv | HyperLogLog | bsz:site_uv:md5(host):md5(path)   |

## SQLite persistence

Busuanzi can periodically snapshot Redis data into a SQLite database in order to keep
long term history or recover from Redis failures. The following configuration options
control the behaviour (all values can also be provided via environment variables):

| key | type | default | description |
| --- | --- | --- | --- |
| `Persistence.Enable` | bool | `false` | Enable background SQLite snapshots. |
| `Persistence.DBPath` | string | `data/busuanzi.db` | Path to the SQLite database file. |
| `Persistence.Interval` | number | `300` | Snapshot interval in seconds. |
| `Persistence.RestoreOnStart` | bool | `false` | Restore the most recent snapshot during startup. |

Snapshots include PV, UV and HyperLogLog states. Restoring a snapshot replaces the
corresponding Redis keys using the stored values and TTL information.

### Snapshot payload mapping

Each Redis data type is normalised into the SQLite `snapshots` table as follows:

| Redis key | Redis type | SQLite fields | Notes |
| --- | --- | --- | --- |
| `site_pv` | String counter | `count` | The plain integer PV value. |
| `site_uv` | HyperLogLog | `count`, `payload` | `payload` stores the raw bytes returned by `DUMP`, encoded as hexadecimal. |
| `page_pv` | ZSet | `count`, `payload` | `payload` is a JSON array of `{path_unique,count}` pairs so that the ranking can be restored. |
| `page_uv` | HyperLogLog | `count`, `payload` | Same as `site_uv`. |

The `ttl_ms` column stores the remaining key lifetime in milliseconds (or `0` for keys without an expiration) so that the Redis
TTL semantics are preserved when restoring a snapshot.

## Data Migration

- You can use the [busuanzi-sync](https://github.com/soxft/busuanzi-sync) tool to sync data from the [original busuanzi](http://busuanzi.ibruce.info) to the self-hosted busuanzi.

## Upgrade Suggestions

- Please be sure to back up your data (dump.rdb) before upgrading.
- New and old version data may not be compatible, please pay attention to the instructions on the Release interface, upgrade cautiously
- 2.5.x - 2.7.x can use the [bsz-transfer](https://github.com/soxft/busuanzi-transfer) tool to migrate data to 2.8.x.

## Other

Logo created by ChatGPT

## Thanks

- CDN acceleration and security protection for this project are sponsored by Tencent EdgeOne: EdgeOne offers a long-term free plan with unlimited traffic and requests, covering Mainland China nodes, with no overage charges. Interested friends can click the link below to claim it ->
    - [Best Asian CDN, Edge, and Secure Solutions - Tencent EdgeOne](https://edgeone.ai/?from=github).
- Thanks to [JetBrains](https://www.jetbrains.com/?from=busuanzi) for providing free student licenses for this project.

<p align="center">
    <a href="https://edgeone.ai/?from=github" style="margin-right: 24px; display: inline-block;">
        <img src="https://raw.githubusercontent.com/soxft/busuanzi/refs/heads/main/static/edgeone.png" alt="Tencent EdgeOne" width="200" style="vertical-align: middle; margin-right: 24px;"/>
    </a>
    <img src="https://resources.jetbrains.com.cn/storage/products/company/brand/logos/jetbrains.png" alt="JetBrains Logo" width="200" style="vertical-align: middle; margin-right: 24px;"/>
    <img src="https://resources.jetbrains.com.cn/storage/products/company/brand/logos/GoLand_icon.png" alt="GoLand Logo" width="50" style="vertical-align: middle;"/>
</p>
