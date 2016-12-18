// package store inplements storage space for session
package store

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
