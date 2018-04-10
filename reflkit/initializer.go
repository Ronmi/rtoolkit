package reflkit

import "reflect"

// InitPtr allocates space and points target to it
//
// target MUST be settable pointer or it panics.
func InitPtr(target reflect.Value) {
	v := reflect.New(target.Type().Elem())
	target.Set(v)
}
