package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Ronmi/rtoolkit/conenv"
)

const rsacode = `
func loadPriv() (key *rsa.PrivateKey, err error) {
	data, err := ioutil.ReadFile("key.priv")
	if err != nil {
		return
	}

	return x509.ParsePKCS1PrivateKey(data)
}
func loadPub() (key *rsa.PublicKey, err error) {
	data, err := ioutil.ReadFile("key.pub")
	if err != nil {
		return
	}

	return x509.ParsePKCS1PublicKey(data)
}
`

func loadPriv(fn string) (key *rsa.PrivateKey, err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}

	return x509.ParsePKCS1PrivateKey(data)
}

func loadPub(fn string) (key *rsa.PublicKey, err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}

	return x509.ParsePKCS1PublicKey(data)
}

func genrsa() {
	var bitSize int
	flag.IntVar(&bitSize, "size", 2048, "key size in bit")
	flag.Parse()

	log.Print("Generating key pair...")
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	priv := x509.MarshalPKCS1PrivateKey(key)
	pub := x509.MarshalPKCS1PublicKey(&key.PublicKey)

	log.Print("Writing public key to pubkey.pub...")
	if err := ioutil.WriteFile("key.pub", pub, 0600); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Writing private key to key.priv...")
	if err := ioutil.WriteFile("key.priv", priv, 0600); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Here's example code:")
	fmt.Println(rsacode)
}

func rsaenc() {
	var (
		keyfile string
		data    string
	)
	flag.StringVar(&keyfile, "k", "key.pub", "public key file.")
	flag.StringVar(&data, "data", "", "data to be encrypted (required).")
	flag.Parse()

	if data == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("Loading public key from %s...", keyfile)
	key, err := loadPub(keyfile)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("encrypting data...")
	code, err := conenv.RSAEncrypt(key, data)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	fmt.Println()
	fmt.Println(code)
}

func rsadec() {
	var (
		keyfile string
		data    string
	)
	flag.StringVar(&keyfile, "k", "key.priv", "private key file.")
	flag.StringVar(&data, "data", "", "cipher text to be decrypted (required).")
	flag.Parse()

	if data == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("Loading private key from %s...", keyfile)
	key, err := loadPriv(keyfile)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("decrypting data...")
	code, err := conenv.RSADecrypt(key, data)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	fmt.Println()
	fmt.Println(code)
}
