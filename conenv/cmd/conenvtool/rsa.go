package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ronmi/rtoolkit/conenv"
)

const rsacode = `
func loadRSAPriv() (key *rsa.PrivateKey, err error) {
	f, err := os.Open("rsa.priv")
	if err != nil {
		return
	}
	defer f.Close()

	key = &rsa.PrivateKey{}
	dec := gob.NewDecoder(f)
	return key, dec.Decode(key)
}
func loadRSAPub() (key *rsa.PublicKey, err error) {
	f, err := os.Open("rsa.pub")
	if err != nil {
		return
	}
	defer f.Close()

	key = &rsa.PublicKey{}
	dec := gob.NewDecoder(f)
	return key, dec.Decode(key)
}
`

const rsaembed = `
var rsaPrivateKey *rsa.PrivateKey

func init() {
	rsPrivateKey = &rsa.PrivateKey{}
	data := %#v
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	dec.Decode(rsaPrivateKey)
}
`

func loadRSAPriv(fn string) (key *rsa.PrivateKey, err error) {
	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()

	key = &rsa.PrivateKey{}
	dec := gob.NewDecoder(f)
	err = dec.Decode(key)
	return
}

func loadRSAPub(fn string) (key *rsa.PublicKey, err error) {
	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()

	key = &rsa.PublicKey{}
	dec := gob.NewDecoder(f)
	err = dec.Decode(key)
	return
}

func genrsa() {
	var (
		bitSize int
		embed   bool
	)
	flag.IntVar(&bitSize, "size", 2048, "key size in bit")
	flag.BoolVar(&embed, "embed", false, "generate codes to embed private key in your program instead")
	flag.Parse()

	log.Print("Generating key pair...")
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Writing public key to rsa.pub...")
	if err := saveGobData("rsa.pub", key.PublicKey); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Writing private key to rsa.priv...")
	if err := saveGobData("rsa.priv", key); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Here's example code:")
	code := rsacode
	if embed {
		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		enc.Encode(key)
		code = fmt.Sprintf(rsaembed, buf.Bytes())
	}
	fmt.Println(code)
}

func rsaenc() {
	var (
		keyfile string
		data    string
	)
	flag.StringVar(&keyfile, "k", "rsa.pub", "public key file.")
	flag.StringVar(&data, "data", "", "data to be encrypted (required).")
	flag.Parse()

	if data == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("Loading public key from %s...", keyfile)
	key, err := loadRSAPub(keyfile)
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
	flag.StringVar(&keyfile, "k", "rsa.priv", "private key file.")
	flag.StringVar(&data, "data", "", "cipher text to be decrypted (required).")
	flag.Parse()

	if data == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("Loading private key from %s...", keyfile)
	key, err := loadRSAPriv(keyfile)
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
