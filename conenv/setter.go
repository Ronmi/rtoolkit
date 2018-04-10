package conenv

import (
	"reflect"
	"time"

	"github.com/Ronmi/rtoolkit/reflkit"
)

func durSetter(str string, target reflect.Value) (ok bool) {
	dur, err := time.ParseDuration(str)
	if err != nil {
		return
	}

	target.Set(reflect.ValueOf(dur))

	return true
}

// DefaultSetter is default value setter if you leave Parser.Setter nil.
//
// It supports all types/kinds that reflkit.DefaultStrConv support, with extra
// time.Duration support (via time.ParseDuration).
var DefaultSetter *reflkit.StrConv

func init() {
	DefaultSetter = reflkit.DefaultStrConv()
	var dur time.Duration
	DefaultSetter.ByType[reflect.TypeOf(dur)] = durSetter
}
