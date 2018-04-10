package conenv

import (
	"errors"
	"reflect"

	"github.com/Ronmi/rtoolkit/reflkit"
)

// Extension handles custom options
//
// There are 4 stages when Parser executes custom handlers. It keeps same order you
// set in struct tag in each stage.
type Extension struct {
	// Key is first stage. It changes the name of envvar before actually load
	// its value.
	//
	// v would be the result of reflect.Value.Field(x)
	Key func(b4 string, opt Options, v reflect.Value) (after string)
	// Value is second stage. It changes the value loaded from envvar.
	//
	// v would be the result of reflect.Value.Field(x)
	//
	// If environment variable is not set, b4 will be empty string.
	Value func(b4 string, opt Options, v reflect.Value) (after string)
	// Validate is the third stage. It validates if anything goes wrong.
	//
	// v would be the result of reflect.Value.Field(x)
	//
	// If environment variable is not set, val will be empty string.
	Validate func(opt Options, v reflect.Value, name, val string) (err error)
	// Fill is last stage. It sets the struct field according to the value
	// string.
	//
	// v would be the result of reflect.Value.Field(x)
	//
	// If environment variable is not set, val will be empty string.
	//
	// Fill() is skipped if v is struct || pointer to struct || pointer with
	// empty val.
	//
	// As setting value multiple times is meaningless, only first matched setter
	// will be executed.
	Fill func(val string, opt Options, v reflect.Value) (err error)
}

// Parser manages extensions and parse struct tag to load value from envvar
//
// Zero value of parser is usable, just without any extension.
type Parser struct {
	handlers map[string]Extension

	// default value setter, use DefaultSetter if nil
	Setter *reflkit.StrConv
}

// DefaultParser returns a parser with default extensions registered.
//
// Supported extensions are:
//
//   - bin: BinIntExtension
//   - oct: OctIntExtension
//   - hex: HexIntExtension
//   - required: RequiredExtension
func DefaultParser() (p *Parser) {
	p = &Parser{}
	p.Register("bin", BinIntExtension)
	p.Register("oct", OctIntExtension)
	p.Register("hex", HexIntExtension)
	p.Register("required", RequiredExtension)

	return
}

// Register registers custom handler, key is case-sensitive
//
// Later with the same key overwrites older one.
//
// This method is not thread-safe, DO NOT call it among multiple goroutines.
func (p *Parser) Register(key string, h Extension) {
	if p.handlers == nil {
		p.handlers = make(map[string]Extension)
	}
	p.handlers[key] = h
}

// Parse parses struct tag and load values from environment variables
//
// Only pointer to struct is supported. Panics with any other type.
func (p *Parser) Parse(data interface{}) (err error) {
	return p.ParseWithPrefix(data, "")
}

// ParseWithPrefix is identical with Parse(), but prefixes some text before envvar
// name.
func (p *Parser) ParseWithPrefix(data interface{}, prefix string) (err error) {
	v := reflect.ValueOf(data).Elem()
	err, _ = p.doParse(v, prefix)
	return
}

// v MUST be settable struct, NOT POINTER!!!!!
func (p *Parser) doParse(v reflect.Value, prefix string) (err error, has bool) {
	for x, y := 0, v.NumField(); x < y; x++ {
		f := v.Type().Field(x)

		err, ok := p.parseField(v.Field(x), f, prefix)
		if err != nil {
			return err, has
		}
		has = has || ok
	}

	return
}

func (p *Parser) parseField(
	v reflect.Value, f reflect.StructField, prefix string,
) (err error, has bool) {
	o := parseOptions(f)

	name := o.envKey(prefix)

	// load handlers specified in custom key
	handlers := make([]Extension, 0, len(o.Custom))
	if p.handlers != nil {
		for _, k := range o.Custom {
			if h, ok := p.handlers[k]; ok {
				handlers = append(handlers, h)
			}
		}
	}

	// name processing
	for _, h := range handlers {
		if h.Key != nil {
			name = h.Key(name, o, v)
		}
	}

	str := o.envValue(name)

	// value processing
	for _, h := range handlers {
		if h.Value != nil {
			str = h.Value(str, o, v)
		}
	}

	// validation
	for _, h := range handlers {
		if h.Validate != nil {
			if err = h.Validate(o, v, name, str); err != nil {
				return
			}
		}
	}

	// struct fields need special treatment
	if v.Kind() == reflect.Struct {
		// inner struct, chaining name
		return p.doParse(v, name)
	}

	// pointer type needs special treatment, too
	if v.Kind() == reflect.Ptr {
		// pointer to struct need special treatment
		if z := v.Type().Elem(); z.Kind() == reflect.Struct {
			x := reflect.New(z)
			err, has = p.doParse(x.Elem(), name)
			if err != nil {
				return
			}
			if has {
				v.Set(x)
			}
			return
		}

		if str == "" {
			// no data, leave pointer type to be nil
			return
		}
		// pointer type, allocate space before set value
		reflkit.InitPtr(v)
		v = reflect.Indirect(v)
	}

	if str != "" {
		has = true
	}

	// fill value
	for _, h := range handlers {
		if h.Fill != nil {
			return h.Fill(str, o, v), has
		}
	}

	return p.defaultSetter(v, o, name, str), has
}

func (p *Parser) defaultSetter(
	v reflect.Value, o Options, name string, str string,
) (err error) {
	s := p.Setter
	if s == nil {
		s = DefaultSetter
	}

	if !s.SetValue(v, str) {
		err = errors.New("unsupported type")
	}

	return
}
