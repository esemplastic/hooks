package hooks

import (
	"reflect"
)

type Hook struct {
	Name     string
	Source   Source
	Callback reflect.Value // any type of function
	Async    bool          // callback will run in goroutine if true
}

func NewHook(name string, callback interface{}) *Hook {
	val := reflect.ValueOf(callback)
	if !isFunc(val.Type()) {
		panic("callback should be a function")
	}
	source := ReadSource(val.Pointer())
	return &Hook{
		Name:     name,
		Source:   source,
		Callback: val,
		Async:    false,
	}
}

func (h *Hook) SetAsync(runInGoRoutine bool) *Hook {
	h.Async = runInGoRoutine
	return h
}
