package main

import (
	"flag"
	"fmt"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

func main() {

	klog.InitFlags(nil)
	flag.Set("v", "1")
	flag.Parse()
	defer klog.Flush()

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthZ", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "200")
	})
	mux.HandleFunc("/401", badRequest)
	mux.HandleFunc("/500", errorRequest)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	app := New("0.0.0.0:8000", Log(mux))
	time.AfterFunc(time.Second, func() {
		app.stop()
	})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
