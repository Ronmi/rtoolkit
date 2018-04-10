package conenv

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"testing"
)

func TestRSA(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("cannot generate key pair: %s", err)
	}
	plain := "my secret text"

	code, err := RSAEncrypt(&key.PublicKey, plain)
	if err != nil {
		t.Fatalf("unexpected error wher encrypting: %s", err)
	}

	if code == plain {
		t.Fatalf("not encrypted!?")
	}

	val, err := RSADecrypt(key, code)
	if err != nil {
		t.Fatalf("unexpected error when decrypting: %s", err)
	}

	if val != plain {
		t.Fatalf(`expected text to be "%s", got "%s"`, plain, val)
	}
}

func ExampleRSAExtension() {
	type dbconf struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
		User string `env:"USER"`
		Pass string `env:"PASS,rsa"`
	}

	// generate key and code
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	code, _ := RSAEncrypt(&key.PublicKey, "mysql password")

	// prepare envvar, simulating what you do in docker.
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("PORT", "3306")
	os.Setenv("USER", "user")
	os.Setenv("PASS", code)

	p := &Parser{}
	p.Register("rsa", RSAExtension(key)) // important!

	var conf dbconf
	p.Parse(&conf)
	fmt.Printf("%+v", conf)

	// output: {Host:127.0.0.1 Port:3306 User:user Pass:mysql password}
}
