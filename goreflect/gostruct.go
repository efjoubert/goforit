package goreflect

import (
	"reflect"
)

type Struct struct {
	tpe      reflect.Type
	mthds    []*Method
	mthdsmap map[string]*Method
	mthdnms  []string
	flds     []*Field
	fldsmap  map[string]*Field
	fldnms   []string
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

func (strct *Struct) CallMethod(owner interface{}, mthdname string, prmsin ...interface{}) (prmsout []interface{}) {
	if len(strct.mthds) > 0 {
		if _, ok := strct.mthdsmap[mthdname]; ok {
			prmsout = strct.mthdsmap[mthdname].Call(prmsin...)
		}
	}
	return
}
