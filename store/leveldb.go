package store

import (
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var (
	_ Store = (*leveldbStore)(nil)

	ErrNotFound = errors.New("not found")
)

type leveldbStore struct {
	dir string
	db  *leveldb.DB
}

func NewLevelDB(dir string) (*leveldbStore, error) {
	opts := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}

	var (
		db  *leveldb.DB
		err error
	)

	if len(dir) == 0 {
		db, err = leveldb.Open(storage.NewMemStorage(), opts)
		if err != nil {
			return nil, fmt.Errorf("failed to init leveldb: %w", err)
		}
	} else {
		db, err = leveldb.OpenFile(dir, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to init leveldb: %w", err)
		}
	}

	return &leveldbStore{
		dir: dir,
		db:  db,
	}, nil
}

func (l *leveldbStore) Close() error {
	return l.db.Close()
}

func (l *leveldbStore) Get(key []byte) ([]byte, error) {
	v, err := l.db.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("faild to get: %w", err)
	}

	return v, nil
}

func (l *leveldbStore) Put(key, value []byte) error {
	return l.db.Put(key, value, nil)
}

func (l *leveldbStore) Delete(key []byte) error {
	return l.db.Delete(key, nil)
}

func (l *leveldbStore) Has(key []byte) (bool, error) {
	return l.db.Has(key, nil)
}
