package conenv

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestAllSupportedTypes(t *testing.T) {
	type inner struct {
		A string
		B string
	}
	type test struct {
		Int       int
		Int8      int8
		Int16     int16
		Int32     int32
		Int64     int64
		Uint      uint
		Uint8     uint8
		Uint16    uint16
		Uint32    uint32
		Uint64    uint64
		Float32   float32
		Float64   float64
		Bool      bool
		Duration  time.Duration
		String    string
		ByteSlice []byte
		RuneSlice []rune
		Inner     inner
		InnerPtr  *inner

		// decorated
		Abc       string `env:"ABC"`
		DecoInner inner  `env:"INNER_"`

		// nil pointer
		Unset      *int
		UnsetInner *inner
	}

	expect := test{
		Int:        1234567890,
		Int8:       123,
		Int16:      12345,
		Int32:      1234567890,
		Int64:      12345678901234,
		Uint:       1234567890,
		Uint8:      123,
		Uint16:     12345,
		Uint32:     1234567890,
		Uint64:     12345678901234,
		Float32:    1.5,
		Float64:    3.1415926,
		Bool:       true,
		Duration:   2 * time.Second,
		String:     "string",
		ByteSlice:  []byte("[]byte"),
		RuneSlice:  []rune("[]rune"),
		Inner:      inner{A: "innera", B: "innerb"},
		InnerPtr:   &inner{A: "ptra", B: "ptrb"},
		Abc:        "abc",
		DecoInner:  inner{A: "inner_a", B: "inner_b"},
		Unset:      nil,
		UnsetInner: nil,
	}

	os.Clearenv()
	os.Setenv("Int", "1234567890")
	os.Setenv("Int8", "123")
	os.Setenv("Int16", "12345")
	os.Setenv("Int32", "1234567890")
	os.Setenv("Int64", "12345678901234")
	os.Setenv("Uint", "1234567890")
	os.Setenv("Uint8", "123")
	os.Setenv("Uint16", "12345")
	os.Setenv("Uint32", "1234567890")
	os.Setenv("Uint64", "12345678901234")
	os.Setenv("Float32", "1.5")
	os.Setenv("Float64", "3.1415926")
	os.Setenv("Bool", "true")
	os.Setenv("Duration", "2s")
	os.Setenv("String", "string")
	os.Setenv("ByteSlice", "[]byte")
	os.Setenv("RuneSlice", "[]rune")
	os.Setenv("InnerA", "innera")
	os.Setenv("InnerB", "innerb")
	os.Setenv("InnerPtrA", "ptra")
	os.Setenv("InnerPtrB", "ptrb")
	os.Setenv("ABC", "abc")
	os.Setenv("INNER_A", "inner_a")
	os.Setenv("INNER_B", "inner_b")

	p := DefaultParser()
	var x test
	if err := p.Parse(&x); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(x.InnerPtr, expect.InnerPtr) {
		t.Fatalf("Unexpected inner pointer result: %+v", x.InnerPtr)
	}
	x.InnerPtr = expect.InnerPtr
	if !reflect.DeepEqual(x, expect) {
		t.Fatalf("Unexpected result: %+v", x)
	}
}

func ExampleParser_ParseWithPrefix() {
	type dbconf struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT,hex" envDefault:"0cea"`
		User string `env:"USER"`
		Pass string `env:"PASS"`
	}

	type conf struct {
		Master dbconf `env:"MASTER_"`
		Slave  dbconf `env:"SLAVE_"`
	}

	os.Setenv("DB_MASTER_HOST", "192.168.0.1")
	os.Setenv("DB_MASTER_PORT", "33fa")
	os.Setenv("DB_MASTER_USER", "user")
	os.Setenv("DB_MASTER_PASS", "pass")
	os.Setenv("DB_SLAVE_HOST", "192.168.0.2")
	os.Setenv("DB_SLAVE_USER", "user")
	os.Setenv("DB_SLAVE_PASS", "pass")

	var c conf
	p := DefaultParser()
	if err := p.ParseWithPrefix(&c, "DB_"); err != nil {
		// error handling
	}

	fmt.Printf("%+v", c)

	// output: {Master:{Host:192.168.0.1 Port:13306 User:user Pass:pass} Slave:{Host:192.168.0.2 Port:3306 User:user Pass:pass}}
}
