package hooks

import (
	"reflect"

	"github.com/satori/go.uuid"
)

type Hook struct {
	// ID will be likely to be used at the future, internally, in order to transfer hooks from other machine to another
	// (yes I have plans to make it net compatible.)
	ID       string
	Owner    *Hub
	Name     string // Multiple hooks can have the same name, is the event.
	Source   Source
	Callback reflect.Value // any type of function
	// if true then the callback will run in goroutine. Per-Name in the HooksMap, but it will run at goroutine so
	// it's a "global" field too, although user can just write go .Notify(...) for global usage among many hook names.
	Async bool
	// if remains zero then order matters on execution,
	// they are defaulted to the "IDLE" which doesn't matters if you don't care,
	// it has nothing to do with performance, is a matter of order.
	// each group of hooks has its own group, so the priority is per Name in the HooksMap
	Priority Priority
	// hiher number is the first.
	// optional descriptionist fields
	Description string
}

func newHook(owner *Hub, name string, callback interface{}) *Hook {
	val := reflect.ValueOf(callback)

	if typ := val.Type(); !isFunc(typ) {
		panic("callback should be a function but got: " + typ.String())
	}

	source := ReadSource(val.Pointer())
	id := uuid.NewV4()
	return &Hook{
		Owner:       owner,
		ID:          id.String(),
		Name:        name,
		Source:      source,
		Callback:    val,
		Async:       false,
		Priority:    Idle,
		Description: "",
	}
}

func (h *Hook) SetAsync(runInGoRoutine bool) *Hook {
	h.Async = runInGoRoutine
	return h
}

func (h *Hook) SetPriority(priority Priority) *Hook {
	oldP := h.Priority

	h.Priority = priority

	if oldP != priority {
		h.Owner.sortHooks(h.Name)
	}

	return h
}

func (h *Hook) PrioritizeAboveOf(otherHook *Hook) {
	h.SetPriority(otherHook.Priority + 1) // +1 is enough, we work with integers, so no need to check the priority "levels".
}
