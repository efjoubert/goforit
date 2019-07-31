package goblet

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"../goio"
)

//RequestResponseMethodHandle RequestResponseMethodHandle
type RequestResponseMethodHandle = func(*Servlet, *Request, *Response) error

//Servlet Servlet
type Servlet struct {
	srvltcntxt       *ServletContext
	reqstresphndls   map[string]RequestResponseMethodHandle
	srvlabsolutepath string
	srvltpath        string
}

//NewServlet NewServlet
func NewServlet(srvltcntxt *ServletContext, a ...interface{}) (srvlt *Servlet) {
	var reqstresphndls = map[string]RequestResponseMethodHandle{}
	for len(a) > 0 && len(a)%2 == 0 {
		if s, sok := a[0].(string); sok {
			if hndl, hndlok := a[1].(RequestResponseMethodHandle); hndlok {
				reqstresphndls[s] = hndl
			} else {
				break
			}
		}

		if len(a) >= 2 {
			a = a[2:]
		}
	}
	if len(reqstresphndls) > 0 {
		srvlt = &Servlet{srvltcntxt: srvltcntxt, reqstresphndls: reqstresphndls}
	}
	return
}

func servletSERVE(servlet *Servlet, request *Request, response *Response) (err error) {
	var wdgtbrkr = newWidgetBroker(response)
	var done = make(chan bool, 1)
	defer func() {
		done <- true
		close(done)
		wdgtbrkr.cleanupWidgetBroker()
		wdgtbrkr = nil
	}()
	wdgtbrkr.AppendNextRequest(request.url)
	var alreadyFoundOnce = false
	var foundActive = false
	var wnotify <-chan bool

	for len(wdgtbrkr.nextRequests) > 0 {
		var nextURL = wdgtbrkr.nextURL()
		var readerhndl = RegisteredEmbededReader(nextURL)
		if readerhndl == nil {
			readerhndl = wdgtbrkr.physicalResourceReaderHandle(nextURL)
		}
		if alreadyFoundOnce {
			if foundActive {
				wdgtbrkr.Println()
			}
		} else {
			alreadyFoundOnce = true
			var mimeType, activeExt = searchMimeTypeByExt(filepath.Ext(nextURL), "text/plain", strings.Split(activeExtensions, ",")...)
			foundActive = activeExt
			response.AddHeader("CONTENT-TYPE", mimeType)
		}

		var wdgtinvhndl, wdgtname, actionpath = SearchWidgetInvokeHandle(nextURL)

		var wdgt Widget
		if wdgtinvhndl != nil {
			wdgt = wdgtinvhndl(wdgtbrkr)
		} else {
			wdgt = NewBaseWidget(wdgtbrkr)
		}

		if wdgt != nil {
			if readerhndl == nil {
				if actionpath != "" {
					readerhndl = wdgt.WidgetMarkupHandle(actionpath)
				} else {
					readerhndl = wdgt.DefaultMarkupHandle()
				}
			}
		}

		func() {
			defer func() {
				if wdgt != nil {
					wdgt.CleanupWidget()
					wdgt = nil
				}
			}()
			if actionpath == "" {
				var params = wdgtbrkr.Parameters()
				if params != nil {
					var wdgtcmdparams = map[string][]string{}
					var wdgtcmds = []string{}
					var wdgtCaptureCmd = func(wdgtcmd string, vals ...string) {
						if wdgtcmd != "" {
							if cmds, cmdsok := wdgtcmdparams[wdgtcmd]; cmdsok {
								if len(vals) > 0 {
									cmds = append(cmds, vals...)
									wdgtcmdparams[wdgtcmd] = cmds
								}
							} else {
								wdgtcmds = append(wdgtcmds, wdgtcmd)
								if len(vals) > 0 {
									wdgtcmdparams[wdgtcmd] = vals[:]
								} else {
									wdgtcmdparams[wdgtcmd] = nil
								}
							}
						}
					}
					for _, pname := range params.StandardKeys() {
						if strings.HasPrefix(pname, strings.ToUpper(wdgtname)+"-") {
							pname = pname[len(wdgtname+"-"):]
							if pname == strings.ToUpper("command") {
								for _, pnme := range params.Parameter(wdgtname + "-" + pname) {
									wdgtCaptureCmd(strings.ToUpper(pnme), params.Parameter(pnme)...)
								}
							} else {
								wdgtCaptureCmd(pname, params.Parameter(wdgtname+"-"+pname)...)
							}
						}
					}
					for len(wdgtcmds) > 0 {
						if readerhndl != nil {
							readerhndl = nil
						}
						var wdgcmdprms = wdgtcmdparams[wdgtcmds[0]]
						if len(wdgcmdprms) > 0 {
							var a = make([]interface{}, len(wdgcmdprms))
							for n := range wdgcmdprms {
								a[n] = wdgcmdprms[n]
							}
							wdgt.CallFunc(wdgtcmds[0], a...)
							wdgtcmdparams[wdgtcmds[0]] = nil
						} else {
							wdgt.CallFunc(wdgtcmds[0])
						}
						delete(wdgtcmdparams, wdgtcmds[0])
						wdgcmdprms = nil
						wdgtcmds = wdgtcmds[1:]
					}
					wdgtcmdparams = nil
					wdgtcmds = nil
				}
			}
			if readerhndl != nil {
				var atvr, atverr = serverActiveContent(readerhndl(), foundActive, "_widget", wdgt, "_out", func() goio.Printer {
					return wdgtbrkr
				})

				if wnotify == nil && foundActive {
					wnotify = response.w.(http.CloseNotifier).CloseNotify()
					go func() {
						for {
							select {
							case ntfy := <-wnotify:
								if ntfy {
									if atvrf, atvrfok := atvr.(*ActiveReader); atvrfok {
										atvrf.interuptReader()
										return
									}
								}
							case dn := <-done:
								if dn {
									return
								}
							}
						}
					}()
				}

				_, err = wdgtbrkr.Print(atvr)
				if err == nil && atverr != nil {
					err = atverr
				}
				if atvrref, atvrok := atvr.(io.ReadCloser); atvrok {
					atvrref.Close()
					atvrref = nil
				}
				atvr = nil
				readerhndl = nil
			}
		}()
	}
	return
}

const activeExtensions string = ".js,.json,.css,.xml,.html,.htm,.svg"

//ServletGET ServletGET
func ServletGET(servlet *Servlet, request *Request, response *Response) (err error) {
	return servletSERVE(servlet, request, response)
}

//ServletPOST ServletPOST
func ServletPOST(servlet *Servlet, request *Request, response *Response) (err error) {
	return servletSERVE(servlet, request, response)
}

var embeddedreaders map[string]func() io.Reader

func init() {
	embeddedreaders = map[string]func() io.Reader{}
}

//RegisterEmbededReader RegisterEmbededReader
func RegisterEmbededReader(path string, readerhndl func() io.Reader) {
	if readerhndl != nil {
		if _, ok := embeddedreaders[path]; !ok {
			embeddedreaders[path] = readerhndl
		}
	}
}

//RegisterEmbededReaders RegisterEmbededReaders
func RegisterEmbededReaders(a ...interface{}) {
	for len(a) > 0 && len(a)%2 == 0 {
		if s, sok := a[0].(string); sok {
			if readerhndl, readerhndlok := a[1].(func() io.Reader); readerhndlok {
				if _, ok := embeddedreaders[s]; !ok {
					embeddedreaders[s] = readerhndl
				}
			} else {
				break
			}
		} else {
			break
		}
		if len(a) >= 2 {
			a = a[2:]
		}
	}
}

//RegisteredEmbededReader RegisteredEmbededReader
func RegisteredEmbededReader(path string) (readerhndl func() io.Reader) {
	readerhndl, _ = embeddedreaders[path]
	return readerhndl
}
