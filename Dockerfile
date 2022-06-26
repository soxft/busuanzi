FROM golang:1.17-alpine as builder

LABEL maintainer="smallqi1@163.com"
WORKDIR /app
ENV GOPROXY=https://goproxy.cn
COPY . .
RUN go build -o busuanzi main.go

FROM alpine:3.16
WORKDIR /app

COPY --from=builder /app/busuanzi /app
COPY --from=builder /app/config.yaml /app
COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/ENTRYPOINT.sh /app

EXPOSE 8080
ENTRYPOINT [ "./ENTRYPOINT.sh", "/app/busuanzi"]