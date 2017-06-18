// package store implements storage space for session
package store

import (
	"crypto/rand"
	"encoding/base64"
)

type Store interface {
	// SetTTL decides how long before data to be considered invalid (in seconds)
	SetTTL(ttl int)

	// Allocate creates a new session id, returns error if store is full
	Allocate() (string, error)

	// Get returns session data (string), returns error if not found
	Get(sessID string) (string, error)

	// Set saves session data, returns error if not found or store is full
	Set(sessID string, data string) error

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
