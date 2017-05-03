package hooks

import (
	"fmt"
	"reflect"
	"runtime"
)

func isFunc(typ reflect.Type) bool {
	return typ.Kind() == reflect.Func
}

func execFunc(val reflect.Value, arguments ...interface{}) ([]reflect.Value, error) {
	// build the arguments for the reflect.Call
	in := make([]reflect.Value, len(arguments))

	for i := range arguments {
		in[i] = reflect.ValueOf(arguments[i])
	}

	// validate the arguments
	typ := val.Type()
	isVariadic := typ.IsVariadic()

	if expected, got := typ.NumIn(), len(in); expected != got {
		if expected > 0 {
			if !isVariadic {
				return nil, fmt.Errorf("expected %d arguments but got %d", got, expected)
			}
		}
	}
	if typ.NumIn() > 0 {
		for i := range in {

			if i == typ.NumIn()-1 {
				if isVariadic && typ.In(typ.NumIn()-1).Kind() == reflect.Slice {
					break
				}
			}

			if expected, got := typ.In(i), in[i].Type(); expected != got {
				if !got.AssignableTo(expected) {
					return nil, fmt.Errorf("arguments are not the same type, expected %s but got %s", expected, got)
				}
			}
		}
	}
	/*
		if isVariadic && typ.NumIn() > 0 {

			// 	slice := reflect.MakeSlice(typ.In(0), typ.NumIn(), typ.NumIn())
			// 	reflect.Append(slice, reflect.ValueOf(arguments))
			// 	return val.CallSlice([]reflect.Value{slice}), nil
			// }
			// ---
			// if variadicType.Kind() == reflect.String {
			// 	var argsStr []string
			// 	for i := range arguments {
			// 		argsStr = append(argsStr, arguments[i].(string))
			// 	}
			// }

			// for i := range arguments {
			// 	in = append(in, reflect.ValueOf(arguments[i]))
			// }
			// ---
			// var convertedIn []interface{}
			// argsVal := reflect.ValueOf(in)

			// println("args val type of: " + argsVal.Index(0).Type().Name())

			// argumentsAsVariadic := argsVal.Convert(typ.In(0)) // reflect.SliceOf(
			// inVariadic := make([]reflect.Value, 0, 1)
			// inVariadic = append(inVariadic, argumentsAsVariadic)
			// ---
			argumentsAsVariadic := []interface{}{
				arguments,
			}
			inPreVariadic := reflect.ValueOf(argumentsAsVariadic)
			inVariadic := inPreVariadic.Convert(typ.In(0))
			// ---
			return val.CallSlice([]reflect.Value{inVariadic}), nil
		}*/

	// execute the function, no need for .CallSlice, it's working  well with .Call on variadics too.
	return val.Call(in), nil
}

func nameOfFunc(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}
