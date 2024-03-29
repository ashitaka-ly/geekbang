package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ashitaka/geekbang/httpserver/metrics"
)

func main() {

	// 日志等级
	// flag.Set("v", "4")
	flag.Parse()
	defer glog.Flush()
	glog.V(3).Info("program starting")
	glog.V(4).Info("DEBUG LEVEL 4")

	metrics.Register()

	// 利用 sigterm 信号关闭服务
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 构建并启动 http server
	srv := buildWebServer()
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			glog.Fatalf("start http server failed, error: %s\n", err.Error())
		}
	}()
	glog.V(3).Info("Server Starting")

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 停止 http server
	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatalf("Server Shutdown Failed:%+v", err)
	}
	glog.V(3).Info("Server Shutdown")
}

// 构建 http server
func buildWebServer() *http.Server {
	glog.V(4).Info("DEBUG msg")
	mux := http.NewServeMux()
	// 加入 debug 信息
	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// root 基本业务
	mux.HandleFunc("/root", rootHandler)
	mux.HandleFunc("/now", timeHandler)
	// 暴露给 prometheus
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/hello", helloHandler)
	// 探活
	mux.HandleFunc("/healthz", healthzHandler)

	srv := &http.Server{
		Addr:    ":80",
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
	glog.V(3).Infof("Success! client ip %s", clientIp)

	w.WriteHeader(200)
	fmt.Fprintf(w, "SUCCESS")
}

func timeHandler(w http.ResponseWriter, req *http.Request) {
	now := time.Now()      //获取当前时间
	year := now.Year()     //年
	month := now.Month()   //月
	day := now.Day()       //日
	hour := now.Hour()     //小时
	minute := now.Minute() //分钟
	second := now.Second() //秒
	fmt.Fprintf(w, "time is [%d-%02d-%02d %02d:%02d:%02d]", year, month, day, hour, minute, second)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(3).Infoln("entering hello")
	et := metrics.NewTimer()
	defer et.ObserveTotal()
	delay := randInt(20, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	glog.V(3).Infof("delay = %d", delay)
	fmt.Fprintf(w, "delay = %d", delay)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
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
