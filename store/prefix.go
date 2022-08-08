package store

import (
	"sync"
)

type PrefixStore struct {
	sync.Mutex

	prefix []byte
	store  Store
}

func NewPrefixStore(store Store, prefix []byte) *PrefixStore {
	return &PrefixStore{
		prefix: prefix,
		store:  store,
	}
}

func (s *PrefixStore) Close() error {
	s.Lock()
	defer s.Unlock()

	return s.store.Close()
}

func (s *PrefixStore) Get(key []byte) ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	return s.store.Get(s.prefixed(key))
}

func (s *PrefixStore) ListAll() ([][]byte, error) {
	s.Lock()
	defer s.Unlock()

	return s.store.List(s.prefix)
}

func (s *PrefixStore) ListAllKey() ([][]byte, error) {
	s.Lock()
	defer s.Unlock()

	return s.store.ListKey(s.prefix)
}

func (s *PrefixStore) Put(key []byte, value []byte) error {
	s.Lock()
	defer s.Unlock()

	return s.store.Put(s.prefixed(key), value)
}

func (s *PrefixStore) Delete(key []byte) error {
	s.Lock()
	defer s.Unlock()

	return s.store.Delete(s.prefixed(key))
}

func (s *PrefixStore) DeleteAll() error {
	s.Lock()
	defer s.Unlock()

	return s.store.DeleteAll(s.prefix)
}

func (s *PrefixStore) Has(key []byte) (bool, error) {
	s.Lock()
	defer s.Unlock()

	return s.store.Has(s.prefixed(key))
}

func (s *PrefixStore) prefixed(key []byte) []byte {
	return append(s.prefix, key...)
}
