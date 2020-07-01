package util

import "reflect"

// BasicToInterface : Converts any basic input type to an interface.
// To be used with the Synchro library to help writing args more easily
func BasicToInterface(input ...interface{}) interface{} {
	return reflect.ValueOf(input).Interface().(interface{})
}
