package hooks

import (
	"reflect"
	"runtime"
)

type Source struct {
	Name string
	File string
	Line int
}

func ReadSourceFunc(fn interface{}) Source {
	return ReadSource(reflect.ValueOf(fn).Pointer())
}

func ReadSource(pointerfn uintptr) Source {
	pcfn := runtime.FuncForPC(pointerfn)
	name := pcfn.Name()
	file, line := pcfn.FileLine(pointerfn)
	return Source{
		Name: name,
		File: file,
		Line: line,
	}
}
