package conenv

import (
	"fmt"
	"os"
	"testing"
)

func TestAES(t *testing.T) {
	cases := []struct {
		plain string
		key   string
	}{
		// AES-128
		{"sdifhwe0f8y2jfnowe8fh9", "@#(RUJf0f9h13f(x"},
		// AES-192
		{"some random text", "AsuPeR5eUrEkeY1234567890"},
		// AES-256
		{"osihfbohf-hhhup9@PJf-293ufjkj", "@#(RUJf0fl9h13f(UF0f923RFRHC($GU"},
	}

	for _, c := range cases {
		code, err := AESEncrypt(c.key, c.plain)
		if err != nil {
			t.Fatalf("unexpected error when encrypting: %s", err)
		}

		if code == c.plain {
			t.Fatalf("not encrypted!?")
		}

		val, err := AESDecrypt(c.key, code)
		if err != nil {
			t.Fatalf("unexpected error when decrypting: %s", err)
		}

		if val != c.plain {
			t.Fatalf(`expected tet to be "%s", got "%s"`, c.plain, val)
		}
	}
}

func ExampleAESExtension() {
	type dbconf struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
		User string `env:"USER"`
		Pass string `env:"PASS,aes"`
	}
	key := "my secret key45678901234" // 24 bytes, AES-192 used

	// prepare envvar, simulating what you do in docker.
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("PORT", "3306")
	os.Setenv("USER", "user")
	// password is encrypted elesewhere using fowllowing code
	// AESEncrypt(key, "mysql password")
	os.Setenv("PASS", "g3Dv1jqTu5hdzNEEDQXdZk1cCUTngaLS")

	p := &Parser{}
	p.Register("aes", DESExtension(key)) // important!

	var conf dbconf
	p.Parse(&conf)
	fmt.Printf("%+v", conf)

	// output: {Host:127.0.0.1 Port:3306 User:user Pass:mysql password}
}
