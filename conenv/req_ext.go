package conenv

import (
	"errors"
	"reflect"
)

// RequiredExtension ensures value of envvar is not empty string
var RequiredExtension = Extension{
	Validate: func(opt Options, v reflect.Value, name, str string) (err error) {
		if str == "" {
			err = errors.New("required field is not set")
		}
		return
	},
}
