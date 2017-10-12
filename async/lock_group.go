package async

import (
	"errors"
	"sync"
)

type lockGroupElement struct {
	counter uint
	sync.Mutex
}

// LockGroup is a container of dynamic numbers of lock
//
// LockGroup is suitable for implementing row-lock. Say you have a
// number of goroutines, each will reclculate an entry of a array.
// Locking the whole array out will make it single-threaded in fact,
// so you have to lock the entries only.
//
// LockGroup ensures only one goroutine can retrive the row-lock
// object. So multiple goroutine may run concurrently: as long as
// they use different lock object.
//
// As the underying mechanism of map works, it might not work as you
// expect when writing to map or changing size of slice/array. You
// still need a global lock when doing these jobs.
//
// To save resources, unused lock object will be cleared out.
type LockGroup struct {
	ch chan map[string]*lockGroupElement
}

func (l *LockGroup) newElement() *lockGroupElement {
	return &lockGroupElement{}
}

func (l *LockGroup) send(m map[string]*lockGroupElement) {
	go func() {
		l.ch <- m
	}()
}

// Lock locks the specified row-lock
func (l *LockGroup) Lock(key string) {
	m := <-l.ch
	e, ok := m[key]
	if !ok {
		e = l.newElement()
		m[key] = e
	}
	defer e.Lock()

	e.counter++
	l.send(m)
}

// Unlock unlocks the specified row-lock.
// It panics if it is not locked on entry to Unlock.
func (l *LockGroup) Unlock(key string) {
	m := <-l.ch
	e, ok := m[key]
	if !ok {
		panic(errors.New("unlocking non-exist lock: " + key))
	}
	defer e.Unlock()

	e.counter--
	if e.counter == 0 {
		delete(m, key)
	}

	l.send(m)
}

// Locker wraps specified row-lock to sync.Locker, so you can use it
// with other tools like sync.Cond
func (l *LockGroup) Locker(key string) sync.Locker {
	return lockInLockGroup{
		l:   l,
		key: key,
	}
}

// NewLockGroup creates a new LockGroup
func NewLockGroup() *LockGroup {
	ret := &LockGroup{
		ch: make(chan map[string]*lockGroupElement),
	}
	ret.send(map[string]*lockGroupElement{})
	return ret
}

type lockInLockGroup struct {
	l   *LockGroup
	key string
}

func (l lockInLockGroup) Lock() {
	l.l.Lock(l.key)
}
func (l lockInLockGroup) Unlock() {
	l.l.Unlock(l.key)
}
