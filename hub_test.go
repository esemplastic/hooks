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

func TestPendingNotifiers(t *testing.T) {
	// first RegisterHook and after declare the Notify,
	// it should be succeed.
	var (
		state           = "hi"
		expectedPayload = "esemplastic"
		hub             = NewHub()
	)

	// notify first, the .Notify
	// will check if a registered hook is exists
	// if not it will add it as a "pending notifier"
	// and will try to execute that notifier
	// against each of the .RegisterHook's name argument.
	// If found then it should remove the pending.
	hub.Notify(state, expectedPayload)

	// it should be added because we don't have a registered hook with that name, yet.
	if expected, got := 1, len(hub.pendingNotifiers); expected != got {
		t.Fatalf("expected pending notifiers len to be %d but got %d", expected, got)
	}

	// register the notify now, it should be call the pending notifier.
	hub.RegisterHook(state, func(username string) {
		if username != expectedPayload {
			t.Fatalf("expected payload to be '%s' but got '%s'", expectedPayload, username)
		}
	})

	// it should be removed now.
	if expected, got := 0, len(hub.pendingNotifiers); expected != got {
		t.Fatalf("expected pending notifiers len to be %d but got %d", expected, got)
	}
}
