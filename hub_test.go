package hooks

import (
	"errors"
	"testing"
)

var (
	hub = NewHub()
)

func TestRegisterHook(t *testing.T) {
	var (
		state1          = "state1"
		expectedPayload = errors.New("Error from state1")
	)

	hub.RegisterHook(state1, func(err error) {
		if expected, got := expectedPayload.Error(), err.Error(); expected != got {
			t.Fatalf("expected error message: '%s' but got: '%s'", expected, got)
		}
	})

	if expected, got := 1, len(hub.hooks); expected != got {
		t.Fatalf("expected hooks len to be %d but got %d", expected, got)
	}

	hub.Notify(state1, expectedPayload)
}

func TestNotify(t *testing.T) {
	var (
		state2           = "state2"
		state3           = "state3"
		expectedPayloads = []string{"payload1", "payload2", "payload3", "payload4"}
	)

	hub.RegisterHook(state2, func(payloads []string) {
		if expected, got := len(expectedPayloads), len(payloads); expected != got {
			t.Fatalf("expected payloads are different by the received, expected: %d but got: %d", expected, got)
		}
	})

	hub.RegisterHook(state3, func(payloads ...string) {
		// println("--> " + payloads[0)
		if expected, got := len(expectedPayloads), len(payloads); expected != got {
			t.Fatalf("expected payloads length to be %d but got: %d", expected, got)
		}
		for index, payload := range payloads {
			if expected := len(expectedPayloads); expected-1 < index {
				t.Fatalf("[%d] - exceed number of expected payloads. Expected maximum len: %d", index, expected)
			}

			if expected, got := expectedPayloads[index], payload; expected != got {
				t.Fatalf("[%d] - expected payload string to be: '%s' but got: '%s'", index, expected, got)
			}
		}
	})

	hub.Notify(state2, expectedPayloads)
	hub.Notify(state3, "payload1", "payload2", "payload3", "payload4")
}
