package hooks

import (
	"fmt"
	"time"
)

func ExampleHub_Register() {
	var (
		myHub   = NewHub()
		SayHook = "SAY_HOOK"
		send    = func(message string) {
			myHub.Run(SayHook, message)
		}
	)
	// Register two listeners on the same hook "SAY_HOOK".
	// Order of registrations order until one of the hook listener's
	// Priority changed, as we see below.

	myHub.Register(SayHook, func(message string) {
		fmt.Println("Processing the incoming message: " + message)
	}).SetPriority(Idle) // default priority

	myHub.Register(SayHook, func(message string) {
		fmt.Println("Incoming message: " + message)
	}).SetPriority(Realtime) // highest priority

	var messages = []string{
		"message1",
		"message2",
	}

	for _, msg := range messages {
		send(msg)
	}

	// Output:
	// Incoming message: message1
	// Processing the incoming message: message1
	// Incoming message: message2
	// Processing the incoming message: message2
}
func ExampleHub_Run() {
	///TODO: fill this example with hook funcs.
	fmt.Println("TODO")

	// Output:
	// TODO
}
func ExampleHub_Run_second() {
	hub := NewHub()

	hub.Register("hook1", func() {
		time.Sleep(1 * time.Second)
		fmt.Println("hey from hook1")
	})

	hub.Register("hook2", func() {
		fmt.Println("hey from hook2")
	})

	go hub.Run("hook1")
	hub.Run("hook2")

	// wait for hook1, it's too long, temporary we don't have a way to wait,but it's todo.
	time.Sleep(2 * time.Second)

	// hook2 should be printed first.
	// because the hooks1' time.Sleep runs in goroutine.
	//
	// At the future we should develop it in order to be able to set
	// the hook2 in Async state in order to be executed with other 'async'
	// hooks inside different goroutines, or inside a group of goroutines, we will see.

	// Also:
	// Be able to .Wait for specific callbacks or a group of them or globally.

	// Output:
	// hey from hook2
	// hey from hook1
}
