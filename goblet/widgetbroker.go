package goblet

import (
	"io"
	"os"
	"strings"

	"../godbc"
	"../goio"
)

//WidgetBroker WidgetBroker
type WidgetBroker struct {
	srvlt        *Servlet
	reqst        *Request
	resp         *Response
	nextRequests []string
	print        goio.Printer
}

func newWidgetBroker(resp *Response) (wdgtbrkr *WidgetBroker) {
	wdgtbrkr = &WidgetBroker{resp: resp, srvlt: resp.srvlt, reqst: resp.reqst, print: resp}
	return
}

//Parameters Parameters
func (wdgtbrkr *WidgetBroker) Parameters() (params *Parameters) {
	if wdgtbrkr.reqst != nil && wdgtbrkr.reqst.params != nil {
		params = wdgtbrkr.reqst.params
	}
	return
}

func (wdgtbrkr *WidgetBroker) nextURL() (url string) {
	if len(wdgtbrkr.nextRequests) > 0 {
		url = wdgtbrkr.nextRequests[0]
		wdgtbrkr.nextRequests = wdgtbrkr.nextRequests[1:]
	} else {
		url = ""
	}
	return
}

//Query Query
func (wdgtbrkr *WidgetBroker) Query(dbalias string, query string) (dbquery *godbc.DBQuery) {
	dbquery = godbc.DatabaseManager().Query(dbalias, query)
	return dbquery
}

//AbsolutePath AbsolutePath
func (wdgtbrkr *WidgetBroker) AbsolutePath() (abslutpath string) {
	abslutpath = wdgtbrkr.srvlt.srvlabsolutepath
	return
}

func (wdgtbrkr *WidgetBroker) physicalResourceReaderHandle(nextURL string) (readerhndl func() io.Reader) {
	var wdgtabsolutepath = wdgtbrkr.AbsolutePath()
	if strings.HasPrefix(nextURL, "/") {
		nextURL = nextURL[1:]
	}
	if fio, err := os.Stat(wdgtabsolutepath + nextURL); fio != nil && err == nil {
		readerhndl = func() io.Reader {
			var f, _ = os.Open(wdgtabsolutepath + nextURL)
			return f
		}
	}
	return
}

//AppendNextRequest AppendNextRequest
func (wdgtbrkr *WidgetBroker) AppendNextRequest(nextrequest ...string) {
	if len(nextrequest) > 0 {
		if wdgtbrkr.nextRequests == nil {
			wdgtbrkr.nextRequests = []string{}
		}
		var preppedNextRequests = []string{}
		for len(nextrequest) > 0 {
			var nextr = strings.TrimSpace(nextrequest[0])
			if strings.Index(nextr, "|") > 0 {
				var extendnextrequest = append(strings.Split(nextr, "|"), nextrequest[1:]...)
				nextrequest = extendnextrequest[:]
				extendnextrequest = nil
			} else {
				if nextr != "" {
					preppedNextRequests = append(preppedNextRequests, nextr)
				}
				nextrequest = nextrequest[1:]
			}
		}
		nextrequest = nil
		if len(preppedNextRequests) > 0 {
			wdgtbrkr.nextRequests = append(wdgtbrkr.nextRequests, preppedNextRequests...)
		}
		preppedNextRequests = nil
	}
}

//Request Request
func (wdgtbrkr *WidgetBroker) Request() *Request {
	return wdgtbrkr.reqst
}

//Response Response
func (wdgtbrkr *WidgetBroker) Response() *Response {
	return wdgtbrkr.resp
}

//Print Print
func (wdgtbrkr *WidgetBroker) Print(a ...interface{}) (n int, err error) {
	n, err = wdgtbrkr.print.Print(a...)
	return
}

//Println Println
func (wdgtbrkr *WidgetBroker) Println(a ...interface{}) (n int, err error) {
	n, err = wdgtbrkr.print.Println(a...)
	return
}

func (wdgtbrkr *WidgetBroker) cleanupWidgetBroker() {
	if wdgtbrkr.resp != nil {
		wdgtbrkr.resp = nil
	}
	if wdgtbrkr.reqst != nil {
		wdgtbrkr.reqst = nil
	}
	if wdgtbrkr.srvlt != nil {
		wdgtbrkr.srvlt = nil
	}
	if wdgtbrkr.print != nil {
		wdgtbrkr.print = nil
	}
}
