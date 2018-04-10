package main

import (
	"bytes"
	"crypto/dsa"
	"crypto/rand"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Ronmi/rtoolkit/conenv"
)

const dsacode = `
func loadDSAPriv() (key *dsa.PrivateKey, err error) {
	data, err := ioutil.ReadFile("dsa.priv")
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	key = &dsa.PrivateKey{}
	err = dec.Decode(key)
	return
}
func loadDSAPub() (key *dsa.PublicKey, err error) {
	data, err := ioutil.ReadFile("dsa.pub")
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	key = &dsa.PublicKey{}
	err = dec.Decode(key)
	return
}
`

const dsaembed = `
var dsaPublicKey *dsa.PublicKey

func init() {
	dsaPublicKey = &dsa.PublicKey{}
	data := %#v
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	dec.Decode(dsaPublicKey)
}
`

func loadDSAPriv(fn string) (key *dsa.PrivateKey, err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	key = &dsa.PrivateKey{}
	err = dec.Decode(key)
	return
}

func loadDSAPub(fn string) (key *dsa.PublicKey, err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	key = &dsa.PublicKey{}
	err = dec.Decode(key)
	return
}

func saveGobData(fn string, data interface{}) (err error) {
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	defer f.Sync()

	enc := gob.NewEncoder(f)
	return enc.Encode(data)
}

func gendsa() {
	var (
		embed bool
		size  int
	)
	flag.BoolVar(&embed, "embed", false, "generate codes to embed pubic key in your program instead")
	flag.IntVar(&size, "size", 2, "parameter size. valid values are 0: L1024N160, 1: L2048N224, 2: L2048N256, 3: L3072N256")
	flag.Parse()

	validSize := map[int]dsa.ParameterSizes{
		0: dsa.L1024N160,
		1: dsa.L2048N224,
		2: dsa.L2048N256,
		3: dsa.L3072N256,
	}
	var sz dsa.ParameterSizes
	sz, ok := validSize[size]
	if !ok {
		flag.PrintDefaults()
		return
	}

	key := &dsa.PrivateKey{}
	log.Print("Generating DSA parameters...")
	err := dsa.GenerateParameters(&key.PublicKey.Parameters, rand.Reader, sz)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Generating key pair...")
	if err = dsa.GenerateKey(key, rand.Reader); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Writing private key to dsa.priv...")
	if err = saveGobData("dsa.priv", key); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}
	log.Print("Writing public key to dsa.pub...")
	if err = saveGobData("dsa.pub", key.PublicKey); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Here's example code:")
	code := dsacode
	if embed {
		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		enc.Encode(key)
		code = fmt.Sprintf(dsaembed, buf.Bytes())
	}
	fmt.Println(code)
}

func dsasign() {
	var (
		keyfile string
		data    string
	)
	flag.StringVar(&keyfile, "k", "dsa.priv", "private key file.")
	flag.StringVar(&data, "data", "", "data to be signed (required).")
	flag.Parse()

	if data == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("Loading private key from %s...", keyfile)
	key, err := loadDSAPriv(keyfile)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Signing data...")
	code, err := conenv.DSASign(key, data)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	fmt.Println()
	fmt.Println(code)
}

func dsavfy() {
	var (
		keyfile string
		data    string
		value   string
	)
	flag.StringVar(&keyfile, "k", "dsa.pub", "public key file.")
	flag.StringVar(&data, "data", "", "data to be verified (required).")
	flag.StringVar(&value, "v", "", "original data (required).")
	flag.Parse()

	if data == "" || value == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("Loading public key from %s...", keyfile)
	key, err := loadDSAPub(keyfile)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	log.Print("Verifying data...")
	if err := conenv.DSAVerify(key, value, data); err != nil {
		log.Fatalf("unexpected error: %s", err)
	}
	log.Print("Data is matching with signature.")
}
