# 不蒜子

> 自建不蒜子API
> 
> 基于 Golang + Redis 的简易访问量统计系统

  - 统计站点的 UV, PV
  - 统计文章页的 UV, PV

# 安装

1. 在release界面下载源码
2. 进入源码目录
3. 在终端执行 go get 安装依赖
4. 配置config.yml
5. 在终端执行 go build 编译为二进制文件
6. 编辑 dist/busuanzi.js 替换链接为自己的, 也可以编辑ts文件自行编译
7. 运行二进制文件
