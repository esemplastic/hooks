package hooks

import (
	"fmt"
	"time"
)

func ExampleHookGoNotify() {
	hub := NewHub()

	hub.RegisterHook("hook1", func() {
		time.Sleep(1 * time.Second)
		fmt.Println("hey from hook1")
	})

	hub.RegisterHook("hook2", func() {
		fmt.Println("hey from hook2")
	})

	go hub.Notify("hook1")
	hub.Notify("hook2")

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

func ExampleHookSetAsync() {
	hub := NewHub()

	hub.RegisterHook("myhook", func() {
		time.Sleep(1 * time.Second)
		fmt.Println("last")
	}).SetAsync(true)

	hub.RegisterHook("myhook", func() {
		fmt.Println("first")
	})

	hub.Notify("myhook")

	// wait for myhook to finish, it's too long, temporary we don't have a way to wait, it's todo.
	time.Sleep(2 * time.Second)

	// At the future we should develop it in order to be able to set
	// the hook2 in Async state in order to be executed with other 'async'
	// hooks inside different goroutines, or inside a group of goroutines, we will see.

	// Also:
	// Be able to .Wait for specific callbacks or a group of them or globally.

	// Output:
	// first
	// last
}

func ExampleHookSetPriority() {
	hub := NewHub()

	hub.RegisterHook("myhook", func() {
		fmt.Println("last")
	}).SetPriority(Idle) // defaults to Idle already.

	hub.RegisterHook("myhook", func() {
		fmt.Println("second")
	}).SetPriority(Normal)

	hub.RegisterHook("myhook", func() {
		fmt.Println("third")
	}).SetPriority(BelowNormal)

	hub.RegisterHook("myhook", func() {
		fmt.Println("first")
	}).SetPriority(High)

	hub.Notify("myhook")

	// Output:
	// first
	// second
	// third
	// last
}

func ExamplePrioritizeAboveOf() {
	hub := NewHub()

	hub.RegisterHook("myhook", func() {
		fmt.Println("last")
	}).SetPriority(Idle) // defaults to Idle already.

	hub.RegisterHook("myhook", func() {
		fmt.Println("third")
	}).SetPriority(Normal)

	hub.RegisterHook("myhook", func() {
		fmt.Println("forth")
	}).SetPriority(BelowNormal)

	firstHook := hub.RegisterHook("myhook", func() {
		fmt.Println("first")
	}).SetPriority(High) // or anything

	secondHook := hub.RegisterHook("myhook", func() {
		fmt.Println("second")
	}).SetPriority(High)

	// even if High, it can be priortized, remember we work with integers, Priority is just a type of int.
	firstHook.PrioritizeAboveOf(secondHook)

	hub.Notify("myhook")

	// Output:
	// first
	// second
	// third
	// forth
	// last
}
