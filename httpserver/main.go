package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("program starting")

	mux := http.NewServeMux()
	// 加入 debug 信息
	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// 作业中 1 2 3
	mux.HandleFunc("/", rootHandler)
	// 4
	mux.HandleFunc("/healthz", healthzHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("start http server failed, error: %s\n", err.Error())
	}
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
