package conenv

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"io"
	"reflect"
	"strings"
)

type sharedDecryptFunc func(key, code string) (val string, err error)

func genSharedKeyExtension(key string, f sharedDecryptFunc) (ext Extension) {
	return Extension{
		Value: func(b4 string, o Options, v reflect.Value) (af string) {
			af, err := f(key, b4)
			if err != nil {
				af = ""
			}

			return
		},
	}
}

func blockEncrypt(val string, b cipher.Block) (code string, err error) {
	// padding data
	if x := len(val) % b.BlockSize(); x != 0 {
		val += strings.Repeat(" ", b.BlockSize()-x)
	}

	ciphertext := make([]byte, b.BlockSize()+len(val))
	iv := ciphertext[:b.BlockSize()]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	enc := cipher.NewCTR(b, iv)
	enc.XORKeyStream(ciphertext[b.BlockSize():], []byte(val))

	// encoding
	code = base64.StdEncoding.EncodeToString(ciphertext)

	return
}

func blockDecrypt(code string, b cipher.Block) (val string, err error) {
	// decoding
	ciphertext := make([]byte, base64.StdEncoding.DecodedLen(len(code)))
	n, err := base64.StdEncoding.Decode(ciphertext, []byte(code))
	if err != nil {
		return
	}
	ciphertext = ciphertext[:n]

	iv := ciphertext[:b.BlockSize()]
	ciphertext = ciphertext[b.BlockSize():]
	dec := cipher.NewCTR(b, iv)
	dec.XORKeyStream(ciphertext, ciphertext)

	// remove padding
	val = string(ciphertext)
	val = strings.TrimRight(val, " ")

	return
}

// DESExtension creates an extension to decrypt the value which is encrypted with
// shared key using 3-DES algorithm.
//
// The extension returns empty string, which means "unsetting" the value, if
// decrypting failed.
//
// key must be 24 bytes string. See DESEncrypt() for implementation detail.
func DESExtension(key string) (ext Extension) {
	return genSharedKeyExtension(key, DESDecrypt)
}

// DESEncrypt encrypts the val with shared key using 3-DES/CTR algorithm.
//
//   * As needed by des.NewTripleDESCipher(), key MUST be 24 bytes string.
//   * Returned code is encoded with base64.StdEncoding.
//   * val is padded with space before encrypting.
func DESEncrypt(key, val string) (code string, err error) {
	// encryption
	b, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		return
	}

	return blockEncrypt(val, b)
}

// DESDecrypt decrypts the code with shared key using 3-DES/CTR.
//
//   * As needed by des.NewTripleDESCipher(), key MUST be 24 bytes string.
//   * code will be decoded with base64.StdEncoding before decryption.
//   * Padded space is trimmed before return.
func DESDecrypt(key, code string) (val string, err error) {
	// decryption
	b, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		return
	}

	return blockDecrypt(code, b)
}
