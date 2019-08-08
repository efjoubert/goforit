package goreflect

import "reflect"

type Field struct {
	rflctfld reflect.StructField
	strct    *Struct
}
