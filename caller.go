package hooks

import (
	"reflect"
)

type Caller struct {
	source Source
	action reflect.Value
}

func NewCaller(fn interface{}) Caller {
	val := reflect.ValueOf(fn)
	if !isFunc(val.Type()) {
		panic("action caller should be a function")
	}
	source := ReadSource(val.Pointer())

	return Caller{
		source: source,
		action: val,
	}
}
