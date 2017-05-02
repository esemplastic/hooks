package hooks

import (
	"fmt"
	"time"
)

func ExampleHookAsync() {
	hub := NewHub()

	hub.RegisterHook("hook1", func() {
		time.Sleep(2 * time.Second)
		fmt.Println("hey from hook1")
	}).SetAsync(true)

	hub.RegisterHook("hook2", func() {
		fmt.Println("hey from hook2")
	})

	hub.Notify("hook1")
	hub.Notify("hook2")

	// wait for hook1, it's too long, temporary we don't have a way to wait,but it's todo.
	time.Sleep(4 * time.Second)

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
