package hooks

import (
	"reflect"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()
var errCacherType = reflect.TypeOf(func(error) {})

func Chain(a ...interface{}) interface{} {
	if len(a) == 0 {
		return nil
	}
	if len(a) < 2 {
		noChain := reflect.ValueOf(a[0])
		return func(payloads ...interface{}) {
			execFunc(noChain, payloads...)
		}
	}

	first := reflect.ValueOf(a[0])
	actionsLen := len(a)

	var errorCacher reflect.Value
	if last := reflect.ValueOf(a[len(a)-1]); last.Type() == errCacherType {
		errorCacher = last
		actionsLen-- // remove the error cacher from the valid actions
	}
	// to valid einai to prwto...
	return func(payloads ...interface{}) {
		// call the first, first of all.
		args, err := execFunc(first, payloads...)
		if err != nil {
			errorCacher.Call([]reflect.Value{reflect.ValueOf(err)})
			return
		}

		for i := 1; i < actionsLen; i++ {

			fn := reflect.ValueOf(a[i])

			if len(args) > 0 {
				fnType := fn.Type()
				numIn := fnType.NumIn()
				if len(args) != numIn {
					continue
				}

				for i := 0; i < numIn; i++ {
					iType := fnType.In(i)
					if !args[0].Type().AssignableTo(iType) {
						return
					}
				}
			}
			oldArgs := args

			args = fn.Call(args)

			if errorCacher.IsValid() {
				// first or second argument should be an error in order to send to the error cacher, if any.
				if len(args) == 1 {
					if err := args[0]; err.Type().Implements(errType) {
						if !err.IsNil() {
							errorCacher.Call([]reflect.Value{err})
							return // call the cacher and break here
						}
						args = oldArgs
					}
				} else if len(args) == 2 {
					if err := args[1]; err.Type().Implements(errType) {
						if !err.IsNil() {
							errorCacher.Call([]reflect.Value{err})
							return // call the cacher and break here
						}
						args = oldArgs // else, re-store the last arguments from the previous call in order to not .Call a nil error
					}
				}
			}

		}
	}
}
