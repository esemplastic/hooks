package hooks

import (
	"sync"
)

type State uint8

type (
	Hook     func(payloads Payloads)
	HooksMap map[State][]Hook
)

type (
	Action     func() Payload
	Actions    []Action
	ActionsMap map[State]Actions
)

type Registrar interface {
	RegisterHook(State, Hook)
	RegisterHooks(HooksMap)
}

type Notifier interface {
	Notify(State, ...interface{})
	NotifyWithAction(state State, action Action)
	NotifyWithActionMany(actionsMap ActionsMap)
}

type Hub struct {
	mu     sync.RWMutex
	states []State
	hooks  HooksMap
}

var _ Registrar = &Hub{}
var _ Notifier = &Hub{}

func NewHub() *Hub {
	return &Hub{
		hooks: make(HooksMap, 0),
	}
}

func (h *Hub) Registrar() Registrar {
	return h
}

func (h *Hub) Notifier() Notifier {
	return h
}

func (h *Hub) addState(state State) {
	h.states = append(h.states, state)
}

func (h *Hub) RegisterHook(state State, hook Hook) {
	h.mu.Lock()
	defer h.mu.Unlock()

	found := false
	for i := range h.states {
		if h.states[i] == state {
			found = true
			break
		}
	}

	if !found {
		h.addState(state)
	}

	h.hooks[state] = append(h.hooks[state], hook)
}

func (h *Hub) RegisterHooks(hooks HooksMap) {
	for k, v := range hooks {
		for i := range v {
			h.RegisterHook(k, v[i])
		}
	}
}

func (h *Hub) Notify(state State, genericPayloads ...interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var allPayloads payloads

	for _, p := range genericPayloads {
		allPayloads = append(allPayloads, ReadPayload(p))
	}

	for k, v := range h.hooks {
		if k != state {
			continue
		}
		for i := range v {
			v[i](allPayloads)
		}
	}
}

func (h *Hub) NotifyWithAction(state State, action Action) {
	payload := action()
	h.Notify(state, payload)
}

func (h *Hub) NotifyWithActionMany(actionsMap ActionsMap) {
	for k, v := range actionsMap {
		for i := range v {
			h.NotifyWithAction(k, v[i])
		}
	}
}
