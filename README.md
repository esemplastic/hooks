# Hooks

## Idea 

- [x] The RegisterHook would be able to catch any parameters
	from the Notify. So inside RegisterHook will be the action.

- [x] The RegisterHook should be de-coupled from the action
	in order to be able to use that library when import cycle
	is not allowed and code refactor is not possible.

- [ ] The RegisterHook should be able to be work as an
	"event listener" too, so it should be able to not work
	as an action too. That makes the Notify to be able to execute
	the action function and notify the listeners.

- [ ] We can do both RegisterHook and
	Notify to accept an interface{} and transfer the
	function parameters or return values between RegisterHook
	and Notify with the reflect package help.
	But we will lose type safety (already losed, the library
	provides some helpers but it's not evaluate the
	action function and its parameters).
	So we need to:

- [ ] Convert any function to a hook, the function can be
	unique with the help of runtime package
	(by getting the full name of the func in the source code).
	But if we do that we lose the de-coupling described on second paragraph.



## License

Unless otherwise noted, the source files are distributed
under the BSD-3 Clause License found in the [LICENSE file](LICENSE).