package async

import (
	"testing"
	"time"
)

func TestOnceAtMost(t *testing.T) {
	expect := 50 * time.Millisecond
	funcs := []func() error{
		OnceAtMost(expect, func() error {
			return nil
		}),
		OnceAtMost(expect, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}),
	}

	for _, f := range funcs {
		begin := time.Now().UnixNano()
		f()
		actual := time.Duration(time.Now().UnixNano() - begin)
		if actual < expect {
			t.Errorf("expect %dns, got %dns", expect, actual)
		}
	}
}
