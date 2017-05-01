package hooks

import (
	"testing"
)

func TestRegisterHook(t *testing.T) {
	var state1 State = 10

	hub := NewHub()

	hub.RegisterHook(state1, func(p Payload) {

	})

	hub.Notify(state1, func() Payload {
		return nil
	})
}
