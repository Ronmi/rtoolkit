package conenv

import (
	"bytes"
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"math/big"
	"os"
	"reflect"
)

// DSAExtension creates an extension to verify value corectness using DSA digital
// signature.
//
// It looks up for signature from envariable variable has "_SIGN" postfix to
// original name: if the variable loads its value from "DB_PASS", signature will be
// loaded from "DB_PASS_SIGN".
//
// As it runs on validating stage, you MUST generate signature against original
// value, not encrypted one.
//
// See DSASign() and DSAVerify() for implementation detail.
func DSAExtension(key *dsa.PublicKey) (ext Extension) {
	return Extension{
		Validate: func(
			o Options, v reflect.Value, name, val string,
		) (err error) {
			code, ok := os.LookupEnv(name + "_SIGN")
			if !ok {
				return errors.New(
					"failed to verify signature of " + name)
			}

			return DSAVerify(key, val, code)
		},
	}
}

// DSASign generates digital signature on val
//
// It uses sha256.Sum256() to generate hashsum of val. The generated signature will
// be encoded in Gob and base64.
func DSASign(key *dsa.PrivateKey, val string) (code string, err error) {
	sum := sha256.Sum256([]byte(val))
	r, s, err := dsa.Sign(rand.Reader, key, sum[0:])
	if err != nil {
		return
	}

	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err = enc.Encode(r); err != nil {
		return
	}
	if err = enc.Encode(s); err != nil {
		return
	}

	code = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}

// DSAVerify verifies correctness of signature generated with DSASign()
//
// It uses sha256.Sum256() to generate hashsum of val.
func DSAVerify(key *dsa.PublicKey, val string, code string) (err error) {
	data := make([]byte, base64.StdEncoding.DecodedLen(len(code)))
	n, err := base64.StdEncoding.Decode(data, []byte(code))
	if err != nil {
		return
	}
	data = data[:n]

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	r := &big.Int{}
	if err = dec.Decode(r); err != nil {
		return
	}
	s := &big.Int{}
	if err = dec.Decode(s); err != nil {
		return
	}

	sum := sha256.Sum256([]byte(val))
	if !dsa.Verify(key, sum[0:], r, s) {
		err = errors.New("failed to verify dsa signature")
	}

	return
}
