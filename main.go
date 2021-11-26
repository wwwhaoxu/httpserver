package main

import (
	"flag"
	"fmt"
	"httpserver/pkg/config"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"
)

var (
	// flagConf is the config flag.
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "/app/config/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {

	klog.InitFlags(nil)
	flag.Set("v", "2")
	flag.Parse()
	defer klog.Flush()

	kvs := make(map[interface{}]interface{})
	conf := config.NewConfig(kvs, flagConf)
	if err := conf.LoadFile(); err != nil {
		panic(err)
	}

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

	ipPort := strings.Join([]string{"0.0.0.0", kvs["port"].(string)}, ":")
	app := New(ipPort, Log(mux))
	time.AfterFunc(time.Second, func() {
		app.stop()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
