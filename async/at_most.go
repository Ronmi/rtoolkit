package async

import (
	"sync"
	"time"
)

// OnceAtMost prevents running specified function too fast
//
// It guarantees to run the function as many times as you call and no more than once
// within the duration.
func OnceAtMost(dur time.Duration, f func() error) func() error {
	lock := new(sync.Mutex)
	last := time.Now().Add(0 - dur)
	return func() error {
		lock.Lock()
		defer lock.Unlock()
		if d := time.Now().Sub(last); d <= dur {
			time.Sleep(dur - d)
		}
		last = time.Now()
		return f()
	}
}

// OnceWithin is identical to OnceAtMost, but calls within duration are ignored
func OnceWithin(dur time.Duration, f func() error) func() error {
	lock := new(sync.RWMutex)
	last := time.Now().Add(0 - dur)
	return func() error {
		lock.RLock()
		if d := time.Now().Sub(last); d <= dur {
			lock.RUnlock()
			return nil
		}
		lock.RUnlock()

		lock.Lock()
		defer lock.Unlock()
		if d := time.Now().Sub(last); d <= dur {
			return nil
		}
		last = time.Now()
		return f()
	}
}
