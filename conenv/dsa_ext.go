package conenv

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"
	"os"
	"reflect"
	"strings"
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
				return errors.New("failed to verify signature")
			}

			return DSAVerify(key, val, code)
		},
	}
}

func encodeBigInt(i *big.Int) (code string, err error) {
	data, err := i.GobEncode()
	if err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(data), err
}

func decodeBigInt(code string) (i *big.Int, err error) {
	data := make([]byte, base64.StdEncoding.DecodedLen(len(code)))
	n, err := base64.StdEncoding.Decode(data, []byte(code))
	if err != nil {
		return
	}
	data = data[:n]

	i = &big.Int{}
	err = i.GobDecode(data)

	return
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

	rcode, err := encodeBigInt(r)
	if err != nil {
		return
	}
	scode, err := encodeBigInt(s)
	if err != nil {
		return
	}

	code = rcode + "," + scode
	return
}

// DSAVerify verifies correctness of signature generated with DSASign()
//
// It uses sha256.Sum256() to generate hashsum of val.
func DSAVerify(key *dsa.PublicKey, val string, code string) (err error) {
	codes := strings.Split(code, ",")
	if len(codes) != 2 {
		return errors.New("format error, not generated with DSASign?")
	}

	r, err := decodeBigInt(codes[0])
	if err != nil {
		return
	}
	s, err := decodeBigInt(codes[1])
	if err != nil {
		return
	}

	sum := sha256.Sum256([]byte(val))
	if !dsa.Verify(key, sum[0:], r, s) {
		err = errors.New("failed to verify dsa signature")
	}

	return
}
