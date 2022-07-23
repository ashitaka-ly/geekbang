package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"os"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("starting")
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/header", headerHandler)
	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/healthz", healthzHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		glog.V(2).Info(err)
	}
}

func headerHandler(w http.ResponseWriter, req *http.Request) {
	header := req.Header
	for key, values := range header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "getHeader")
}

func versionHandler(w http.ResponseWriter, req *http.Request) {
	version := os.Getenv("VERSION")
	w.Header().Add("VERSION", version)
	fmt.Fprintf(w, "getVersion")
}

func welcomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "welcome")
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "200")
}
