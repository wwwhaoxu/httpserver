package main

import (
	"context"
	"errors"
	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
	"httpserver/pkg/metrics"
	"io"
	"k8s.io/klog/v2"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

//func (o *responseObserver) Write(p []byte) (n int, err error) {
//	if !o.wroteHeader {
//		o.WriteHeader(http.StatusOK)
//	}
//	n, err = o.ResponseWriter.Write(p)
//	o.written += int64(n)
//	return
//}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

// Logs incoming requests.

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w}
		defer func(o *responseObserver) {
			if p := recover(); p != nil {
				klog.Error(p)
				o.WriteHeader(500)
			}
			if o.status >= 400 {
				klog.Errorf("[URI:%s][ClientIP:%s][Status:%d]", r.URL.Path, ClientIP(r), o.status)
			} else {
				klog.Infof("[URI:%s][ClientIP:%s][Status:%d]", r.URL.Path, ClientIP(r), o.status)
			}
		}(o)
		o.status = 200
		h.ServeHTTP(o, r)

	})
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(401)
}

func errorRequest(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(500)
}

func ClientIP(r *http.Request) string {
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
		if ip == "::1" {
			return "127.0.0.1"
		}
		return ip
	}

	return ""
}

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func serve(addr string, handler http.Handler) error {
	s := http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	return s.ListenAndServe()
}

type App struct {
	addr    string
	ctx     context.Context
	cancel  func()
	handler http.Handler
}

func New(addr string, handler http.Handler) *App {

	ctx, cancle := context.WithCancel(context.Background())

	return &App{
		addr:    addr,
		ctx:     ctx,
		cancel:  cancle,
		handler: handler,
	}
}

func (a *App) Run() error {
	g, ctx := errgroup.WithContext(a.ctx)
	for {
		g.Go(func() error {
			return serve(a.addr, a.handler)
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, shutdownSignals...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				a.stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	klog.V(4).Info("entering root handler")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	delay := randInt(10, 200)
	time.Sleep(time.Millisecond * time.Duration(delay))

	io.WriteString(w, "=========start send requests======")
	service, _ := os.LookupEnv("SERVICE")
	req, err := http.NewRequest("GET", "http://"+service, nil)
	if err != nil {
		klog.Warningf("request failed %s", service)
	}

	lowerCaseHeader := make(http.Header)
	for key, value := range r.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	klog.V(4).Info("headers:", lowerCaseHeader)

	req.Header = lowerCaseHeader
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		glog.Info("HTTP get failed with error: ", "error", err)
	} else {
		glog.Info("HTTP get succeeded")
	}
	if resp != nil {
		resp.Write(w)
	}
	//if r.URL.Path != "/" {
	//	w.WriteHeader(404)
	//	return
	//}
	//for k, v := range r.Header {
	//	for _, s := range v {
	//		if k == "Name" {
	//			w.Header().Add(k, s)
	//			io.WriteString(w, fmt.Sprintf("%s=%s\n", k, s))
	//			klog.V(4).Info(k, s)
	//		}
	//	}
	//
	//}
	//if env, b := os.LookupEnv("VERSION"); b {
	//	w.Header().Set("VERSION", env)
	//}
	klog.V(4).Infof("Respond in %d ms", delay)

}
