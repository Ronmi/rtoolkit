package conenv

import (
	"errors"
	"reflect"

	"github.com/Ronmi/rtoolkit/reflkit"
)

// HexIntExtension parses hexdecimal string into integer field
//
// Use it with any type other than integer might cause panic.
var HexIntExtension = Extension{
	Fill: func(val string, opt Options, v reflect.Value) (err error) {
		c := reflkit.IntHexConverter
		v = reflect.Indirect(v)
		if k := v.Kind(); k >= reflect.Uint && k <= reflect.Uint64 {
			c = reflkit.UintHexConverter
		}

		if !c(val, v) {
			err = errors.New("hexdecimal format error")
		}
		return
	},
}

// OctIntExtension parses octal string into integer field
//
// Use it with any type other than integer might cause panic.
var OctIntExtension = Extension{
	Fill: func(val string, opt Options, v reflect.Value) (err error) {
		c := reflkit.IntOctConverter
		v = reflect.Indirect(v)
		if k := v.Kind(); k >= reflect.Uint && k <= reflect.Uint64 {
			c = reflkit.UintOctConverter
		}

		if !c(val, v) {
			err = errors.New("octal format error")
		}
		return
	},
}

// BinIntExtension parses binary string into integer field
//
// Use it with any type other than integer might cause panic.
var BinIntExtension = Extension{
	Fill: func(val string, opt Options, v reflect.Value) (err error) {
		c := reflkit.IntBinConverter
		v = reflect.Indirect(v)
		if k := v.Kind(); k >= reflect.Uint && k <= reflect.Uint64 {
			c = reflkit.UintBinConverter
		}

		if !c(val, v) {
			err = errors.New("binary format error")
		}
		return
	},
}
