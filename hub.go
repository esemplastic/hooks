package hooks

import (
	"fmt"
	"sync"
)

type Hooks []*Hook
type HooksMap map[string]Hooks

type Registry interface {
	RegisterHook(name string, callback interface{}) *Hook
}

type Notifier interface {
	Notify(name string, payloads ...interface{})
}

type pendingEntry struct {
	name     string
	payloads []interface{}
}

type Hub struct {
	logger           Logger
	mu               sync.RWMutex // locks the hooks
	hooks            HooksMap
	pendingNotifiers []pendingEntry // slice of names(hooks' name)
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

func (h *Hub) RegisterHook(name string, callback interface{}) *Hook {
	hook := NewHook(name, callback)
	h.registerHook(hook)
	h.callPendingNotifiers(name)
	return hook
}

func (h *Hub) registerHook(hook *Hook) {
	h.mu.Lock()
	h.hooks[hook.Name] = append(h.hooks[hook.Name], hook)
	h.mu.Unlock()
}

func (h *Hub) Notify(name string, payloads ...interface{}) {
	if hooks, has := h.GetHooks(name); has {
		h.callHooks(hooks, name, payloads...)
		return
	}
	h.addPendingNotifier(name, payloads)
}

func (h *Hub) addPendingNotifier(name string, payloads []interface{}) {
	entry := pendingEntry{
		name:     name,
		payloads: payloads,
	}
	h.pendingNotifiers = append(h.pendingNotifiers, entry)
}

func (h *Hub) callPendingNotifiers(registeredHookName string) {
	entries := h.pendingNotifiers
	for i, entry := range entries {
		if entry.name == registeredHookName {
			// remove that entry when found (we don't care about the order)
			entries[i] = entries[len(entries)-1]
			h.pendingNotifiers = entries[:len(entries)-1]
			// finally, do the Notify now.
			h.Notify(entry.name, entry.payloads...)
		}
	}
}

func (h *Hub) GetHooks(name string) (Hooks, bool) {
	h.mu.RLock()
	hooks, has := h.hooks[name]
	h.mu.RUnlock()
	return hooks, has
}

func (h *Hub) callHooks(hooks Hooks, name string, arguments ...interface{}) {
	for _, hook := range hooks {
		if hook.Async {
			go h.callHook(hook, arguments...)
		} else {
			h.callHook(hook, arguments...)
		}
	}
}

func (h *Hub) callHook(hook *Hook, arguments ...interface{}) {
	_, err := execFunc(hook.Callback, arguments...)
	if err != nil {
		h.logger(fmt.Sprintf("error: %s\n callback metadata:\n  name: %s\n  file: %s\n  line: %d\n notification: '%s'",
			err.Error(), hook.Source.Name, hook.Source.File, hook.Source.Line, hook.Name))
	}
}
