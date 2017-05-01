package hooks

import (
	"errors"
	"reflect"
	"testing"
)

var (
	hub = NewHub()
)

func TestRegisterHook(t *testing.T) {
	var (
		state1          State = 10
		expectedPayload       = errors.New("Error from state1")
	)

	hub.RegisterHook(state1, func(p Payloads) {
		if expected, got := expectedPayload.Error(), p.First().Err().Error(); expected != got {
			t.Fatalf("expected error message: '%s' but got: '%s'", expected, got)
		}
	})

	if expected, got := 1, len(hub.hooks); expected != got {
		t.Fatalf("expected hooks len to be %d but got %d", expected, got)
	}

	if expected, got := 1, len(hub.states); expected != got {
		t.Fatalf("expected states len to be %d but got %d", expected, got)
	}

	hub.Notify(state1, expectedPayload)
}

func TestNotify(t *testing.T) {
	var (
		state2           State = 20
		state3           State = 30
		expectedPayloads       = []string{"payload1", "payload2", "payload3", "payload4"}
	)

	hub.RegisterHook(state2, func(p Payloads) {
		if expected, got := reflect.TypeOf(expectedPayloads).Kind(), p.First().Kind(); expected != got {
			t.Fatalf("expected to receive kind of: %#v but got: %#v", expected, got)
		}
	})

	hub.RegisterHook(state3, func(p Payloads) {
		if expected, got := len(expectedPayloads), len(p); expected != got {
			t.Fatalf("expected payloads are different by the received, expected: %d but got: %d", expected, got)
		}

		p.Iterate(func(index int, payload Payload) {

			if expected := len(expectedPayloads); expected-1 < index {
				t.Fatalf("[%d] - exceed number of expected payloads. Expected maximum len: %d", index, expected)
			}

			if expected, got := expectedPayloads[index], payload.String(); expected != got {
				t.Fatalf("[%d] - expected payload string to be: '%s' but got: '%s'", index, expected, got)
			}
		})
	})

	hub.Notify(state2, expectedPayloads)
	hub.Notify(state3, "payload1", "payload2", "payload3", "payload4")
}
