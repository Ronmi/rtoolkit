package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/Ronmi/rtoolkit/conenv"
)

func desenc() {
	var (
		key  string
		data string
		h    bool
	)
	flag.StringVar(&key, "k", "", "key string (required).")
	flag.StringVar(&data, "data", "", "data to be encrypted (required).")
	flag.BoolVar(&h, "hex", false, "decode key string with hex.Decode().")
	flag.Parse()

	if key == "" || data == "" {
		flag.PrintDefaults()
		return
	}

	k := []byte(key)
	if h {
		log.Print("Decoding key with hex.Decode()...")
		var err error
		k, err = hex.DecodeString(key)
		if err != nil {
			log.Fatalf("unexpected error: %s", err)
		}
	}

	log.Print("Encrypting data...")
	code, err := conenv.DESEncrypt(string(k), data)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	fmt.Println()
	fmt.Println(code)
}

func desdec() {
	var (
		key  string
		data string
		h    bool
	)
	flag.StringVar(&key, "k", "", "key string (required).")
	flag.StringVar(&data, "data", "", "data to be decrypted (required).")
	flag.BoolVar(&h, "hex", false, "decode key string with hex.Decode().")
	flag.Parse()

	if key == "" || data == "" {
		flag.PrintDefaults()
		return
	}

	k := []byte(key)
	if h {
		log.Print("Decoding key with hex.Decode()...")
		var err error
		k, err = hex.DecodeString(key)
		if err != nil {
			log.Fatalf("unexpected error: %s", err)
		}
	}

	log.Print("Decrypting data...")
	code, err := conenv.DESDecrypt(string(k), data)
	if err != nil {
		log.Fatalf("unexpected error: %s", err)
	}

	fmt.Println()
	fmt.Println(code)
}
