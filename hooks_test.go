package hooks

import (
	"errors"
	"testing"
)

func TestRegisterHook(t *testing.T) {
	var (
		state1          State = 10
		expectedPayload       = errors.New("Error from state1")
	)

	hub := NewHub()

	hub.RegisterHook(state1, func(p PayloadWrapper) {
		if expected, got := expectedPayload.Error(), p.First().Err().Error(); expected != got {
			t.Fatalf("Expected error message: '%s' but got: '%s'", expected, got)
		}
	})

	if expected, got := 1, len(hub.hooks); expected != got {
		t.Fatalf("Expected hooks len to be %d but got %d", expected, got)
	}

	if expected, got := 1, len(hub.states); expected != got {
		t.Fatalf("Expected states len to be %d but got %d", expected, got)
	}

	hub.Notify(state1, expectedPayload)
}
