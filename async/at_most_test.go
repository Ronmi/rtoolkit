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
		f()
		begin := time.Now().UnixNano()
		f()
		actual := time.Duration(time.Now().UnixNano() - begin)
		if actual < expect {
			t.Errorf("expect %dns, got %dns", expect, actual)
		}
	}
}

func TestOnceWithin(t *testing.T) {
	expect := 50 * time.Millisecond
	a := 0
	funcs := []struct {
		expect int
		f      func() error
	}{
		{
			expect: 2,
			f: OnceWithin(expect, func() error {
				a++
				return nil
			}),
		},
		{
			expect: 4,
			f: OnceWithin(expect, func() error {
				a++
				time.Sleep(100 * time.Millisecond)
				return nil
			}),
		},
	}

	for _, c := range funcs {
		a = 0
		c.f()
		c.f()
		time.Sleep(expect)
		c.f()
		c.f()

		if a != c.expect {
			t.Errorf("expected %d, got %d", c.expect, a)
		}
	}
}
