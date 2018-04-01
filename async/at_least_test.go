package async

import (
	"errors"
	"testing"
	"time"
)

func TestRunAtLeast(t *testing.T) {
	expect := 50 * time.Millisecond
	funcs := []func() error{
		RunAtLeast(expect, func() error {
			return nil
		}),
		RunAtLeast(expect, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}),
	}

	for _, f := range funcs {
		f()
		begin := time.Now().UnixNano()
		f()
		actual := time.Duration(time.Now().UnixNano() - begin)
		if actual < expect {
			t.Errorf("expect %dns, got %dns", expect, actual)
		}
	}
}

func TestRunCondAtLeast(t *testing.T) {
	expect := 50 * time.Millisecond
	cases := []struct {
		f     func() error
		count bool
	}{
		{RunSuccessAtLeast(expect, func() error {
			return nil
		}), true},
		{RunSuccessAtLeast(expect, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}), true},
		{RunSuccessAtLeast(expect, func() error {
			return errors.New("")
		}), false},
		{RunSuccessAtLeast(expect, func() error {
			time.Sleep(100 * time.Millisecond)
			return errors.New("")
		}), true},
		{RunFailedAtLeast(expect, func() error {
			return nil
		}), false},
		{RunFailedAtLeast(expect, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}), true},
		{RunFailedAtLeast(expect, func() error {
			return errors.New("")
		}), true},
		{RunFailedAtLeast(expect, func() error {
			time.Sleep(100 * time.Millisecond)
			return errors.New("")
		}), true},
	}

	for _, c := range cases {
		c.f()
		begin := time.Now().UnixNano()
		c.f()
		actual := time.Duration(time.Now().UnixNano() - begin)
		if c.count && actual < expect {
			t.Errorf("expect %dns, got %dns", expect, actual)
		}
		if !c.count && actual > expect {
			t.Errorf("expect %dns, got %dns", expect, actual)
		}
	}
}
