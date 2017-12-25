package async

import (
	"errors"
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

func TestOnceSuccessAtMost(t *testing.T) {
	expect := 50 * time.Millisecond
	var e error
	f := OnceSuccessAtMost(expect, func() error {
		return e
	})

	e = nil
	f()
	begin := time.Now().UnixNano()
	e = errors.New("")
	f()
	actual := time.Duration(time.Now().UnixNano() - begin)
	if actual <= expect {
		t.Errorf("expect %dns, got %dns", expect, actual)
	}

	begin = time.Now().UnixNano()
	e = nil
	f()
	actual = time.Duration(time.Now().UnixNano() - begin)
	if actual >= expect {
		t.Errorf("expect not more than %dns, got %dns", expect, actual)
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

func TestOnceSuccessWithin(t *testing.T) {
	expect := 50 * time.Millisecond
	var e error
	a := 0
	f := OnceSuccessWithin(expect, func() error {
		if e == nil {
			a++
		}
		return e
	})

	// ok
	e = nil
	f()
	if a != 1 {
		t.Fatalf("expected run once, got %d", a)
	}

	// no
	time.Sleep(expect)
	e = errors.New("")
	f()
	if a != 1 {
		t.Fatalf("expected run once, got %d", a)
	}

	// ok
	time.Sleep(expect)
	e = nil
	f()
	if a != 2 {
		t.Fatalf("expected run twice, got %d", a)
	}

	// no
	e = errors.New("")
	f()
	if a != 2 {
		t.Fatalf("expected run twice, got %d", a)
	}

	// no
	e = nil
	f()
	if a != 2 {
		t.Fatalf("expected run twice, got %d", a)
	}
}
