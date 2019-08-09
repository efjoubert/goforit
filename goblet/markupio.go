package goblet

import (
	"strings"

	"github.com/efjoubert/goforit/goio"
)

//StartElem StartElem
func StartElem(out goio.Printer, elemName string, props ...interface{}) {
	out.Print("<", elemName)
	PrintProperties(out, "\"", props...)
	out.Print(">")
}

//PrintProperties PrintProperties
func PrintProperties(out goio.Printer, parenthasis string, props ...interface{}) {
	if len(props) > 0 {
		var ouputProp = func(pname string, pvalue string) {
			out.Print(" ", pname, "=", parenthasis, pvalue, parenthasis)
		}
		for _, prop := range props {
			if sprop, spropok := prop.(string); spropok && strings.Index(sprop, "=") > 0 {
				ouputProp(sprop[:strings.Index(sprop, "=")], sprop[strings.Index(sprop, "=")+1:])
			} else if mprop, mpropok := prop.(map[string]string); mpropok && len(mprop) > 0 {
				for pname, pvalue := range mprop {
					ouputProp(pname, pvalue)
				}
			} else if pprops, ppropsok := prop.([]string); ppropsok && len(props) > 0 {
				for _, prop := range pprops {
					if strings.Index(prop, "=") > 0 {
						ouputProp(prop[:strings.Index(prop, "=")], prop[strings.Index(prop, "=")+1:])
					}
				}
			}
		}
	}
}

//EndElem EndElem
func EndElem(out goio.Printer, elemName string) {
	out.Print("</", elemName, ">")
}

//SingleElem SingleElem
func SingleElem(out goio.Printer, elemName string, props ...interface{}) {
	out.Print("<", elemName)
	PrintProperties(out, "\"", props...)
	out.Print("/>")
}

//Elem Elem
func Elem(out goio.Printer, elemName string, a ...interface{}) {
	var props, funcs = SplitPropsIntoFuncAndProps(a...)
	StartElem(out, elemName, props...)
	for _, f := range funcs {
		f()
	}
	EndElem(out, elemName)
}

//SplitPropsIntoFuncAndProps SplitPropsIntoFuncAndProps
func SplitPropsIntoFuncAndProps(propstosplit ...interface{}) (props []interface{}, funcs []func()) {
	if len(propstosplit) > 0 {
		for _, propToTest := range propstosplit {
			if f, ok := propToTest.(func()); ok {
				if funcs == nil {
					funcs = []func(){}
				}
				funcs = append(funcs, f)
			} else {
				if props == nil {
					props = []interface{}{}
				}
				props = append(props, propToTest)
			}
		}
	}
	return props, funcs
}
