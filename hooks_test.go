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

	hub.RegisterHook(state1, func(p PayloadWrapper) {
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

	hub.RegisterHook(state2, func(p PayloadWrapper) {
		if expected, got := reflect.TypeOf(expectedPayloads).Kind(), p.First().Kind(); expected != got {
			t.Fatalf("expected to receive kind of: %#v but got: %#v", expected, got)
		}
	})

	hub.RegisterHook(state3, func(p PayloadWrapper) {

		p.Iterate(func(index int, payload SinglePayloadWrapper) {

			if expected := len(expectedPayloads); expected-1 < index {
				t.Fatalf("[%d] - exceed number of expected payloads. Expected maximum len: %d", index, expected)
			}

			if expected, got := expectedPayloads[index], payload.String(); expected != got {
				t.Fatalf("[%d] - expected payload string to be: '%s' but got: '%s'", index, expected, got)
			}
		})
	})

	hub.Notify(state2, []Payload{expectedPayloads})
	hub.Notify(state3, Payload("payload1"), Payload("payload2"), Payload("payload3"), Payload("payload4"))
	// state3 note: This doesn't works []Payload{expectedPayloads}...)
}
