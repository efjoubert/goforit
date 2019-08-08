package goreflect

import (
	"fmt"
	"reflect"
)

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
	return
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
