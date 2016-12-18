package store

import "testing"

func TestMemoryStore(t *testing.T) {
	s := InMemory(10)
	var id string

	t.Run("try allocating", func(t *testing.T) {
		var err error
		id, err = s.Allocate()
		if err != nil {
			t.Fatalf("cannot allocate: %s", err)
		}
	})

	t.Run("try first get", func(t *testing.T) {
		data, err := s.Get(id)
		if err != nil {
			t.Fatalf("cannot do first get: %s", err)
		}

		if data != "" {
			t.Errorf("expected data in first get to be empty string, got %s", data)
		}
	})

	lst := []string{
		"0",
		"-1",
		"1",
		"-1.1",
		"1.1",
		"asd",
		"null",
		"true",
		"false",
	}
	for _, str := range lst {
		t.Run("try set and get", func(t *testing.T) {

			if err := s.Set(id, str); err != nil {
				t.Fatalf("cannot set data in %s: %s", id, err)
			}

			data, err := s.Get(id)
			if err != nil {
				t.Fatalf("as data %s, unexpected error when get %s: %s", str, id, err)
			}

			if data != str {
				t.Errorf("expected data to be %s, got %s", str, data)
			}
		})
	}

	s.Release(id)

	t.Run("try get after release", func(t *testing.T) {
		if data, err := s.Get(id); err == nil {
			t.Errorf("expected to get error after released, got no error but data %s", data)
		}
	})

	t.Run("try set after release", func(t *testing.T) {
		if err := s.Set(id, "test"); err == nil {
			t.Error("expected error when set after released, got no error")
		}
	})
}
