package goblet

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"../goio"
)

//Widget Widget
type Widget interface {
	goio.Printer
	CleanupWidget()
	WidgetMarkupHandle(...string) func() io.Reader
	DefaultMarkupHandle() func() io.Reader
	CallFunc(funcname string, a ...interface{}) (err error)
}

//BaseWidget BaseWidget
type BaseWidget struct {
	goio.Printer
	wdgtbrkr         *WidgetBroker
	dynwgdtfuncs     map[string]func(Widget)
	dynwgdtprmfuncs  map[string]func(Widget, ...interface{})
	dynfuncs         map[string]func()
	dynprmfuncs      map[string]func(...interface{})
	lckFunc          *sync.RWMutex
	dynwdgtmrkphndls map[string]func() io.Reader
	wdgtmrkphndl     func() io.Reader
}

func (bsewdgt *BaseWidget) assignWidgetFunc(funcname string, funcimpl interface{}) (found bool) {
	funcname = strings.ToUpper(funcname)
	var wdgtmrkphndl, wdgtmrkphndlok = funcimpl.(func() io.Reader)
	if wdgtmrkphndlok {
		var _, dynwdgtmrkphndlsok = bsewdgt.dynwdgtmrkphndls[funcname]
		if !dynwdgtmrkphndlsok {
			if _, found = bsewdgt.dynwdgtmrkphndls[funcname]; !found {
				bsewdgt.dynwdgtmrkphndls[funcname] = wdgtmrkphndl
				found = true
			}
		}
	} else {
		var _, dynwgdtprmfuncsok = bsewdgt.dynwgdtprmfuncs[funcname]
		if dynwgdtprmfuncsok {
			return
		}
		var _, dynwgdtfuncsok = bsewdgt.dynwgdtfuncs[funcname]
		if dynwgdtfuncsok {
			return
		}
		var _, dynprmfuncok = bsewdgt.dynprmfuncs[funcname]
		if dynprmfuncok {
			return
		}
		var _, dynfuncok = bsewdgt.dynfuncs[funcname]
		if dynfuncok {
			return
		}

		var paramswdgtfunc, paramswdgtfuncok = funcimpl.(func(Widget, ...interface{}))
		var nonparamswdgtfunc, nonparamswdgtfuncok = funcimpl.(func(Widget))
		var paramsfunc, paramsfuncok = funcimpl.(func(...interface{}))
		var nonparamfunc, nonparamfuncok = funcimpl.(func())
		if paramsfuncok && !nonparamfuncok && !paramswdgtfuncok && !nonparamswdgtfuncok {
			bsewdgt.dynprmfuncs[funcname] = paramsfunc
			found = true
		} else if !paramsfuncok && nonparamfuncok && !paramswdgtfuncok && !nonparamswdgtfuncok {
			bsewdgt.dynfuncs[funcname] = nonparamfunc
			found = true
		} else if !paramsfuncok && !nonparamfuncok && paramswdgtfuncok && !nonparamswdgtfuncok {
			bsewdgt.dynwgdtprmfuncs[funcname] = paramswdgtfunc
			found = true
		} else if !paramsfuncok && !nonparamfuncok && !paramswdgtfuncok && nonparamswdgtfuncok {
			bsewdgt.dynwgdtfuncs[funcname] = nonparamswdgtfunc
			found = true
		}
	}
	return
}

func invokeWidgetFunc(bsewdgt *BaseWidget, funcname string, a ...interface{}) {
	funcname = strings.ToUpper(funcname)
	var _, dynwgdtprmfuncsok = bsewdgt.dynwgdtprmfuncs[funcname]
	var _, dynwgdtfuncsok = bsewdgt.dynwgdtfuncs[funcname]
	var _, dynprmfuncok = bsewdgt.dynprmfuncs[funcname]
	var _, dynfuncok = bsewdgt.dynfuncs[funcname]
	if dynfuncok {
		var funcimpl = bsewdgt.dynfuncs[funcname]
		if funcimpl != nil {
			funcimpl()
		}
	} else if dynprmfuncok {
		var funcimpl = bsewdgt.dynprmfuncs[funcname]
		if funcimpl != nil {
			funcimpl(a...)
		}
	} else if dynwgdtprmfuncsok {
		var funcimpl = bsewdgt.dynwgdtprmfuncs[funcname]
		if funcimpl != nil {
			funcimpl(bsewdgt, a...)
		}
	} else if dynwgdtfuncsok {
		var funcimpl = bsewdgt.dynwgdtfuncs[funcname]
		if funcimpl != nil {
			funcimpl(bsewdgt)
		}
	}
}

//MapWidgetFunction MapWidgetFunction
func (bsewdgt *BaseWidget) MapWidgetFunction(funcname string, funcimpl interface{}, a ...interface{}) {
	if funcname != "" && funcimpl != nil {
		bsewdgt.lckFunc.RLock()
		defer bsewdgt.lckFunc.RUnlock()
		if !bsewdgt.assignWidgetFunc(funcname, funcimpl) {
			return
		}
	} else {
		return
	}

	for len(a) > 0 && len(a)%2 == 0 {
		if fncnme, fncnmeok := a[0].(string); fncnmeok {
			if fncnme != "" && a[1] != nil {
				if bsewdgt.assignWidgetFunc(fncnme, a[1]) {
					a = a[2:]
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
	}
}

//CallFunc CallFunc
func (bsewdgt *BaseWidget) CallFunc(funcname string, a ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
		}
	}()
	invokeWidgetFunc(bsewdgt, funcname, a...)
	return
}

//Print Print
func (bsewdgt *BaseWidget) Print(a ...interface{}) (int, error) {
	return bsewdgt.wdgtbrkr.Print(a...)
}

//Println Println
func (bsewdgt *BaseWidget) Println(a ...interface{}) (int, error) {
	return bsewdgt.wdgtbrkr.Println(a...)
}

//DefaultWidgetMarkupHandle DefaultWidgetMarkupHandle
func (bsewdgt *BaseWidget) DefaultWidgetMarkupHandle(wdgtmrkphndl func() io.Reader) {
	if bsewdgt.wdgtmrkphndl == nil {
		bsewdgt.wdgtmrkphndl = wdgtmrkphndl
	}
}

//DefaultMarkupHandle DefaultMarkupHandle
func (bsewdgt *BaseWidget) DefaultMarkupHandle() func() io.Reader {
	return bsewdgt.wdgtmrkphndl
}

//WidgetMarkupHandle WidgetMarkupHandle
func (bsewdgt *BaseWidget) WidgetMarkupHandle(funcname ...string) (wdgtmrkphndl func() io.Reader) {
	if len(funcname) == 1 && funcname[0] != "" {
		wdgtmrkphndl = bsewdgt.dynwdgtmrkphndls[strings.ToUpper(funcname[0])]
	} else {
		wdgtmrkphndl = bsewdgt.wdgtmrkphndl
	}
	return
}

//Broker Broker
func (bsewdgt *BaseWidget) Broker() *WidgetBroker {
	return bsewdgt.wdgtbrkr
}

//CleanupWidget CleanupWidget
func (bsewdgt *BaseWidget) CleanupWidget() {
	if bsewdgt.wdgtbrkr != nil {
		bsewdgt.wdgtbrkr = nil
	}

	if bsewdgt.dynfuncs != nil {
		if len(bsewdgt.dynfuncs) > 0 {
			var funcnames = []string{}
			for funcnme := range bsewdgt.dynfuncs {
				funcnames = append(funcnames, funcnme)
			}
			for len(funcnames) > 0 {
				bsewdgt.dynfuncs[funcnames[0]] = nil
				delete(bsewdgt.dynfuncs, funcnames[0])
				funcnames = funcnames[1:]
			}
			funcnames = nil
		}
		bsewdgt.dynfuncs = nil
	}
	if bsewdgt.dynprmfuncs != nil {
		if len(bsewdgt.dynprmfuncs) > 0 {
			var funcnames = []string{}
			for funcnme := range bsewdgt.dynprmfuncs {
				funcnames = append(funcnames, funcnme)
			}
			for len(funcnames) > 0 {
				bsewdgt.dynprmfuncs[funcnames[0]] = nil
				delete(bsewdgt.dynprmfuncs, funcnames[0])
				funcnames = funcnames[1:]
			}
			funcnames = nil
		}
		bsewdgt.dynprmfuncs = nil
	}
	if bsewdgt.dynwgdtfuncs != nil {
		if len(bsewdgt.dynwgdtfuncs) > 0 {
			var funcnames = []string{}
			for funcnme := range bsewdgt.dynwgdtfuncs {
				funcnames = append(funcnames, funcnme)
			}
			for len(funcnames) > 0 {
				bsewdgt.dynwgdtfuncs[funcnames[0]] = nil
				delete(bsewdgt.dynwgdtfuncs, funcnames[0])
				funcnames = funcnames[1:]
			}
			funcnames = nil
		}
		bsewdgt.dynwgdtfuncs = nil
	}
	if bsewdgt.dynwgdtprmfuncs != nil {
		if len(bsewdgt.dynwgdtprmfuncs) > 0 {
			var funcnames = []string{}
			for funcnme := range bsewdgt.dynwgdtprmfuncs {
				funcnames = append(funcnames, funcnme)
			}
			for len(funcnames) > 0 {
				bsewdgt.dynwgdtprmfuncs[funcnames[0]] = nil
				delete(bsewdgt.dynwgdtprmfuncs, funcnames[0])
				funcnames = funcnames[1:]
			}
			funcnames = nil
		}
		bsewdgt.dynwgdtprmfuncs = nil
	}

	if bsewdgt.dynwdgtmrkphndls != nil {
		if len(bsewdgt.dynwdgtmrkphndls) > 0 {
			var funcnames = []string{}
			for funcnme := range bsewdgt.dynwdgtmrkphndls {
				funcnames = append(funcnames, funcnme)
			}
			for len(funcnames) > 0 {
				bsewdgt.dynwdgtmrkphndls[funcnames[0]] = nil
				delete(bsewdgt.dynwdgtmrkphndls, funcnames[0])
				funcnames = funcnames[1:]
			}
			funcnames = nil
		}
		bsewdgt.dynwdgtmrkphndls = nil
	}
	if bsewdgt.wdgtmrkphndl != nil {
		bsewdgt.wdgtmrkphndl = nil
	}
}

//NewBaseWidget NewBaseWidget
func NewBaseWidget(wdgtbrkr *WidgetBroker) (bsewdgt *BaseWidget) {
	bsewdgt = &BaseWidget{wdgtbrkr: wdgtbrkr, dynfuncs: map[string]func(){}, dynprmfuncs: map[string]func(...interface{}){}, dynwgdtfuncs: map[string]func(Widget){}, dynwgdtprmfuncs: map[string]func(Widget, ...interface{}){}, lckFunc: &sync.RWMutex{}, dynwdgtmrkphndls: map[string]func() io.Reader{}}
	return
}

var widgetpaths = map[string]map[string]func(*WidgetBroker) Widget{}

//RegisterWidgetPath RegisterWidgetPath
func RegisterWidgetPath(path string, a ...interface{}) {
	var widgets = widgetpaths[path]
	if widgets == nil {
		widgets = map[string]func(*WidgetBroker) Widget{}
		widgetpaths[path] = widgets
	}

	for len(a) > 0 && len(a)%2 == 0 {
		if wdgtnme, wdgtnmeok := a[0].(string); wdgtnmeok {
			if wdgtinvhndl, wdgtinvhndlok := a[1].(func(*WidgetBroker) Widget); wdgtinvhndlok {
				if _, widgetsexist := widgets[wdgtnme]; !widgetsexist {
					widgets[wdgtnme] = wdgtinvhndl
				}
				a = a[2:]
			} else {
				break
			}
		} else {
			break
		}
	}
}

//RegisterWidget RegisterWidget
func RegisterWidget(path string, widgetname string, widgetinvokehndle func(*WidgetBroker) Widget, a ...interface{}) {
	if widgetname != "" && widgetinvokehndle != nil {

		var widgets = widgetpaths[path]
		if widgets != nil {
			var wdgtinvhandle = widgets[widgetname]
			if wdgtinvhandle == nil {
				wdgtinvhandle = widgetinvokehndle
				widgets[widgetname] = wdgtinvhandle
			}
		}

		for len(a) > 0 && len(a)%2 == 0 {
			if wdgtnme, wdgtnmeok := a[0].(string); wdgtnmeok {
				if wdgtinvhndl, wdgtinvhndlok := a[1].(func(*WidgetBroker) Widget); wdgtinvhndlok {
					if _, widgetsexist := widgets[wdgtnme]; !widgetsexist {
						widgets[wdgtnme] = wdgtinvhndl
					}
					a = a[2:]
				} else {
					break
				}
			} else {
				break
			}
		}
	}
}

//SearchWidgetInvokeHandle SearchWidgetInvokeHandle
func SearchWidgetInvokeHandle(wdgtpath string) (wdgtinvhandle func(*WidgetBroker) Widget, wdgtname string, actionpath string) {

	if wdgtpath != "" && strings.HasPrefix(wdgtpath, "/") && strings.LastIndex(wdgtpath, "/") < strings.LastIndex(wdgtpath, ".") {
		wdgtname = wdgtpath[strings.LastIndex(wdgtpath, "/")+1:]

		if wdgtname != "" {
			wdgtpath = wdgtpath[:len(wdgtpath)-len(wdgtname)]
			if strings.Index(wdgtpath, "/") < strings.LastIndex(wdgtpath, "/") && strings.HasSuffix(wdgtpath, "/") {
				wdgtpath = wdgtpath[:strings.LastIndex(wdgtpath, "/")]
			}
			if wdgtpath != "" {
				var wdgtext = filepath.Ext(wdgtname)
				if wdgtext != "" {
					wdgtname = wdgtname[:len(wdgtname)-len(wdgtext)]
					var wdgtactions = strings.Split(wdgtname, "-")
					if len(wdgtactions) > 1 {
						wdgtname = wdgtname[:len(wdgtname)-len(wdgtactions[len(wdgtactions)-1])-1]
						actionpath = wdgtactions[len(wdgtactions)-1]
					}
					wdgtactions = nil
				}

				if wdgts := widgetpaths[wdgtpath]; wdgts != nil {
					wdgtinvhandle = wdgts[wdgtname+wdgtext]
				}
			}
		}
	}

	return
}

//InvokeWidgetByHandle InvokeWidgetByHandle
func InvokeWidgetByHandle(wdgtinvhandle func(*WidgetBroker) Widget, wdgtbrkr *WidgetBroker) (wdgt Widget) {
	if wdgtinvhandle != nil && wdgtbrkr != nil {
		wdgt = wdgtinvhandle(wdgtbrkr)
	}
	return wdgt
}
