package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"time"
)

type Options struct {
	Port     int
	PacFile  string
	Proxy    string
	ProxyPac string
}

var DefaultOptions = Options{
	Port:     3128,
	PacFile:  "proxy.pac",
	Proxy:    "",
	ProxyPac: "",
}

// httputil.ReverseProxy implements basically everything for a baseline proxy.
//
//  1. It removes all hop-by-hop headers sent by the client.
//  2. Initiates TCP connection between backend, and sends the request.
//  3. Handles HTTP/1.1 upgrade request and its response (101), eventually by
//     hijacking to establish a TCP tunnel between the client and backend.
//  4. If the backend sends Trailer, the values are valid only after all body
//     is read. The standard implementation reads the entire body, and then
//     sends the trailer using Transfer-Encoding: chunked.
type Proxy struct {
	httputil.ReverseProxy
	Options *Options
}

// Create new proxy instance.
func NewProxy() *Proxy {
	return &Proxy{
		ReverseProxy: httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetXForwarded()
			},
			Transport: &http.Transport{
				// TODO: option should be used
				Proxy: http.ProxyFromEnvironment,
			},
		},
		Options: &DefaultOptions,
	}
}

func (p *Proxy) ProxyStart() {
	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", p.Options.Port),
		Handler: http.HandlerFunc(p.rootHandler),
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Println("proxy server listening on port:", p.Options.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		println(err.Error())
	}
}

func (p *Proxy) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		p.httpsHandler(w, r)
		return
	}
	if r.URL.Path == "/proxy.pac" {
		p.servePacFile(w, r)
		return
	}
	p.ServeHTTP(w, r)
}

func (p *Proxy) isProxySet() bool {
	return p.Options.Proxy != ""
}

//lint:ignore U1000 ""
func notImplemented(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "not implemented", http.StatusInternalServerError)
}
