package hooks

import (
	"reflect"
	"sort"
	"sync"
)

type Hooks []*Hook
type HooksMap map[string]Hooks

type Registry interface {
	Register(name string, callback interface{}) *Hook
}

type Notifier interface {
	Run(name string, payloads ...interface{})
}

type pendingEntry struct {
	name     string
	payloads []interface{}
}

type Hub struct {
	logger         Logger
	mu             sync.RWMutex // locks for the hooks
	hooks          HooksMap
	pendingRunners []pendingEntry // if .Run before .Register, then the caller goes there and runs whenever it can.
}

var (
	_ Registry = &Hub{}
	_ Notifier = &Hub{}
)

func NewHub() *Hub {
	return &Hub{
		logger: DefaultLogger(),
		hooks:  make(HooksMap, 0),
	}
}

func (h *Hub) AttachLogger(logger Logger) {
	h.logger = logger
}

func (h *Hub) RegisterFunc(hookFunc interface{}, callback interface{}) *Hook {
	return h.Register(NameOfFunc(hookFunc), callback)
}

func (h *Hub) Register(name string, callback interface{}) *Hook {
	hook := newHook(h, name, callback)
	h.registerHook(hook)
	h.callPendingNotifiers(name)
	return hook
}

func (h *Hub) registerHook(hook *Hook) {
	h.mu.Lock()
	h.hooks[hook.Name] = append(h.hooks[hook.Name], hook)
	h.mu.Unlock()
}

// RemoveHooks removes all registered hooks sharing the same name.
//
// Returns true if the removal succeed.
func (h *Hub) RemoveHooks(name string) bool {
	return h.removeHooks(name)
}

func (h *Hub) removeHooks(name string) (ok bool) {
	if _, has := h.GetHooks(name); has {
		h.mu.Lock()
		delete(h.hooks, name)
		ok = true
		h.mu.Unlock()
	}
	return
}

// Remove removes a hook based on a function name and its callback.
//
// Same as Remove(NameOfFunc(fn), callback).
//
// Returns true if the removal succeed.
func (h *Hub) RemoveFunc(fn interface{}, callback interface{}) bool {
	return h.Remove(NameOfFunc(fn), callback)
}

// Remove removes a hook based on a hook name and its callback.
//
// Returns true if the removal succeed.
func (h *Hub) Remove(name string, callback interface{}) bool {
	return h.removeHook(name, callback)
}

func (h *Hub) removeHook(name string, callback interface{}) bool {
	if hooks, has := h.GetHooks(name); has {
		callbackPointer := reflect.ValueOf(callback).Pointer() // we could use the nameOfFunc too, but pointer is much safer.
		for i := range hooks {
			if hooks[i].Callback.Pointer() == callbackPointer {
				// remove that entry when found (we don't care about the order)
				hooks[i] = hooks[len(hooks)-1]
				h.mu.Lock()
				h.hooks[name] = hooks[:len(hooks)-1]
				h.mu.Unlock()
				return true
			}
		}
	}

	return false
}

func (h *Hub) RunFunc(hookFunc interface{}, payloads ...interface{}) {
	h.Run(NameOfFunc(hookFunc), payloads...)
}

func (h *Hub) Run(name string, payloads ...interface{}) {
	if hooks, has := h.GetHooks(name); has {
		for _, hook := range hooks {
			_, err := hook.Run(payloads...)
			if err != nil {
				h.logger(err.Error())
			}
		}
		return
	}
	h.addPendingRunner(name, payloads)
}

func (h *Hub) addPendingRunner(name string, payloads []interface{}) {
	entry := pendingEntry{
		name:     name,
		payloads: payloads,
	}
	h.pendingRunners = append(h.pendingRunners, entry)
}

func (h *Hub) callPendingNotifiers(registeredHookName string) {
	entries := h.pendingRunners
	for i, entry := range entries {
		if entry.name == registeredHookName {
			// remove that entry when found (we don't care about the order)
			entries[i] = entries[len(entries)-1]
			h.pendingRunners = entries[:len(entries)-1]
			// finally, do the Run now.
			h.Run(entry.name, entry.payloads...)
		}
	}
}

func (h *Hub) GetHooksFunc(fn interface{}) (Hooks, bool) {
	return h.GetHooks(NameOfFunc(fn))
}

func (h *Hub) GetHooks(name string) (Hooks, bool) {
	h.mu.RLock()
	hooks, has := h.hooks[name]
	h.mu.RUnlock()
	return hooks, has
}

// // this event will be notified to all Hooks, because the name of the func will be the same,
// // that is that we want here.
//
// forget it, it produces an overflow of stack if the user do an accident...
// func (h *Hub) RegisterHookChanged(changedFunction interface{}, callback func(hookChanged *Hook)) {
// 	h.RegisterFunc(changedFunction, callback)
// }

func (h *Hub) sortHooks(name string) {
	// per-group of hook maps, select and re-sort only these
	// that are inside the same hook map
	if hooks, has := h.GetHooks(name); has {
		h.mu.Lock()
		// sorts by the higher number of priority level
		sort.Slice(hooks, func(i, j int) bool {
			return hooks[i].Priority >= hooks[j].Priority
		})
		h.mu.Unlock()
	}
}
