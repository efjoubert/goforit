package goblet

import (
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/kardianos/osext"
)

//ServletContext ServletContext
type ServletContext struct {
	servlets     map[string]*Servlet
	physicalpath string
	srvltreqlock *sync.RWMutex
}

//RegisterServlet RegisterServlet
func (srvltcntxt *ServletContext) RegisterServlet(srvletpath string, a interface{}, amethods ...interface{}) (srvlt *Servlet) {
	srvltcntxt.srvltreqlock.RLock()
	defer srvltcntxt.srvltreqlock.RUnlock()

	if srvlt, _ = srvltcntxt.servlets[srvletpath]; srvlt == nil {
		if srvltref, srvlok := a.(*Servlet); srvlok {
			srvlt = srvltref
			srvltcntxt.servlets[srvletpath] = srvlt
			srvlt.srvltpath = srvletpath
			srvlt.srvlabsolutepath = srvltcntxt.physicalpath
			if strings.HasPrefix(srvlt.srvltpath, "/") {
				if strings.HasSuffix(srvlt.srvlabsolutepath, "/") {
					srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath[1:]
				} else {
					srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath
				}
			} else {
				if strings.HasSuffix(srvlt.srvlabsolutepath, "/") {
					srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath
				} else {
					srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + "/" + srvlt.srvltpath
				}
			}
		} else {
			if amethods == nil {
				amethods = []interface{}{"GET", ServletGET, "POST", ServletPOST}
			}
			if srvlt = NewServlet(srvltcntxt, amethods...); srvlt != nil {
				srvltcntxt.servlets[srvletpath] = srvlt
				srvlt.srvltpath = srvletpath
				srvlt.srvlabsolutepath = srvltcntxt.physicalpath
				if strings.HasPrefix(srvlt.srvltpath, "/") {
					if strings.HasSuffix(srvlt.srvlabsolutepath, "/") {
						srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath[1:]
					} else {
						srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath
					}
				} else {
					if strings.HasSuffix(srvlt.srvlabsolutepath, "/") {
						srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath
					} else {
						srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + "/" + srvlt.srvltpath
					}
				}
			}
		}
	} else if a == nil {
		if amethods == nil {
			amethods = []interface{}{"GET", ServletGET, "POST", ServletPOST}
		}
		if srvlt = NewServlet(srvltcntxt, amethods...); srvlt != nil {
			srvltcntxt.servlets[srvletpath] = srvlt
			srvlt.srvltpath = srvletpath
			srvlt.srvlabsolutepath = srvltcntxt.physicalpath
			if strings.HasPrefix(srvlt.srvltpath, "/") {
				if strings.HasSuffix(srvlt.srvlabsolutepath, "/") {
					srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath[1:]
				} else {
					srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath
				}
			} else {
				if strings.HasSuffix(srvlt.srvlabsolutepath, "/") {
					srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + srvlt.srvltpath
				} else {
					srvlt.srvlabsolutepath = srvlt.srvlabsolutepath + "/" + srvlt.srvltpath
				}
			}
		}
	}
	return
}

func (srvltcntxt *ServletContext) servletHTTPRequest(w http.ResponseWriter, r *http.Request, srvlRequestPath string) {
	var srvlt *Servlet
	srvltcntxt.srvltreqlock.RLock()
	srvlt, srvlRequestPath = srvltcntxt.searchRegisteredServlet(srvlRequestPath)
	srvltcntxt.srvltreqlock.RUnlock()
	queueServletRequestResponse(w, r, srvlt, srvlRequestPath)
}

func (srvltcntxt *ServletContext) registeredServlet(srvltpath string) (srvlt *Servlet) {
	srvlt, _ = srvltcntxt.servlets[srvltpath]
	return
}

func (srvltcntxt *ServletContext) searchRegisteredServlet(srvltpath string) (srvlt *Servlet, cntxsrvltremainderpath string) {
	cntxsrvltremainderpath = "/"
	if !strings.HasPrefix(srvltpath, "/") {
		srvltpath = "/" + srvltpath
	}
	if strings.Index(srvltpath, "/") < strings.Index(srvltpath, ".") {
		cntxsrvltremainderpath = cntxsrvltremainderpath + srvltpath[strings.LastIndex(srvltpath[0:strings.Index(srvltpath, ".")], "/")+1:]
		srvltpath = srvltpath[:strings.LastIndex(srvltpath[0:strings.Index(srvltpath, ".")], "/")]
	}
	if srvltpath == "" {
		srvltpath = "/"
	}

	if srvltpath != "/" && strings.HasSuffix(srvltpath, "/") {
		srvltpath = srvltpath[:len(srvltpath)-1]
	}

	for srvltpath != "" {
		if _, srvltpathok := srvltcntxt.servlets[srvltpath]; srvltpathok {
			srvlt = srvltcntxt.servlets[srvltpath]
			break
		}

		if srvltpath != "/" && strings.LastIndex(srvltpath, "/") > -1 {
			cntxsrvltremainderpath = srvltpath[strings.LastIndex(srvltpath, "/"):] + cntxsrvltremainderpath
			srvltpath = srvltpath[:strings.LastIndex(srvltpath, "/")]
			if srvltpath == "" {
				if _, srvltpathok := srvltcntxt.servlets["/"]; srvltpathok {
					srvlt = srvltcntxt.servlets["/"]
					break
				}
			}
		} else {
			break
		}
	}
	return
}

var servletcontexts map[string]*ServletContext

//RegisterServletContextPath RegisterServletContextPath
func RegisterServletContextPath(contextpah string, physicalpath string, a ...interface{}) (srvltcntxt *ServletContext) {
	if _, cntxpathok := servletcontexts[contextpah]; cntxpathok {
		srvltcntxt = servletcontexts[contextpah]
	} else {
		if physicalpath == "" {
			execfolder, _ := osext.ExecutableFolder()
			physicalpath = strings.ReplaceAll(execfolder, "\\", "/")
			if !strings.HasSuffix(physicalpath, "/") {
				physicalpath = physicalpath + "/"
			}
		}
		srvltcntxt = &ServletContext{servlets: map[string]*Servlet{}, physicalpath: physicalpath, srvltreqlock: &sync.RWMutex{}}
		servletcontexts[contextpah] = srvltcntxt
	}

	if len(a) > 0 {
		for len(a) > 0 && len(a)&2 == 0 {
			if s, sok := a[0].(string); sok {
				if srvlt, srvltok := a[1].(*Servlet); srvltok {
					srvltcntxt.RegisterServlet(s, srvlt)
				} else {
					break
				}
			}

			if len(a) > 2 {
				a = a[2:]
			} else {
				a = nil
				break
			}
		}
	}

	return
}

func init() {
	if servletcontexts == nil {
		servletcontexts = map[string]*ServletContext{}
	}
	RegisterServletContextPath("/", "").RegisterServlet("/", nil)
}

//RegisteredServletContext RegisteredServletContext
func RegisteredServletContext(contextpath string) (srvltcntxt *ServletContext) {
	if _, cntxpathok := servletcontexts[contextpath]; cntxpathok {
		srvltcntxt = servletcontexts[contextpath]
	}
	return
}

func searchRegisteredServletContext(contextpath string) (srvltcntxt *ServletContext, cntxremainderPath string) {
	cntxremainderPath = "/"
	if !strings.HasPrefix(contextpath, "/") {
		contextpath = "/" + contextpath
	}
	if strings.Index(contextpath, "/") < strings.Index(contextpath, ".") {
		cntxremainderPath = cntxremainderPath + contextpath[strings.LastIndex(contextpath[0:strings.Index(contextpath, ".")], "/")+1:]
		contextpath = contextpath[:strings.LastIndex(contextpath[0:strings.Index(contextpath, ".")], "/")]
	}
	if contextpath == "" {
		contextpath = "/"
	}

	if contextpath != "/" && strings.HasSuffix(contextpath, "/") {
		contextpath = contextpath[:len(contextpath)-1]
	}

	for contextpath != "" {
		if _, cntxpathok := servletcontexts[contextpath]; cntxpathok {
			srvltcntxt = servletcontexts[contextpath]
			break
		}

		if contextpath != "/" && strings.LastIndex(contextpath, "/") > -1 {
			cntxremainderPath = contextpath[strings.LastIndex(contextpath, "/"):] + cntxremainderPath
			contextpath = contextpath[:strings.LastIndex(contextpath, "/")]
			if contextpath == "" {
				if _, cntxpathok := servletcontexts["/"]; cntxpathok {
					srvltcntxt = servletcontexts["/"]
					break
				}
			}
		} else {
			break
		}
	}

	return
}

//PerformServletRequest PerformServletRequest
func PerformServletRequest(w io.Writer, r io.Reader) {

}

//PerformHTTPServletRequest PerformHttpServletRequest
func PerformHTTPServletRequest(w http.ResponseWriter, r *http.Request) {
	var srvlcntx, cntxremainderPath = searchRegisteredServletContext(r.URL.Path)
	if srvlcntx != nil {
		srvlcntx.servletHTTPRequest(w, r, cntxremainderPath)
	}
}
