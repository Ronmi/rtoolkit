package goredistore

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Ronmi/rtoolkit/session/store"
	"github.com/go-redis/redis"
)

func genSeed() string {
	return strings.Repeat(string(store.SeedChars[rand.Intn(len(store.SeedChars))]), store.SeedLength)
}

func TestStore(t *testing.T) {
	constr := os.Getenv("REDIS_CONSTR")
	if constr == "" {
		log.Fatal("You must set constr in REDIS_CONSTR.")
	}

	o, err := redis.ParseURL(constr)
	if err != nil {
		log.Fatalf("Invalid constr: %s", err)
	}

	s := &GoRedisStore{
		Options: o,
	}
	s.SetTTL(1)

	// coupled tests, stop if anything goes wrong
	seed := genSeed()
	id, err := s.Allocate(seed)
	if err != nil {
		t.Fatalf("unexpected error when allocating: %s", err)
	}

	c := s.GetClient()
	str, err := c.Get(id).Result()
	if err != nil {
		t.Fatalf("unexpected error when validating: %s", err)
	}

	if str != seed {
		t.Fatalf("expect to get only seed data `%s` got `%s`", seed, str)
	}

	aSeed, aData, err := s.Get(id)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if aSeed != seed {
		t.Fatalf("seed mismatch, expect `%s`, got `%s`", seed, aSeed)
	}
	if aData != "" {
		t.Fatalf("data mismatch, expect empty, got `%s`", aData)
	}

	// individule tests
	t.Run("Set", func(t *testing.T) {
		seed := genSeed()
		id, _ := s.Allocate(seed)
		expect := "expect"

		if err := s.Set(id, seed, expect); err != nil {
			t.Fatalf("unexpected error when setting: %s", err)
		}

		aSeed, aData, err := s.Get(id)
		if err != nil {
			t.Fatalf("unexpected error when validating: %s", err)
		}
		if aSeed != seed {
			t.Fatalf("seed mismatch, expect `%s`, got `%s`", seed, aSeed)
		}
		if aData != expect {
			t.Fatalf("data mismatch, expect `%s`, got `%s`", expect, aData)
		}
	})

	t.Run("Release", func(t *testing.T) {
		seed := genSeed()
		id, _ := s.Allocate(seed)

		s.Release(id)

		aSeed, aData, err := s.Get(id)
		if err == nil {
			t.Fatalf("sesion not released, dumping data `%s`, seed `%s`", aData, aSeed)
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		seed := genSeed()
		id, _ := s.Allocate(seed)

		time.Sleep(s.ttl)

		aSeed, aData, err := s.Get(id)
		if err == nil {
			t.Fatalf("sesion not expired, dumping data `%s`, seed `%s`", aData, aSeed)
		}
	})
}
