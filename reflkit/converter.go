package reflkit

import (
	"reflect"
	"strconv"
)

// ConvFunc fills target using the string.
//
// Implementations do not have to call target.CanSet(), it is validated in Set().
type ConvFunc func(str string, target reflect.Value) (ok bool)

// StrConv provides few methods to convert string to other types.
type StrConv struct {
	ByType map[reflect.Type]ConvFunc
	ByKind map[reflect.Kind]ConvFunc
}

// Set sets value of target using str.
//
// Only integers/floats/boolean/strings are supported.
//
// If target is not settable or supported, return false.
func (c StrConv) Set(target interface{}, str string) (ok bool) {
	v := reflect.Indirect(reflect.ValueOf(target))
	return c.SetValue(v, str)
}

// SetValue sets value of target using str.
//
// Only integers/floats/boolean/strings are supported.
//
// If v is not settable or supported, return false.
func (c StrConv) SetValue(v reflect.Value, str string) (ok bool) {
	if !v.CanSet() {
		return
	}

	if conv, ok := c.ByType[v.Type()]; ok {
		return conv(str, v)
	}

	conv, ok := c.ByKind[v.Kind()]
	if !ok {
		return
	}

	return conv(str, v)
}

// BoolConverter converts string to bool using strconv.ParseBool()
func BoolConverter(str string, target reflect.Value) (ok bool) {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return
	}

	target.SetBool(b)
	return true
}

// IntDecConverter converts string to 10-based signed intergers using
// strconv.ParseInt()
func IntDecConverter(str string, target reflect.Value) (ok bool) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}

	target.SetInt(i)
	return true
}

// IntHexConverter converts string to 16-based signed intergers using
// strconv.ParseInt()
func IntHexConverter(str string, target reflect.Value) (ok bool) {
	i, err := strconv.ParseInt(str, 16, 64)
	if err != nil {
		return
	}

	target.SetInt(i)
	return true
}

// IntOctConverter converts string to 8-based signed intergers using
// strconv.ParseInt()
func IntOctConverter(str string, target reflect.Value) (ok bool) {
	i, err := strconv.ParseInt(str, 8, 64)
	if err != nil {
		return
	}

	target.SetInt(i)
	return true
}

// IntBinConverter converts string to 2-based signed intergers using
// strconv.ParseInt()
func IntBinConverter(str string, target reflect.Value) (ok bool) {
	i, err := strconv.ParseInt(str, 2, 64)
	if err != nil {
		return
	}

	target.SetInt(i)
	return true
}

// UintDecConverter converts string to 10-based unsigned intergers using
// strconv.ParseUint()
func UintDecConverter(str string, target reflect.Value) (ok bool) {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return
	}

	target.SetUint(i)
	return true
}

// UintHexConverter converts string to 16-based signed uintergers using
// strconv.ParseUint()
func UintHexConverter(str string, target reflect.Value) (ok bool) {
	i, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return
	}

	target.SetUint(i)
	return true
}

// UintOctConverter converts string to 8-based signed uintergers using
// strconv.ParseUint()
func UintOctConverter(str string, target reflect.Value) (ok bool) {
	i, err := strconv.ParseUint(str, 8, 64)
	if err != nil {
		return
	}

	target.SetUint(i)
	return true
}

// UintBinConverter converts string to 2-based signed uintergers using
// strconv.ParseUint()
func UintBinConverter(str string, target reflect.Value) (ok bool) {
	i, err := strconv.ParseUint(str, 2, 64)
	if err != nil {
		return
	}

	target.SetUint(i)
	return true
}

// Floatconverter converts string to floating pointers using strconv.ParseFloat()
func FloatConverter(str string, target reflect.Value) (ok bool) {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return
	}

	target.SetFloat(f)
	return true
}

// StringConverter does nothing, just calling target.SetString()
func StringConverter(str string, target reflect.Value) (ok bool) {
	v := reflect.ValueOf(str).Convert(target.Type())
	target.Set(v)
	return true
}

// DefaultStrConv retrieves default converter implementations.
//
// It does not support pointer types, you have to handle it on your own.
func DefaultStrConv() (ret *StrConv) {
	var (
		b []byte
		r []rune
	)

	return &StrConv{
		ByType: map[reflect.Type]ConvFunc{
			reflect.TypeOf(b): StringConverter,
			reflect.TypeOf(r): StringConverter,
		},
		ByKind: map[reflect.Kind]ConvFunc{
			reflect.Bool:    BoolConverter,
			reflect.Int:     IntDecConverter,
			reflect.Int8:    IntDecConverter,
			reflect.Int16:   IntDecConverter,
			reflect.Int32:   IntDecConverter,
			reflect.Int64:   IntDecConverter,
			reflect.Uint:    UintDecConverter,
			reflect.Uint8:   UintDecConverter,
			reflect.Uint16:  UintDecConverter,
			reflect.Uint32:  UintDecConverter,
			reflect.Uint64:  UintDecConverter,
			reflect.Float32: FloatConverter,
			reflect.Float64: FloatConverter,
			reflect.String:  StringConverter,
		},
	}
}
