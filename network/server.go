package network

import (
	"net/http"
	"strings"
	"time"

	"github.com/efjoubert/goforit/goblet"
)

//Server conveniance struct wrapping arround *http.Server
type Server struct {
	port        string
	svr         *http.Server
	istls       bool
	certFile    string
	keyFile     string
	srvHTTPCall ServeHTTPCall
	serveQueue  chan func()
	active      chan bool
}

//NewServer return *Server
//istls is tls listener
//certfile string - path to certificate file
//keyFile string - path to key file
func NewServer(port string, istls bool, certFile string, keyFile string, srvHttPReqstCall ...ServeHTTPCall) *Server {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	var svr = &Server{port: port, istls: istls, certFile: certFile, keyFile: keyFile, serveQueue: make(chan func()), active: make(chan bool, 1)}
	if len(srvHttPReqstCall) == 1 && srvHttPReqstCall[0] != nil {
		svr.srvHTTPCall = srvHttPReqstCall[0]
	} else {
		svr.srvHTTPCall = DefaultServeHTTPCall
	}
	return svr
}

//InvokeAndListen Invoke Server and Start Listening
func InvokeAndListen(port string, istls bool, certFile string, keyFile string, srvHttPReqstCall ...ServeHTTPCall) (err error) {
	var srvr = NewServer(port, istls, certFile, keyFile, srvHttPReqstCall...)
	err = srvr.Listen()
	srvr = nil
	return
}

func (svr *Server) processServes() {
	for {
		select {
		case serve := <-svr.serveQueue:
			go serve()
		case atv := <-svr.active:
			if atv {
				svr.active <- true
				return
			}
		}
	}
}

//Listen start listening for connection(s)
func (svr *Server) Listen() (err error) {
	svr.svr = &http.Server{Addr: svr.port, Handler: svr, ReadHeaderTimeout: 3 * 1024 * time.Millisecond, ReadTimeout: 30 * 1024 * time.Millisecond, WriteTimeout: 60 * 1024 * time.Millisecond}
	if svr.istls {
		go svr.processServes()
		err = svr.svr.ListenAndServeTLS(svr.certFile, svr.keyFile)
		svr.active <- true
		if <-svr.active {
			svr.svr.Close()
		}
		svr.svr = nil
	} else {
		go svr.processServes()
		err = svr.svr.ListenAndServe()
		svr.active <- true
		if <-svr.active {
			svr.svr.Close()
		}
		svr.svr = nil
	}
	close(svr.active)
	close(svr.serveQueue)
	return err
}

//ServeHTTP server http.Handler interface implementation of ServeHTTP
func (svr *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool, 1)
	defer close(done)
	svr.serveQueue <- func() {
		svr.srvHTTPCall(svr, w, r)
		done <- true
	}
	<-done
}

//ServeHTTPCall ServeHTTPCall
type ServeHTTPCall = func(srv *Server, w http.ResponseWriter, r *http.Request)

//DefaultServeHTTPCall DefaultServeHTTPCall
func DefaultServeHTTPCall(srv *Server, w http.ResponseWriter, r *http.Request) {
	goblet.PerformHTTPServletRequest(w, r)
}
