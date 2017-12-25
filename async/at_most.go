package async

import (
	"sync"
	"time"
)

// OnceAtMost prevents running specified function too fast
//
// It guarantees:
//      - run the function as many times as you call
//      - only one call is running at the same time
//      - no more than once within the duration.
//
// Time is recorded before calling real function, which means the duration includes
// function execution time.
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

// OnceSuccessAtMost is identical to OnceAtMost, but only successful call counts
func OnceSuccessAtMost(dur time.Duration, f func() error) func() error {
	lock := new(sync.Mutex)
	last := time.Now().Add(0 - dur)
	return func() error {
		lock.Lock()
		defer lock.Unlock()
		if d := time.Now().Sub(last); d <= dur {
			time.Sleep(dur - d)
		}

		now := time.Now()
		ret := f()
		if ret == nil {
			last = now
		}
		return ret
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

// OnceSuccessWithin is identical to OnceWithin, but only success call counts
func OnceSuccessWithin(dur time.Duration, f func() error) func() error {
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

		now := time.Now()
		ret := f()
		if ret == nil {
			last = now
		}
		return ret
	}
}
