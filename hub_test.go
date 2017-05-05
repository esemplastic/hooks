package hooks

import (
	"errors"
	"testing"
)

var (
	hub = NewHub()
)

func TestHubRegister(t *testing.T) {
	var (
		state1          = "state1"
		expectedPayload = errors.New("Error from state1")
	)

	hub.Register(state1, func(err error) {
		if expected, got := expectedPayload.Error(), err.Error(); expected != got {
			t.Fatalf("expected error message: '%s' but got: '%s'", expected, got)
		}
	})

	if expected, got := 1, len(hub.hooks); expected != got {
		t.Fatalf("expected hooks len to be %d but got %d", expected, got)
	}

	hub.Run(state1, expectedPayload)
}

func TestHubRun(t *testing.T) {
	var (
		state2           = "state2"
		state3           = "state3"
		expectedPayloads = []string{"payload1", "payload2", "payload3", "payload4"}
	)

	hub.Register(state2, func(payloads []string) {
		if expected, got := len(expectedPayloads), len(payloads); expected != got {
			t.Fatalf("expected payloads are different by the received, expected: %d but got: %d", expected, got)
		}
	})

	hub.Register(state3, func(payloads ...string) {
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

	hub.Run(state2, expectedPayloads)
	hub.Run(state3, "payload1", "payload2", "payload3", "payload4")
}

func TestHubPendingRunners(t *testing.T) {
	// first Register and after declare the Run,
	// it should be succeed.
	var (
		state           = "hi"
		expectedPayload = "esemplastic"
		hub             = NewHub()
	)

	// notify first, the .Run
	// will check if a registered hook is exists
	// if not it will add it as a "pending notifier"
	// and will try to execute that notifier
	// against each of the .Register's name argument.
	// If found then it should remove the pending.
	hub.Run(state, expectedPayload)

	// it should be added because we don't have a registered hook with that name, yet.
	if expected, got := 1, len(hub.pendingRunners); expected != got {
		t.Fatalf("expected pending notifiers len to be %d but got %d", expected, got)
	}

	// register the notify now, it should be call the pending notifier.
	hub.Register(state, func(username string) {
		if username != expectedPayload {
			t.Fatalf("expected payload to be '%s' but got '%s'", expectedPayload, username)
		}
	})

	// it should be removed now.
	if expected, got := 0, len(hub.pendingRunners); expected != got {
		t.Fatalf("expected pending notifiers len to be %d but got %d", expected, got)
	}
}

var (
	removeHookHub = NewHub()
)

func myHook(message string) {
	removeHookHub.RunFunc(myHook, message)
}

func TestHubRemove(t *testing.T) {
	t.Parallel()

	var (
		// we 're testing the notify and register again before the remove hook.
		expectedMessage = "hello"
	)

	var callback = func(message string) {
		if message != expectedMessage {
			t.Fatalf("expected incoming message to be: '%s' but got: '%s' ", expectedMessage, message)
		}
	}

	removeHookHub.RegisterFunc(myHook, callback)
	myHook(expectedMessage)

	hooks, _ := removeHookHub.GetHooksFunc(myHook)
	if expected, got := 1, len(hooks); expected != got {
		t.Fatalf("expected hooks len to be %d but got %d", expected, got)
	}

	// remove and test the len again.

	removed := removeHookHub.RemoveFunc(myHook, callback)
	if removed {
		hooks, _ := removeHookHub.GetHooksFunc(myHook)
		if expected, got := 0, len(hooks); expected != got {
			t.Fatalf("removed is true but expected hooks len to be %d but got %d", expected, got)
		}
	} else {
		t.Fatalf("remove action failed")
	}

	// try to remove the already removed, it should be give us false isntead.
	removed = removeHookHub.RemoveFunc(myHook, callback)
	if removed {
		t.Fatalf("remove action should be failed because we already remove that hook!")
	}
}

func TestHubRemoveHooks(t *testing.T) {
	///TODO: or add it inside TestRemoveHook
}

var messaging = NewHub()

func n() {
	messaging.RunFunc(n)
}

func TestGetCurrentRunner(t *testing.T) {
	messaging.RegisterFunc(n, func() {
		cRunner := GetCurrentRunner().Name
		fnName := ReadSourceFunc(n).Name

		if expected, got := "github.com/esemplastic/hooks.TestGetCurrentRunner", cRunner; expected != got {
			t.Fatalf("expected current runner's name source to be: '%s' but got: '%s'", expected, got)
		}

		if expected, got := "github.com/esemplastic/hooks.n", fnName; expected != got {
			t.Fatalf("expected hook's source to be: '%s' but got: '%s'", expected, got)
		}
	})

	n()
}
