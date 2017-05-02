# Hooks

This project is not ready yet.

Why?

- Code is not commented at all.
- Documentation is missing.
- Export of functions that is not part of the public API.
- Readme is just a blank paper of work.
- Travis Integration is missing.
- The [Idea](#idea) --features is not yet implemented. 
 

## Idea 

- [x] The RegisterHook would be able to catch any parameters
	from the Notify. So inside RegisterHook will be the action.

- [x] The RegisterHook should be de-coupled from the action
	in order to be able to use that library when import cycle
	is not allowed and code refactor is not possible.

- [x] The RegisterHook should be able to be work as an
	"event listener" too, so it should be able to not work
	as an action too. That makes the Notify to be able to execute
	the action function and notify the listeners.

- [x] We can do both RegisterHook and
	Notify to accept an interface{} and transfer the
	function parameters or return values between RegisterHook
	and Notify with the reflect package help.
	But we will lose type safety (already losed, the library
	provides some helpers but it's not evaluate the
	action function and its parameters).
	So we need to:

- [x] Convert any function to a hook, the function can be
	unique with the help of runtime package
	(by getting the full name of the func in the source code).
	But if we do that we lose the de-coupling described on second paragraph.
	
	Possible solution:
	- Keep the hook state as uint8 and make a function
	which will convert the func full name to
	a unique number -- Or just set the sate or rename to id with type of string and we're done with the funcname --
	so callers can still call
	and notifirs still notify without knowning each other, when needed and when no needed then
	the user can use the functionality without 
	touching a lot of existing code. 

- [ ] At the future we should develop it in order to be able to set a hook in Async state
    in order to be executed with other 'async' hooks inside different goroutines,
    or inside a group of goroutines, we will see. 

	Also:
	Be able to .Wait for specific callbacks or a group of them or globally.

- [x] The Notify should be able to be registered, as well, before RegisterHook in order to be able to be used
	when init functions are being used or when the order does not matters -- pending notifiers, remove them when RegisterHook registers the correct.

## License

Unless otherwise noted, the source files are distributed
under the BSD-3 Clause License found in the [LICENSE file](LICENSE).