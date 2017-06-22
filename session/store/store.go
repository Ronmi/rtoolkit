// package store implements storage space for session
package store

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	// Characters a seed value can use
	SeedChars = "abcdefghijklmnopqrstuvwxyz1234567890"

	SeedLength = 32
)

type Store interface {
	// SetTTL decides how long before data to be considered invalid (in seconds)
	SetTTL(ttl int)

	// Allocate creates a new session id, returns error if store is full
	//
	// Size of session id depends on store.
	//
	// Implementation MUST allocate space for the session id before returning it,
	// and ttl value MUST follow what was set by SetTTL().
	//
	// seed is used for session validating, see session.Manager.Start() for detail
	Allocate(seed string) (string, error)

	// Get returns session data (string), returns error if not found or something goes wrong
	//
	// It MUST refresh ttl value.
	Get(sessID string) (seed, data string, err error)

	// Set saves session data, returns error if not found or something goes wrong
	//
	// It MUST refresh ttl value.
	Set(sessID string, seed, data string) error

	// Release clears a session, never fail
	Release(sessID string)
}

// GenerateRandomKey is a helper function to generate random string as session key
//
// size is in bytes, and will be carried to multiple of 4 due to base64 encoding.
func GenerateRandomKey(size int, isExist func(string) bool) string {
	if m := size % 4; m != 0 {
		size += 4 - m
	}

	var ret string
	enc := base64.StdEncoding
	src := make([]byte, size/4*3)
	for ok := false; !ok; ok = isExist(ret) {
		// create a 64 bytes random data
		_, _ = rand.Read(src)

		ret = enc.EncodeToString(src)
	}

	return ret
}
