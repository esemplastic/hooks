package hooks

import (
	"fmt"
	"sync"
)

type Hooks map[string][]*Hook

type Hub struct {
	logger Logger
	mu     sync.RWMutex // locks the hooks
	hooks  Hooks
}

func NewHub() *Hub {
	return &Hub{
		logger: DefaultLogger(),
		hooks:  make(Hooks, 0),
	}
}

func (h *Hub) AttachLogger(logger Logger) {
	h.logger = logger
}

func (h *Hub) RegisterHook(name string, callback interface{}) *Hook {
	hook := NewHook(name, callback)
	h.registerHook(hook)
	return hook
}

func (h *Hub) registerHook(hook *Hook) {
	h.mu.Lock()
	h.hooks[hook.Name] = append(h.hooks[hook.Name], hook)
	h.mu.Unlock()
}

func (h *Hub) Notify(name string, payloads ...interface{}) {
	callback := h.getHookCallbackFrom(name, payloads...)
	callback()
}

func (h *Hub) getHookCallbackFrom(name string, arguments ...interface{}) func() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	// choose to return a function which is de-coupled from the map
	// fast as possible them in order to return the release the lock as fast as possible.

	for hookName, hooks := range h.hooks {
		if hookName != name {
			continue
		}

		return func() {
			for _, hook := range hooks {
				if hook.Async {
					go h.callHook(hook, arguments...)
				} else {
					h.callHook(hook, arguments...)
				}

			}
		}
	}

	return func() {}
}

func (h *Hub) callHook(hook *Hook, arguments ...interface{}) {
	_, err := execFunc(hook.Callback, arguments...)
	if err != nil {
		h.logger(fmt.Sprintf("error: %s\n callback metadata:\n  name: %s\n  file: %s\n  line: %d\n notification: '%s'",
			err.Error(), hook.Source.Name, hook.Source.File, hook.Source.Line, hook.Name))
	}
}
