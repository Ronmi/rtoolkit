package conenv

import (
	"os"
	"reflect"
	"strings"
)

// Options is the data parsed from struct tag, which is needed by extension.
type Options struct {
	Name    string
	Default string
	Custom  []string
}

func (o *Options) envKey(prefix string) (ret string) {
	return prefix + o.Name
}

func (o *Options) envValue(key string) (ret string) {
	ret, ok := os.LookupEnv(key)
	if !ok {
		ret = o.Default
	}

	return
}

func parseOptions(f reflect.StructField) (ret Options) {
	e := f.Tag.Get("env")
	if e == "" {
		ret.Name = f.Name
		return
	}

	arr := strings.Split(e, ",")
	ret.Name = arr[0]
	ret.Custom = arr[1:]
	ret.Default = f.Tag.Get("envDefault")

	return
}
