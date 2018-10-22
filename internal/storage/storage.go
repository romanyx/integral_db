package storage

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrNotFound returns when a key is not
	// found in the storage.
	ErrNotFound = errors.New("not found")
)

// Storage represents abstraction
// for key/value strorage.
type Storage interface {
	Set(ctx context.Context, key, value interface{})
	Get(key interface{}) (value interface{}, err error)
}

// New returns initialzied storage
// implementation. When context
// of the Set method will be done, storage
// will remove key from it. When same key
// will be passed to the Get method, previous
// context done wait will be expired.
func New() Storage {
	m := muxMap{
		Mutex:   &sync.Mutex{},
		storage: make(map[interface{}]data),
	}

	return &m
}

type muxMap struct {
	*sync.Mutex
	storage map[interface{}]data
}

type data struct {
	value interface{}
	reset chan struct{}
}

func (m muxMap) Get(key interface{}) (interface{}, error) {
	m.Lock()
	defer m.Unlock()

	data, ok := m.storage[key]
	if !ok {
		return nil, ErrNotFound
	}

	close(data.reset)
	delete(m.storage, key)

	return data.value, nil
}

func (m muxMap) Set(ctx context.Context, key, value interface{}) {
	m.Lock()
	{

		d, ok := m.storage[key]
		if ok {
			// signal reset channel to
			// prevent key deletion on
			// <- ctx.Done().
			close(d.reset)
		}

		reset := make(chan struct{})

		m.storage[key] = data{
			value: value,
			reset: reset,
		}

		go func() {
			select {
			case <-ctx.Done():
				delete(m.storage, key)
				ctxDoneCall()
			case <-reset:
				return
			}
		}()
	}
	m.Unlock()
}

var ctxDoneCall = func() {}
