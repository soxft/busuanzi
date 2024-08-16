FROM golang:1.21-alpine AS builder
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN set -evx -o pipefail        \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories               \
    && apk add --no-cache git   \
    && rm -rf /var/cache/apk/*  \
    && GO111MODULE=on go build -o busuanzi main.go

FROM alpine:3.16
WORKDIR /app

COPY --from=builder /app/busuanzi /app
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/entrypoint.sh /app

EXPOSE 8080
ENTRYPOINT [ "sh", "entrypoint.sh" ]