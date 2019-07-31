package goblet

import "net/http"

//Request Request
type Request struct {
	params   *Parameters
	method   string
	url      string
	protocol string
	srvlt    *Servlet
	r        *http.Request
}

//Header Request Header
func (reqst *Request) Header(name string) (val string) {
	return reqst.r.Header.Get(name)
}

func newHTTPRequest(srvlt *Servlet, r *http.Request, requestpath string) (reqst *Request) {
	reqst = &Request{params: NewParameters(), method: r.Method, url: requestpath, protocol: r.Proto, srvlt: srvlt, r: r}
	loadParametersFromHTTPRequest(reqst.params, r)
	return
}

func cleanupRequest(reqst *Request) {
	if reqst.params != nil {
		reqst.params.CleanupParameters()
		reqst.params = nil
	}
	if reqst.srvlt != nil {
		reqst.srvlt = nil
	}
	if reqst.r != nil {
		reqst.r = nil
	}
}
