package main

import (
	"fmt"
	"os"
)

var handlers = map[string]func(){
	"help":    help,
	"genrsa":  genrsa,
	"rsaenc":  rsaenc,
	"rsadec":  rsadec,
	"aesenc":  aesenc,
	"aesdec":  aesdec,
	"desenc":  desenc,
	"desdec":  desdec,
	"gendsa":  gendsa,
	"dsasign": dsasign,
	"dsavfy":  dsavfy,
}

func main() {
	var act string

	for x, o := range os.Args[1:] {
		if o[0] != '-' {
			act = o
			os.Args = append(os.Args[:x], os.Args[x+1:]...)
			break
		}
	}

	h, ok := handlers[act]
	if !ok {
		h = handlers["help"]
	}

	h()
}

func help() {
	fmt.Print(`Usage: conenvtool action [options...]

Supported actions:

  help:    display this text.
  genrsa:  generate rsa key pair, with short piece of example code.
  rsaenc:  encrypt string data with previously generated key.
  rsadec:  decrypt cipher text encrypted using this tool.
  aesenc:  encrypt string data with AESEncrypt.
  aesdec:  decrypt string data with AESDecrypt.
  desenc:  encrypt string data with DESEncrypt.
  desdec:  decrypt string data with DESDecrypt.
  gendsa:  generate dsa key pairs.
  dsasign: sign string data with DSASign.
  dsavfy:  verify signature with DSAVerify.

Example:
  
  conenvtool genrsa -h
     show supported options. plz ignore usage line it shows. at this time, I am
     too lazy to fix.

  conenvtool aesenc -k "my super secure secret key123456" -data "my secret data"
     encrypt data with 32 bytes key

  conenvtool aesdec -k "my super secure secret key123456" -data "xIF5hGhzymxaGExs7cO5Bp75nTV1f9vmKP7KRxaDMtk="
     decrypt data (debug purpose)
`)
}
