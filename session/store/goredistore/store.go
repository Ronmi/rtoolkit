// package goredistore implements Redis based session store using go-redis/redis
//
// Since seed length is fixed, storing data is simplified to only one command:
//
//      redis.Set(sessionID, seed+data, ttl)
//
// But reading data need two command, one for refreshing ttl, one for getting
// data.
package goredistore

import (
	"errors"
	"sync"
	"time"

	"github.com/Ronmi/rtoolkit/session/store"
	"github.com/go-redis/redis"
)

type GoRedisStore struct {
	*redis.Options
	client *redis.Client
	lock   sync.Mutex

	ttl time.Duration
}

// GetClient returns go-redis client instance
func (s *GoRedisStore) GetClient() *redis.Client {
	if s.client != nil {
		return s.client
	}

	// ensure only one thread can create connection
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.client == nil {
		s.client = redis.NewClient(s.Options)
	}
	return s.client
}

func (s *GoRedisStore) SetTTL(ttl int) {
	s.ttl = time.Duration(ttl) * time.Second
}

func (s *GoRedisStore) Allocate(seed string) (string, error) {
	c := s.GetClient()
	var err error
	id := store.GenerateRandomKey(32, func(id string) bool {
		var ret bool
		ret, err = c.SetNX(id, seed, s.ttl).Result()
		if err != nil {
			return true
		}

		return ret
	})

	return id, err
}

func (s *GoRedisStore) Get(sessID string) (seed, data string, err error) {
	c := s.GetClient()
	if _, err := c.Expire(sessID, s.ttl).Result(); err != nil {
		return "", "", err
	}

	str, err := c.Get(sessID).Result()
	if err != nil {
		return "", "", err
	}

	if len(str) < store.SeedLength {
		return "", "", errors.New("rtoolkit/session/store/goredis: incorrect format data detected in store")
	}

	return str[:store.SeedLength], str[store.SeedLength:], err
}

func (s *GoRedisStore) Set(sessID, seed, data string) error {
	c := s.GetClient()

	_, err := c.SetXX(sessID, seed+data, s.ttl).Result()
	return err
}

func (s *GoRedisStore) Release(sessID string) {
	s.GetClient().Expire(sessID, time.Duration(-1)*time.Second)

}
