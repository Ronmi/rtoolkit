package async

import "time"

// OnceAtMost prevents running specified function too fast
func OnceAtMost(dur time.Duration, f func() error) func() error {
	return func() error {
		begin := time.Now().UnixNano()
		ret := f()
		dur -= time.Duration(time.Now().UnixNano() - begin)
		if dur > 0 {
			time.Sleep(dur)
		}

		return ret
	}
}
