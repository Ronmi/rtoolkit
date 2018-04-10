package conenv

import (
	"crypto/dsa"
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

var testDSAPubKey *dsa.PublicKey

func init() {
	testDSAPubKey = &testDSAKey.PublicKey
}

func Example() {
	// simulating your setup in container
	os.Setenv("MYSQL_MASTER_HOST", "192.168.0.1")
	os.Setenv("MYSQL_MASTER_PORT", "33fa")
	os.Setenv("MYSQL_MASTER_USER", "user")
	os.Setenv("MYSQL_MASTER_PASS", "ML7b/Db8fo0tuzmsgONZPjc7chzFa3wXTSViPyDa12X8inaHqvzgEigyqBAyR/rmL/69c+eCOJqs5eiasp4TF1pLYBYyp+4KWcizwqsioWgiSU5cxd4I7LxYg1UOItmsPCHshzP0lOZo4hDHn+GoiVDuRFkuPIhYhO9Co6g54z7aolHAcUT5JgVkNpfENFxiSTBBLxQVCQEoJcT29RXC7u1tmQ7drQTH9Dz0lDDMCGqugFz1k6FzBOOuzqT5mnYEniqJFTagdOHs2DU/Rv+35ZeU1vIqXvYor4U8yMm1S2mBvMRgXBE/AVeip1EY7Y8eewra7UFt5UXiW0cMvdHsfA==")
	os.Setenv("MYSQL_MASTER_PASS_SIGN", "Cv+BBQEC/4QAAAAl/4IAIQJ1gEUv23YL5j5YkuS1Tzf5fcTcR/C8Q/QFWaUjIb1A5CX/ggAhArfRGm5HDHWPGnsnSareByKudJ6MWp+mhrzlOXJjVhXa")
	os.Setenv("MYSQL_SLAVE_HOST", "192.168.0.2")
	os.Setenv("MYSQL_SLAVE_USER", "user")
	os.Setenv("MYSQL_SLAVE_PASS", "SKSG78sOMHr+7k6vIvM+/BEb41v0QYWtVnvoCh/iikEs9bKU6859l0m2XjSkKVxSlGy3Wo8sx/rXNl8FwWl/O29ori4JAep1bM/kPCXknOFe8HmPeqWNBk5TWc5zxnhKLXigds3Sqf06jkQ6/LqJgaUrdMccHWFTImHiLa5YX0svV86TRxpKEcZ/H/J+ALRaZFtDSs9JvPn8DpB3T96SWKp3bDX63hlf+F2ijJnk89kL6xB8uKBi9huNSuS9Vy2kdtIVdpLw6I7Wq3tpeAoGbSIay5h/i6kDtVUR5jQxeNZgtIFaWaIChoARgeDoIf4FttUwwmPNG/lzXzcWPfO1qA==")
	os.Setenv("MYSQL_SLAVE_PASS_SIGN", "Cv+BBQEC/4QAAAAl/4IAIQKITi303J64pPQ99mDDkcaXhGx+L5iwvofJUrNp2MqljyX/ggAhAioAW3nUpvIXjpBfoxLWNgEEJ1b/3Sjag2Ql9fGsqJwm")
	// Best practice is generating cipher tex and signatures for each variable
	// separately, even if they are same value in fact.

	// =========== real code begin here

	// define config types
	type dbconf struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT,hex" envDefault:"0cea"`
		User string `env:"USER"`
		Pass string `env:"PASS,enc,required,sign"`
	}
	type conf struct {
		Master dbconf `env:"MASTER_"`
		Slave  dbconf `env:"SLAVE_"`
	}

	// load config from envvar
	var c conf
	p := DefaultParser()
	p.Register("enc", RSAExtension(testRSAKey))
	p.Register("sign", DSAExtension(testDSAPubKey))
	if err := p.ParseWithPrefix(&c, "MYSQL_"); err != nil {
		// handle error here
	}
	// =========== end of real code

	// print values for package testing purpose
	fmt.Printf("%+v", c)
	// output: {Master:{Host:192.168.0.1 Port:13306 User:user Pass:pass} Slave:{Host:192.168.0.2 Port:3306 User:user Pass:pass}}
}
