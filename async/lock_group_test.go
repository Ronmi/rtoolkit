package async

import (
	"sync"
	"testing"
	"time"
)

func TestLockGroup(t *testing.T) {
	l := NewLockGroup()
	runner := func(l sync.Locker, wg *sync.WaitGroup) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()

		time.Sleep(100 * time.Millisecond)
	}

	t.Run("Different Lock", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		wg.Add(2)

		begin := time.Now().UnixNano()

		go runner(l.Locker("a"), wg)
		go runner(l.Locker("b"), wg)

		wg.Wait()
		delta := time.Duration(time.Now().UnixNano() - begin)

		if delta > 2*100*time.Millisecond {
			t.Error("running too long, is it really running concurrently?")
		}
	})

	t.Run("Same Lock", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		wg.Add(2)

		begin := time.Now().UnixNano()

		go runner(l.Locker("a"), wg)
		go runner(l.Locker("a"), wg)

		wg.Wait()
		delta := time.Duration(time.Now().UnixNano() - begin)

		if delta < 2*100*time.Millisecond {
			t.Error("running too fast, is it really locked?")
		}
	})
}
