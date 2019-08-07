package main

import (
	"fmt"
	"os"
	"reflect"

	godataagent ".."
	"github.com/efjoubert/goforit/goplatform"
)

type TestS struct {
}

func (tsts *TestS) Test(t string) {

}

var test *TestS

func main() {
	test = &TestS{}
	interprateStruct(sTruct(test))
	var srvs, _ = goplatform.NewService("godataagent", "", "",
		startGoDataAgent,
		runGoDataAgent,
		stopGoDataAgent)
	if srvs != nil {
		srvs.Execute(os.Args)
	}
}

type Struct struct {
	tpe      reflect.Type
	mthds    []*Method
	mthdsmap map[string]*Method
	mthdnms  []string
	flds     []*Field
	fldsmap  map[string]*Field
	fldnms   []string
}

func (strct *Struct) CallMethod(owner interface{}, mthdname string, prmsin ...interface{}) (prmsout []interface{}) {
	if len(strct.mthds) > 0 {
		if _, ok := strct.mthdsmap[mthdname]; ok {
			prmsout = strct.mthdsmap[mthdname].Call(prmsin...)
		}
	}
	return
}

type Field struct {
	rflctfld reflect.StructField
	strct    *Struct
}

type Parameter struct {
	mthd *Method
	tpe  reflect.Type
}

type INParameter struct {
	*Parameter
}

type OUTParameter struct {
	*Parameter
}

type Method struct {
	strct     *Struct
	rflctmthd reflect.Method
	prmsin    []*INParameter
	prmsout   []*OUTParameter
}

func (mthd *Method) Call(prmsin ...interface{}) (prmsout []interface{}) {
	var rflctvalsin []reflect.Value
	if len(prmsin) > 0 {
		rflctvalsin = make([]reflect.Value, len(prmsin))
		for n := range prmsin {
			var pval = prmsin[n]
			rflctvalsin[n] = reflect.ValueOf(pval)
		}
	} else {
		rflctvalsin = []reflect.Value{}
	}
	var rflctvalsout = mthd.rflctmthd.Func.Call(rflctvalsin)
	if len(rflctvalsout) > 0 {
		prmsout = make([]interface{}, len(rflctvalsout))
		for n := range rflctvalsout {
			prmsout[n] = rflctvalsout[n].Interface()
		}
	}
}

func sTruct(a interface{}) (strct *Struct) {
	strct = &Struct{tpe: reflect.TypeOf(a)}
	return
}

func interprateStruct(strct *Struct) {
	var ptr = strct.tpe
	if mthdnum := ptr.NumMethod(); mthdnum > 0 {
		if strct.mthds == nil {
			strct.mthds = make([]*Method, mthdnum)
			strct.mthdnms = make([]string, mthdnum)
			strct.mthdsmap = map[string]*Method{}
		}
		for mthdnum > 0 {
			mthdnum--

			var mthd = &Method{strct: strct}
			mthd.strct = strct
			mthd.rflctmthd = ptr.Method(mthdnum)
			interprateParameters(mthd)
			strct.mthds[mthdnum] = mthd
			strct.mthdnms[mthdnum] = mthd.rflctmthd.Name
			strct.mthdsmap[strct.mthdnms[mthdnum]] = mthd
		}
	}
	if ptr.Kind() == reflect.Struct {
		if fldnum := ptr.NumField(); fldnum > 0 {
			if strct.flds == nil {
				strct.flds = make([]*Field, fldnum)
				strct.fldnms = make([]string, fldnum)
				strct.fldsmap = map[string]*Field{}
			}
			for fldnum > 0 {
				fldnum--
				var fld = &Field{strct: strct, rflctfld: ptr.Field(fldnum)}
				strct.fldnms[fldnum] = fld.rflctfld.Name
				strct.flds[fldnum] = fld
				strct.fldsmap[strct.fldnms[fldnum]] = fld
			}
		}
	}
}

func interprateParameters(mthd *Method) {

	if mthd != nil {
		var inNum = mthd.rflctmthd.Type.NumIn()
		var inIajust = 0
		if mthd.strct != nil && mthd.rflctmthd.Type.In(0) == mthd.strct.tpe {
			inIajust = 1
			inNum--
		}

		if inNum > 0 {
			mthd.prmsin = make([]*INParameter, inNum)
		}

		for i := 0; i < inNum; i++ {
			var prmin = &INParameter{Parameter: &Parameter{mthd: mthd, tpe: mthd.rflctmthd.Type.In(i + inIajust)}}
			fmt.Println(prmin.tpe.Name())
			mthd.prmsin[i] = prmin
		}

		var outNum = mthd.rflctmthd.Type.NumOut()

		if outNum > 0 {
			mthd.prmsout = make([]*OUTParameter, outNum)
		}

		for i := 0; i < outNum; i++ {
			var prmout = &OUTParameter{Parameter: &Parameter{mthd: mthd, tpe: mthd.rflctmthd.Type.Out(i)}}
			mthd.prmsout[i] = prmout
		}
	}
}

var dtaAgntsMngr *godataagent.DataAgentsManager

func startGoDataAgent(svs *goplatform.Service, args ...string) {
	if dtaAgntsMngr == nil {
		dtaAgntsMngr = godataagent.NewDataAgentsManager()
	}
}

func runGoDataAgent(svs *goplatform.Service, args ...string) {
	if dtaAgntsMngr != nil {
		dtaAgntsMngr.Startup("", "", "", ":1111")
	}
}

func stopGoDataAgent(svs *goplatform.Service, args ...string) {
	if dtaAgntsMngr != nil {
		dtaAgntsMngr.Shutdown()
		dtaAgntsMngr = nil
	}
}
