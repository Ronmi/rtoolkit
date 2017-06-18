package store

import (
	"errors"
	"sync"
	"time"
)

type memoryElement struct {
	*sync.Mutex // protects only data, lastUsed is not protected
	data        string
	lastUsed    int64
}

func newMemEle() *memoryElement {
	return &memoryElement{
		&sync.Mutex{},
		"",
		time.Now().UnixNano(),
	}
}

func (e *memoryElement) isValid(ttl int64) bool {
	return time.Now().UnixNano() <= ttl+e.lastUsed
}

func (e *memoryElement) invalid(ttl int64) {
	// sould be now - ttl - 1, but set zero is faster
	e.lastUsed = 0
}

func (e *memoryElement) get() string {
	e.lastUsed = time.Now().UnixNano()
	return e.data
}

func (e *memoryElement) set(data string) {
	e.lastUsed = time.Now().UnixNano()
	e.data = data
}

type memoryStore struct {
	data   map[string]*memoryElement
	ttl    int64
	lock   *sync.Mutex // for allocate/release/gc
	gcing  bool
	lastgc int64
}

func (s *memoryStore) SetTTL(ttl int) {
	s.ttl = int64(ttl) * int64(time.Second)
}

func (s *memoryStore) Allocate() (id string, err error) {
	go s.gc() // run gc

	id = GenerateRandomKey(32, func(id string) bool {
		s.lock.Lock()
		defer s.lock.Unlock()

		_, ok := s.data[id]
		if !ok {
			s.data[id] = newMemEle()
		}

		return !ok
	})

	return
}

func (s *memoryStore) Release(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.doRelease(id)
}

func (s *memoryStore) doRelease(id string) {
	if data, ok := s.data[id]; ok {
		data.Lock()
		defer data.Unlock()
		data.invalid(s.ttl)

		delete(s.data, id)
	}
}

func (s *memoryStore) canGC() (ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	now := time.Now().UnixNano()
	if s.gcing {
		return
	}

	// at most 1 times in 1 sec
	if s.lastgc+int64(time.Second) >= now {
		return
	}

	s.gcing = true
	s.lastgc = now
	return true
}

// gc clears all invalid entries using Release, so no lock is required
func (s *memoryStore) gc() {
	if !s.canGC() {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for id, ele := range s.data {
		if ele.isValid(s.ttl) {
			continue
		}

		s.doRelease(id)
	}

	s.gcing = false
}

func (s *memoryStore) getElement(id string) (*memoryElement, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	e, ok := s.data[id]
	if !ok {
		return nil, errors.New("session not exists: " + id)
	}
	return e, nil
}

func (s *memoryStore) Get(id string) (data string, err error) {
	e, err := s.getElement(id)
	if err != nil {
		return
	}

	e.Lock()
	defer e.Unlock()
	if !e.isValid(s.ttl) {
		return "", errors.New("session expired: " + id)
	}

	return e.get(), nil
}

func (s *memoryStore) Set(id string, data string) error {
	e, err := s.getElement(id)
	if err != nil {
		return err
	}

	e.Lock()
	defer e.Unlock()
	if !e.isValid(s.ttl) {
		return errors.New("session expired: " + id)
	}
	e.set(data)

	return nil
}

// InMemory creates a memory store
func InMemory(ttl int) Store {
	ret := &memoryStore{
		data: make(map[string]*memoryElement),
		lock: &sync.Mutex{},
	}
	ret.SetTTL(ttl)
	return ret
}
