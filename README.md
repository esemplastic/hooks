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

- [x] Be able to catch raw function execution without removing the dynamic behavior too. 

	Explanation: In order to keep the statical typing feature of the language, 
	I must think a way to add notifiers and hooks with a custom function form. 

	The same function form which should be able to be splitted into: 
	1. name and callback for the notifier
	2. name and payloads for the registrar

	For example:
	- /myhub contains the form which in the same time calls the (new) NotifyFunc which
		should gets the func and converts that to a name and calls the dynamic .Notify: 
	- (myhub) func Add(item){ hub.NotifyFunc(Add, item) }
	- (notifier) myhub.Add(item)
	- (registrar) hub.RegisterFunc(myhub.Add, func(item){})
	 
	 We keep the de-coupling. The registrar doesn't knows the notifier and notifier doesn't knows about registrar at all.
	 The registrar can import the notifier with empty statement (_ importPath), inside that init the notifier will use
	 the myhub's notifier. Remember: The notifier executes first, which is anorthodox BUT at the previous commit
	 I made it to be able to 'wait' for the correct .RegisterHook in order to notify the hooks. So we don't have any issue
	 with that. The lib should be able to work both ways, as explained before (as an event listener and as a down-up notifier for func execution).

- [x] Prioritize per hook map's entry. Hooks are grouped to a Name (hook.Name and HooksMap's key, hook.Name is there to provide debug messages when needed).

	- The type should be an integer, but with some default Priority "levels".
	- Highest executes first.
	- No limits to the number that developer can use, by-default we will have 5-6 levels with iota * 100 starting from Idle (the lowest Priority).
	- Should be able to a hook to be prioritized from another hook, at runtime on build time, at any time it wants.
	- Should be able to accept negative values in order to reduce the priority when needed.
	- The hub will sort the routes per hook's changed priority .Name, so only hooks that are "linked" will be sorted.
	- The hook should call the hub's sort, so we should add an 'Owner *Hub' field.

## License

Unless otherwise noted, the source files are distributed
under the BSD-3 Clause License found in the [LICENSE file](LICENSE).