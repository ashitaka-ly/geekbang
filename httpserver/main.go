package main

import (
	"flag"
	"github.com/golang/glog"
	"net/http"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("starting")
	http.HandleFunc("/", rootHandler)

}

func rootHandler(writer http.ResponseWriter, request *http.Request) {

}
