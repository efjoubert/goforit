package goblet

import (
	"net/http"

	"../goio"
)

//Response Response
type Response struct {
	w     http.ResponseWriter
	srvlt *Servlet
	reqst *Request
	iw    *goio.IORW
}

//AddHeader Response
func (resp *Response) AddHeader(name string, value string) {
	resp.w.Header().Add(name, value)
}

//SetHeader Response
func (resp *Response) SetHeader(name string, value string) {
	resp.w.Header().Set(name, value)
}

//FlushResponseHeader Flush Response Header
func (resp *Response) FlushResponseHeader(status int) (err error) {
	resp.w.WriteHeader(status)
	return err
}

func newHTTPResponse(srvlt *Servlet, reqst *Request, w http.ResponseWriter) (resp *Response) {
	resp = &Response{srvlt: srvlt, reqst: reqst, w: w}
	resp.iw, _ = goio.NewIORW(resp.w)
	return
}

//Print -. conveniant method that works the same as fmt.Fprint but writing to Response.w
func (resp *Response) Print(a ...interface{}) (n int, err error) {
	if resp.iw != nil {
		n, err = resp.iw.Print(a...)
	}
	return
}

//Println -. conveniant method that works the same as fmt.Fprint but writing to Response.w
func (resp *Response) Println(a ...interface{}) (n int, err error) {
	if resp.iw != nil {
		n, err = resp.iw.Println(a...)
	}
	return
}

func cleanupResponse(resp *Response) {
	if resp.iw != nil {
		resp.iw.Close()
		resp.iw = nil
	}
	if resp.reqst != nil {
		resp.reqst = nil
	}
	if resp.w != nil {
		resp.w = nil
	}
	if resp.srvlt != nil {
		resp.srvlt = nil
	}
}
