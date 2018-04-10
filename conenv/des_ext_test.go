package conenv

import (
	"fmt"
	"os"
	"testing"
)

func TestDES(t *testing.T) {
	cases := []struct {
		plain string
		key   string
	}{
		{"some random text", "AsuPeR5eUrEkeY1234567890"},
		{"osihfbohf-hhhup9@PJf-293ufjkj", "@#(RUJf0f9h13f(UF0f923RF"},
	}

	for _, c := range cases {
		code, err := DESEncrypt(c.key, c.plain)
		if err != nil {
			t.Fatalf("unexpected error when encrypting: %s", err)
		}

		if code == c.plain {
			t.Fatalf("not encrypted!?")
		}

		val, err := DESDecrypt(c.key, code)
		if err != nil {
			t.Fatalf("unexpected error when decrypting: %s", err)
		}

		if val != c.plain {
			t.Fatalf(`expected tet to be "%s", got "%s"`, c.plain, val)
		}
	}
}

func ExampleDESExtension() {
	type dbconf struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
		User string `env:"USER"`
		Pass string `env:"PASS,des"`
	}
	key := "my secret key45678901234"

	// prepare envvar, simulating what you do in docker.
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("PORT", "3306")
	os.Setenv("USER", "user")
	// password is encrypted elsewhere using following code
	// DESEncrypt(key, "mysql password")
	os.Setenv("PASS", "g3Dv1jqTu5hdzNEEDQXdZk1cCUTngaLS")

	p := &Parser{}
	p.Register("des", DESExtension(key)) // important!

	var conf dbconf
	p.Parse(&conf)
	fmt.Printf("%+v", conf)

	// output: {Host:127.0.0.1 Port:3306 User:user Pass:mysql password}
}
