package hooks

import (
	"github.com/satori/go.uuid"
	"reflect"
)

type Hook struct {
	// ID will be likely to be used at the future, internally, in order to transfer hooks from other machine to another
	// (yes I have plans to make it net compatible.)
	ID string

	Name     string // Multiple hooks can have the same name, is the event.
	Source   Source
	Callback reflect.Value // any type of function
	Async    bool          // callback will run in goroutine if true
}

func NewHook(name string, callback interface{}) *Hook {
	val := reflect.ValueOf(callback)

	if typ := val.Type(); !isFunc(typ) {
		panic("callback should be a function but got: " + typ.String())
	}

	source := ReadSource(val.Pointer())
	id := uuid.NewV4()
	return &Hook{
		ID:       id.String(),
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
