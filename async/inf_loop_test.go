package async

import (
	"context"
	"testing"
	"time"
)

type errFailed string

func (e errFailed) Error() string { return string(e) }

func TestInfiniteLoop(t *testing.T) {
	cnt := 0
	task := func() error {
		if cnt >= 3 {
			return errFailed("jos done before cancel")
		}

		cnt++
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	cancel, errchan := InfiniteLoop(task)
	time.Sleep(10 * time.Millisecond)
	cancel()

	if err := <-errchan; err != context.Canceled {
		t.Error(err)
	}
}

func TestInfiniteLoopCancelAfterDone(t *testing.T) {
	cnt := 0
	task := func() error {
		if cnt >= 3 {
			return errFailed("jos done before cancel")
		}

		cnt++
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	cancel, errchan := InfiniteLoop(task)
	time.Sleep(35 * time.Millisecond)
	cancel()

	if err := <-errchan; err == context.Canceled {
		t.Error("cancel should not work after job done")
	}
}
