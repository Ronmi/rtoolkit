package apitool

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

// TOTPMiddleware represents a middleware to auth client with TOTP before execution
//
// The TOTP implementation is this middleware has following restrictions:
//
//   - Accepts only two windows: current and previous one.
type TOTPMiddleware struct {
	// 80bits binary secret, REQUEIRED
	Secret [10]byte
	// how many digits should a code be, at least 6
	//
	// any value less than 6 will be forced to 6
	Digit int
	// how to get code from request, leave nil to use default implementation,
	// which loads code from X-OTP-CODE header
	//
	// for best security, it is not suggested to load code from URL query
	GetCode func(r *http.Request) string
	// which kinds of error is returned if failed to auth, leave nil to use
	// default implementation, which returns 403 error
	Failed func(r *http.Request) error
}

// HOTP implements RFC4226
//
// The signature of this function is designed for use within TOTPMiddleware.
func (m TOTPMiddleware) HOTP(c int64) (code string) {
	// safe to set struct member as it is passed by value
	if m.Digit < 6 {
		m.Digit = 6
	}

	// https://tools.ietf.org/html/rfc4226#section-5.3
	// HS = HMAC-SHA-1(K,C)
	hash := hmac.New(sha1.New, m.Secret[:])
	if err := binary.Write(hash, binary.BigEndian, c); err != nil {
		return
	}
	hs := hash.Sum(nil)

	// Sbits = DT(HS)
	// Snum  = StToNum(Sbits)
	offset := hs[19] & 0x0f
	snum := binary.BigEndian.Uint32(hs[offset : offset+4])
	snum &= 0x7fffffff

	// Return D = Snum mod 10^Digit
	code = strconv.Itoa(int(snum))
	l := len(code)
	if l < m.Digit {
		return strings.Repeat("0", m.Digit-l) + code
	}
	return code[l-m.Digit:]
}

func (m TOTPMiddleware) verify(code string) (ok bool) {
	frame := time.Now().Unix() / 30
	if m.HOTP(frame) == code {
		return true
	}

	return m.HOTP(frame-1) == code
}

func (m TOTPMiddleware) Middleware(h jsonapi.Handler) (ret jsonapi.Handler) {
	return func(r jsonapi.Request) (data interface{}, err error) {
		// safe to set struct member as it is passed by value
		if m.GetCode == nil {
			m.GetCode = OTPCodeByHeader("X-OTP-CODE")
		}
		if m.Failed == nil {
			m.Failed = DefaultOTPFailHandler
		}

		if !m.verify(m.GetCode(r.R())) {
			return nil, m.Failed(r.R())
		}

		return h(r)
	}
}

var E403TOTP = jsonapi.E403.SetData("failed to auth with TOTP")

// DefaultOTPFailHandler is the default implementation for TOTP failure handler
//
// It just returns E403TOTP
func DefaultOTPFailHandler(r *http.Request) (err error) {
	return E403TOTP
}

// OTPCodeByHeader grabs TOTP value from custom header
func OTPCodeByHeader(key string) func(*http.Request) string {
	return func(r *http.Request) (code string) {
		return r.Header.Get(key)
	}
}

// OTPCodeByHeader grabs TOTP value from post form
func OTPCodeByForm(key string) func(*http.Request) string {
	return func(r *http.Request) (code string) {
		return r.PostFormValue(key)
	}
}

// OTPInHeader is a helper function to create OTPMiddleware
//
//   - OTP code is passed in custom HTTP header specified by headerKey.
//   - secret is 10 bytes binary string, DO NOT USE PLAIN TEXT for best security.
//
// It is identical to the following code, which is also actual implementation:
//
//     return (TOTPMiddleware{
//         Secret: secret,
//         GetCode: OTPCodeByHeader(headerKey),
//         Failed: DefaultOTPFailHandler,
//     }).Middleware
func TOTPInHeader(secret [10]byte, headerKey string) jsonapi.Middleware {
	return (TOTPMiddleware{
		Secret:  secret,
		GetCode: OTPCodeByHeader(headerKey),
		Failed:  DefaultOTPFailHandler,
	}).Middleware
}

// OTPInForm is a helper function to create OTPMiddleware
//
//   - OTP code is passed in post form specified in formKey
//   - secret is 10 bytes binary string, DO NOT USE PLAIN TEXT for best security.
//
// It is identical to the following code, which is also actual implementation:
//
//     return (TOTPMiddleware{
//         Secret: secret,
//         GetCode: OTPCodeByForm(formKey),
//         Failed: DefaultOTPFailHandler,
//     }).Middleware
func TOTPInForm(secret [10]byte, formKey string) jsonapi.Middleware {
	return (TOTPMiddleware{
		Secret:  secret,
		GetCode: OTPCodeByForm(formKey),
		Failed:  DefaultOTPFailHandler,
	}).Middleware
}
