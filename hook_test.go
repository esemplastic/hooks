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

	time.Sleep(2 * time.Second)

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
	}).SetPriority(High)

	hub.RegisterHook("myhook", func() {
		fmt.Println("forth")
	}).SetPriority(Normal)

	firstHook := hub.RegisterHook("myhook", func() {
		fmt.Println("first")
	}).SetPriority(Realtime) // or anything

	secondHook := hub.RegisterHook("myhook", func() {
		fmt.Println("second")
	}).SetPriority(Realtime)

	// even if Realtime, it can be priortized, remember we work with integers, Priority is just a type of int.
	firstHook.PrioritizeAboveOf(secondHook)

	hub.Notify("myhook")

	// Output:
	// first
	// second
	// third
	// forth
	// last
}

var messaging = NewHub()

func say(message string) {
	messaging.NotifyFunc(say, message)
}

func ExampleHookSource() {
	messaging.RegisterHookFunc(say, func(message string) {
		fmt.Printf("%s from %s via func %s\n", message, GetCurrentNotifier().Name, ReadSourceFunc(say).Name)
	})

	var messages = []string{
		"hello",
		"hi",
		"yo",
		"hola",
		"hey",
	}

	for _, s := range messages {
		say(s)
	}

	// Output:
	// hello from github.com/esemplastic/hooks.ExampleHookSource via func github.com/esemplastic/hooks.say
	// hi from github.com/esemplastic/hooks.ExampleHookSource via func github.com/esemplastic/hooks.say
	// yo from github.com/esemplastic/hooks.ExampleHookSource via func github.com/esemplastic/hooks.say
	// hola from github.com/esemplastic/hooks.ExampleHookSource via func github.com/esemplastic/hooks.say
	// hey from github.com/esemplastic/hooks.ExampleHookSource via func github.com/esemplastic/hooks.say
}
