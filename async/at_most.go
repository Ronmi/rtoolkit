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
//
// Typical usecase: You have to call external API for every connected clients, but
// not more than 10 times per minute or you get banned.
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
//
// The f SHOULD be time-consuming function. If not, use RunAtLeast() instead.
//
// Typical usecase: You need to grab a web page at random time, but not more than
// 10 times per minute or you get banned.
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
