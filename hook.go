package hooks

import (
	"fmt"
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
	Callback reflect.Value
	// if remains zero then order matters on execution,
	// they are defaulted to the "IDLE" which doesn't matters if you don't care,
	// it has nothing to do with performance, is a matter of order.
	// each group of hooks has its own group, so the priority is per Name in the HooksMap.
	//
	// Read-only value, use SetPriority if you want to alt it.
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
		Priority:    Idle,
		Description: "",
	}
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

func (h *Hook) Run(payloads ...interface{}) ([]reflect.Value, error) {
	return h.run(payloads...)
}

func (h *Hook) RunAsync(payloads ...interface{}) *HookAsyncRunner {
	asyncRunner := new(HookAsyncRunner)
	go func() {
		returnValues, err := h.run(payloads...)
		asyncRunner.fireComplete(returnValues, err)
	}()

	return asyncRunner
}

type (
	HookAsyncResultListener func([]reflect.Value, error)
	HookAsyncRunner         struct {
		completeListeners   []HookAsyncResultListener
		pendingReturnValues []reflect.Value
		pendingErr          error
	}
)

func (runner *HookAsyncRunner) OnComplete(listeners ...HookAsyncResultListener) {
	runner.completeListeners = append(runner.completeListeners, listeners...)
	if runner.pendingReturnValues != nil || runner.pendingErr != nil {
		runner.fireComplete(runner.pendingReturnValues, runner.pendingErr)
		runner.pendingReturnValues = nil
		runner.pendingErr = nil
	}
}

func (runner *HookAsyncRunner) fireComplete(returnValues []reflect.Value, err error) {
	// if fire before OnComplete registered, save the results to fire them on
	// the first call of .OnComplete (can have multiple listeners)
	if len(runner.completeListeners) == 0 {
		runner.pendingReturnValues = returnValues
		runner.pendingErr = err
		return
	}

	for i := range runner.completeListeners {
		runner.completeListeners[i](returnValues, err)
	}
}

func (h *Hook) run(payloads ...interface{}) (returnValues []reflect.Value, err error) { // maybe return values are useless and confusing here.

	returnValues, err = execFunc(h.Callback, payloads...)
	if err != nil {
		err = fmt.Errorf("error: %s\n callback metadata:\n  name: %s\n  file: %s\n  line: %d\n notification: '%s'",
			err.Error(), h.Source.Name, h.Source.File, h.Source.Line, h.Name)
	}

	return
}

func (h *Hook) Use(preProcessor interface{}) *Hook { // func(callback interface{}, payloads ...interface{})) *Hook {

	// oldCallback := h.Callback
	// oldNext := h.preProcessor
	// newNext := preProcessor.(func(func(string), string))
	// next := func(message string) {
	// 	if oldNext.IsValid() {
	// 		_, err := execFunc(oldNext, message)
	// 		if err != nil {
	// 			println(err.Error())
	// 		}
	// 	}
	// 	// 	println("from next")
	// 	newNext(func(newMessage string) {
	// 		_, err := execFunc(oldCallback, newMessage)
	// 		if err != nil {
	// 			println(err.Error())
	// 		}
	// 	}, message)
	// }

	// h.preProcessor = reflect.ValueOf(next)

	return h
}
