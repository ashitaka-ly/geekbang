FROM golang:1.17 AS builder

## 指定环境变量，这里开启 gomod 是因为这个项目用了 go mod
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

## 指定工作目录，并且把当前目录整个 copy 过去，
## 因为开启了 go mod，所以不能只 copy httpserver 目录
WORKDIR /app
COPY . .
## 在工作目录下 build httpserver 下的应用，应用名 app
RUN go build -o app httpserver/main.go

## 把上阶段的 build 的 app 可执行应用，复制到容器指定位置
## 注意在 docker run 时要指定端口，-p 8080:8080, 否则本机不能直接访问
FROM busybox
COPY --from=builder /app/app /app/httpserver
WORKDIR /app/
EXPOSE 8080
ENTRYPOINT ["./httpserver -v=3 -log_dir=/data/log" ]




