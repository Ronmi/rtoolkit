package async

import "time"

// RunAtLeast ensures the execution time is greater than the duration
//
// Typical usecase: You have to grab a web page repeatly, but not more than 10 times
// per minute or you get banned.
func RunAtLeast(dur time.Duration, f func() error) func() error {
	return func() (err error) {
		begin := time.Now()
		err = f()
		if d := time.Now().Sub(begin); d <= dur {
			time.Sleep(dur - d)
		}
		return
	}
}

// RunSuccessAtLeast is identical to RunAtLeast, but only successful call counts
func RunSuccessAtLeast(dur time.Duration, f func() error) func() error {
	return func() (err error) {
		begin := time.Now()
		err = f()
		if d := time.Now().Sub(begin); err == nil && d <= dur {
			time.Sleep(dur - d)
		}
		return
	}
}

// RunFailedAtLeast is identical to RunAtLeast, but only failed call counts
func RunFailedAtLeast(dur time.Duration, f func() error) func() error {
	return func() (err error) {
		begin := time.Now()
		err = f()
		if d := time.Now().Sub(begin); err != nil && d <= dur {
			time.Sleep(dur - d)
		}
		return
	}
}
