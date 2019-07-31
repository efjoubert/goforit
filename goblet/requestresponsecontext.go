package goblet

import (
	"net/http"
	"sync"
)

//ReqRespContext ReqRespContext
type ReqRespContext struct {
	w             http.ResponseWriter
	r             *http.Request
	srvlt         *Servlet
	reqst         *Request
	requestpath   string
	resp          *Response
	err           error
	doneReqstResp chan bool
	rlock         *sync.RWMutex
}

func queueServletRequestResponse(w http.ResponseWriter, r *http.Request, srvlt *Servlet, requestpath string) (err error) {
	if srvlt == nil {
		return
	}
	var reqRespContext = &ReqRespContext{w: w, r: r, srvlt: srvlt, doneReqstResp: make(chan bool, 1), requestpath: requestpath, rlock: &sync.RWMutex{}}
	defer func() {
		reqRespContext.rlock.RUnlock()
		cleanupRequestReponseContext(reqRespContext)
		reqRespContext = nil
	}()
	reqRespContext.rlock.RLock()
	queueCaptureReqstContext <- reqRespContext
	if <-reqRespContext.doneReqstResp {
		err = reqRespContext.err
	}
	return
}

func cleanupRequestReponseContext(reqRespCntxt *ReqRespContext) {
	if reqRespCntxt.doneReqstResp != nil {
		close(reqRespCntxt.doneReqstResp)
		reqRespCntxt.doneReqstResp = nil
	}
	if reqRespCntxt.err != nil {
		reqRespCntxt.err = nil
	}
	if reqRespCntxt.r != nil {
		reqRespCntxt.r = nil
	}
	if reqRespCntxt.w != nil {
		reqRespCntxt.w = nil
	}
	if reqRespCntxt.reqst != nil {
		cleanupRequest(reqRespCntxt.reqst)
		reqRespCntxt.reqst = nil
	}
	if reqRespCntxt.resp != nil {
		cleanupResponse(reqRespCntxt.resp)
		reqRespCntxt.resp = nil
	}
	if reqRespCntxt.rlock != nil {
		reqRespCntxt.rlock = nil
	}
}

var queueCaptureReqstContext chan *ReqRespContext

func initCaptureReqstContext() {
	if queueCaptureReqstContext == nil {
		queueCaptureReqstContext = make(chan *ReqRespContext)
		go func() {
			for {
				select {
				case reqstContext := <-queueCaptureReqstContext:
					go func() {
						reqstContext.reqst = newHTTPRequest(reqstContext.srvlt, reqstContext.r, reqstContext.requestpath)
						queueCaptureRespContext <- reqstContext
					}()
				}
			}
		}()
	}
}

var queueCaptureRespContext chan *ReqRespContext

func initCaptureRespContext() {
	if queueCaptureRespContext == nil {
		queueCaptureRespContext = make(chan *ReqRespContext)
		go func() {
			for {
				select {
				case respContext := <-queueCaptureRespContext:
					go func() {
						respContext.resp = newHTTPResponse(respContext.srvlt, respContext.reqst, respContext.w)
						executeServletRequestResponse(respContext)
					}()
				}
			}
		}()
	}
}

func init() {
	initCaptureReqstContext()
	initCaptureRespContext()
}

func executeServletRequestResponse(reqstrespcntxt *ReqRespContext) {
	defer func() { reqstrespcntxt.doneReqstResp <- true }()
	var srvlmethodhndl RequestResponseMethodHandle
	if srvlmethodhndl, _ = reqstrespcntxt.srvlt.reqstresphndls[reqstrespcntxt.reqst.method]; srvlmethodhndl != nil {
		reqstrespcntxt.err = srvlmethodhndl(reqstrespcntxt.srvlt, reqstrespcntxt.reqst, reqstrespcntxt.resp)
		srvlmethodhndl = nil
	}
}
