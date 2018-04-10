package conenv

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"reflect"
)

// RSAExtension creates an extension to decrypt the value which is encrypted with
// RSA-OAEP algorithm.
//
// You SHOULD keep both public and private key secret to prevent attacks.
//
// The extension returns empty string, which means "unsetting" the value, if
// decrypting failed.
//
// You MUST generate your key pair long enough, as described in rsa.EncryptOAEP().
// See RSAEncrypt() for implementation detail.
func RSAExtension(key *rsa.PrivateKey) (ext Extension) {
	return Extension{
		Value: func(b4 string, o Options, v reflect.Value) (af string) {
			af, err := RSADecrypt(key, b4)
			if err != nil {
				af = ""
			}

			return
		},
	}
}

// RSAEncrypt encrypts data with RSA-OAEP algorithm.
//
//   * key size matters! See rsa.EncryptOAEP() for detail.
//   * Unlike RSA is designed for, you SHOULD keep both public and private key in
//     secret.
//   * code will be decoded with base64.StdEncoding before decryption.
//   * The hash method is sha256.New(), which is documented in rsa.EncryptOAEP().
//   * The random source is crypto/rand.Reader.
func RSAEncrypt(key *rsa.PublicKey, val string) (code string, err error) {
	c, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, []byte(val), nil)
	if err != nil {
		return
	}

	code = base64.StdEncoding.EncodeToString(c)
	return
}

// RSADecrypt decrypts data with RSA-OAEP algorithm.
//
//   * key size matters! See rsa.EncryptOAEP() for detail.
//   * Unlike RSA is designed for, you SHOULD keep both public and private key in
//     secret.
//   * code will be decoded with base64.StdEncoding before decryption.
//   * The hash method is sha256.New(), which is documented in rsa.EncryptOAEP().
//   * The random source is crypto/rand.Reader.
func RSADecrypt(key *rsa.PrivateKey, code string) (val string, err error) {
	ciphertext := make([]byte, base64.StdEncoding.DecodedLen(len(code)))
	n, err := base64.StdEncoding.Decode(ciphertext, []byte(code))
	if err != nil {
		return
	}
	ciphertext = ciphertext[:n]

	ret, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, ciphertext, nil)
	if err != nil {
		return
	}

	val = string(ret)
	return
}
