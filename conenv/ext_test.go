package conenv

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestExtExecuted(t *testing.T) {
	keyExecuted := false
	valueExecuted := false
	validateExecuted := false
	fillExecuted := false
	ext := Extension{
		Key: func(b4 string, o Options, v reflect.Value) (af string) {
			keyExecuted = true
			return strings.ToUpper(b4)
		},
		Value: func(b4 string, o Options, v reflect.Value) (af string) {
			valueExecuted = true
			i, _ := strconv.ParseInt(b4, 10, 64)
			return strconv.FormatInt(i, 16)
		},
		Validate: func(o Options, f reflect.Value, n, v string) (err error) {
			validateExecuted = true
			if v == "" {
				err = errors.New("empty")
			}
			return
		},
		Fill: func(s string, o Options, v reflect.Value) (err error) {
			fillExecuted = true
			i, err := strconv.ParseInt(s, 16, 64)
			if err != nil {
				return
			}
			v.SetInt(i)
			return
		},
	}

	type test struct {
		Field int `env:"field,ext"`
	}
	var o test

	p := &Parser{}
	p.Register("ext", ext)

	os.Setenv("FIELD", "10")

	if err := p.Parse(&o); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !keyExecuted {
		t.Error("Key() is not executed")
	}
	if !valueExecuted {
		t.Error("Value() is not executed")
	}
	if !validateExecuted {
		t.Error("Validate() is not executed")
	}
	if !fillExecuted {
		t.Error("Fill() is not executed")
	}

	if o.Field != 10 {
		t.Fatalf("unexpected result: %d", o.Field)
	}
}
