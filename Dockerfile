FROM golang:1.18-alpine as builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN set -evx -o pipefail        \
    && apk update               \
    && apk add --no-cache git   \
    && rm -rf /var/cache/apk/*  \
    && go build -ldflags="-s -w" -o busuanzi main.go


FROM alpine:3.16
WORKDIR /app

COPY --from=builder /app/busuanzi /app
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/entrypoint.sh /app

EXPOSE 8080
ENTRYPOINT [ "./entrypoint.sh" ]