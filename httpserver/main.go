package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
)

func main() {
	// 日志等级
	flag.Set("v", "4")
	glog.V(2).Info("program starting")

	// 利用 sigterm 信号关闭服务
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 构建并启动 http server
	srv := buildWebServer()
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("start http server failed, error: %s\n", err.Error())
		}
	}()
	glog.V(2).Info("Server Starting")

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 停止 http server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Shutdown")
}

// 构建 http server
func buildWebServer() *http.Server {
	mux := http.NewServeMux()
	// 加入 debug 信息
	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// root 基本业务
	mux.HandleFunc("/", rootHandler)
	// 探活
	mux.HandleFunc("/healthz", healthzHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return srv
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	// header 信息写入
	header := req.Header
	for key, values := range header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 获取 version 配置
	os.Setenv("VERSION", "v0.0.1")
	version := os.Getenv("VERSION")
	w.Header().Add("VERSION", version)

	// 获取 IP 地址
	clientIp := getIp(req)
	log.Printf("Success! client ip %s", clientIp)

	w.WriteHeader(200)
	fmt.Fprintf(w, "SUCCESS")
}

// 获取 IP 地址
func getIp(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "running")
}
