FROM golang:1.17 AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY . .
RUN go build -o app httpserver/main.go


FROM busybox
COPY --from=builder /app/app /app/httpserver
WORKDIR /app/
EXPOSE 8080
ENTRYPOINT [ "./httpserver" ]




