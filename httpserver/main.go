package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"net/http"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("starting")
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/healthz", healthzHandler)
	//http.Handle("/healthz", http.HandleFunc(http.TimeoutHandler(healthzHandler, time.Second * 5, "")))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		glog.V(2).Info(err)
	}
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "welcome")
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	header := req.Header
	for key, values := range header {
		for index, value := range values {
			w.Header().Add(key, value)
			fmt.Println(index)
		}
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "200")
}
