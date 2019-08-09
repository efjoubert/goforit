package goblet

import (
	"github.com/efjoubert/goforit/goio"
)

//StartActiveContent - to be used with webactions.js embedded resource
func StartActiveContent(out goio.Printer, contentName string) {
	out.Print("replace-content||" + contentName + "||")
}

//EndActiveContent - to be used with webactions.js embedded resource
func EndActiveContent(out goio.Printer, contentName string) {
	out.Print("||replace-content")
}

//ActiveContent - to be used with webactions.js embedded resource
func ActiveContent(out goio.Printer, contentName string, a ...interface{}) {
	var _, funcs = SplitPropsIntoFuncAndProps(a...)
	StartActiveContent(out, contentName)
	for _, f := range funcs {
		f()
	}
	EndActiveContent(out, contentName)
}

//StartActiveScript - to be used with webactions.js embedded resource
func StartActiveScript(out goio.Printer) {
	out.Print("script||")
}

//EndActiveScript - to be used with webactions.js embedded resource
func EndActiveScript(out goio.Printer) {
	out.Print("||script")
}

//ActiveScript - to be used with webactions.js embedded resource
func ActiveScript(out goio.Printer, a ...interface{}) {
	var _, funcs = SplitPropsIntoFuncAndProps(a...)
	StartActiveScript(out)
	for _, f := range funcs {
		f()
	}
	EndActiveScript(out)
}
