package hooks

import (
	"fmt"
	"time"
)

func ExampleHook_SetAsync() {
	hub := NewHub()

	hub.Register("myhook", func() {
		time.Sleep(1 * time.Second)
		fmt.Println("last")
	}).SetAsync(true)

	hub.Register("myhook", func() {
		fmt.Println("first")
	})

	hub.Run("myhook")

	time.Sleep(2 * time.Second)

	// Output:
	// first
	// last
}

func ExampleHook_SetPriority() {
	hub := NewHub()

	hub.Register("myhook", func() {
		fmt.Println("last")
	}).SetPriority(Idle) // defaults to Idle already.

	hub.Register("myhook", func() {
		fmt.Println("second")
	}).SetPriority(Normal)

	hub.Register("myhook", func() {
		fmt.Println("third")
	}).SetPriority(BelowNormal)

	hub.Register("myhook", func() {
		fmt.Println("first")
	}).SetPriority(High)

	hub.Run("myhook")

	// Output:
	// first
	// second
	// third
	// last
}

func ExampleHook_PrioritizeAboveOf() {
	hub := NewHub()

	hub.Register("myhook", func() {
		fmt.Println("last")
	}).SetPriority(Idle) // defaults to Idle already.

	hub.Register("myhook", func() {
		fmt.Println("third")
	}).SetPriority(High)

	hub.Register("myhook", func() {
		fmt.Println("forth")
	}).SetPriority(Normal)

	firstHook := hub.Register("myhook", func() {
		fmt.Println("first")
	}).SetPriority(Realtime) // or anything

	secondHook := hub.Register("myhook", func() {
		fmt.Println("second")
	}).SetPriority(Realtime)

	// even if Realtime, it can be priortized, remember we work with integers, Priority is just a type of int.
	firstHook.PrioritizeAboveOf(secondHook)

	hub.Run("myhook")

	// Output:
	// first
	// second
	// third
	// forth
	// last
}
