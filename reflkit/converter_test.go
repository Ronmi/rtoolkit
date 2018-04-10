package reflkit

import (
	"testing"
)

func TestSetOK(t *testing.T) {
	conv := DefaultStrConv()

	t.Run("bool", func(t *testing.T) {
		var b bool
		ok := conv.Set(&b, "true")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if !b {
			t.Errorf("expected to be true, got false")
		}
	})

	t.Run("int", func(t *testing.T) {
		var i int
		ok := conv.Set(&i, "-1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != -1 {
			t.Errorf("expected to be -1, got %d", i)
		}
	})
	t.Run("int8", func(t *testing.T) {
		var i int8
		ok := conv.Set(&i, "-1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != -1 {
			t.Errorf("expected to be -1, got %d", i)
		}
	})
	t.Run("int16", func(t *testing.T) {
		var i int16
		ok := conv.Set(&i, "-1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != -1 {
			t.Errorf("expected to be -1, got %d", i)
		}
	})
	t.Run("int32", func(t *testing.T) {
		var i int32
		ok := conv.Set(&i, "-1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != -1 {
			t.Errorf("expected to be -1, got %d", i)
		}
	})
	t.Run("int64", func(t *testing.T) {
		var i int64
		ok := conv.Set(&i, "-1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != -1 {
			t.Errorf("expected to be -1, got %d", i)
		}
	})

	t.Run("uint", func(t *testing.T) {
		var i uint
		ok := conv.Set(&i, "1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != 1 {
			t.Errorf("expected to be 1, got %d", i)
		}
	})
	t.Run("uint8", func(t *testing.T) {
		var i uint8
		ok := conv.Set(&i, "1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != 1 {
			t.Errorf("expected to be 1, got %d", i)
		}
	})
	t.Run("uint16", func(t *testing.T) {
		var i uint16
		ok := conv.Set(&i, "1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != 1 {
			t.Errorf("expected to be 1, got %d", i)
		}
	})
	t.Run("uint32", func(t *testing.T) {
		var i uint32
		ok := conv.Set(&i, "1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != 1 {
			t.Errorf("expected to be 1, got %d", i)
		}
	})
	t.Run("uint64", func(t *testing.T) {
		var i uint64
		ok := conv.Set(&i, "1")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if i != 1 {
			t.Errorf("expected to be 1, got %d", i)
		}
	})

	t.Run("float32", func(t *testing.T) {
		var f float32
		ok := conv.Set(&f, "1.0")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if f != 1.0 {
			t.Errorf("expected to be 1.0, got %f", f)
		}
	})
	t.Run("float64", func(t *testing.T) {
		var f float64
		ok := conv.Set(&f, "1.0")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if f != 1.0 {
			t.Errorf("expected to be 1.0, got %f", f)
		}
	})

	t.Run("string", func(t *testing.T) {
		var s string
		ok := conv.Set(&s, "abc")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if s != "abc" {
			t.Errorf("expected to be abc, got %s", s)
		}
	})

	t.Run("byte-slice", func(t *testing.T) {
		var s []byte
		ok := conv.Set(&s, "abc")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if string(s) != "abc" {
			t.Errorf("expected to be abc, got %s", string(s))
		}
	})
	t.Run("rune-slice", func(t *testing.T) {
		var s []rune
		ok := conv.Set(&s, "abc")
		if !ok {
			t.Errorf("expected to be ok, but actually not")
		}
		if string(s) != "abc" {
			t.Errorf("expected to be abc, got %s", string(s))
		}
	})
}

func TestSetUnsupport(t *testing.T) {
	conv := DefaultStrConv()

	var x []int
	if ok := conv.Set(&x, "1"); ok {
		t.Errorf("expected not to be ok, but actually passed")
	}
}
