package store

import (
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
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

func (l *leveldbStore) List(prefix []byte) (results [][]byte, err error) {
	// NOTE:
	// ReadOptions holds the optional parameters for 'read operation'. The
	// 'read operation' includes Get, Find and NewIterator.
	// type ReadOptions struct {
	// 	// DontFillCache defines whether block reads for this 'read operation'
	// 	// should be cached. If false then the block will be cached. This does
	// 	// not affects already cached block.
	// 	//
	// 	// The default value is false.
	// 	DontFillCache bool

	// 	// Strict will be OR'ed with global DB 'strict level' unless StrictOverride
	// 	// is present. Currently only StrictReader that has effect here.
	// 	Strict Strict
	// }
	var ro *opt.ReadOptions
	iter := l.db.NewIterator(util.BytesPrefix(prefix), ro)
	defer iter.Release()

	for iter.Next() {
		var (
			val    = iter.Value()
			result = make([]byte, len(val))
		)
		copy(result, val)

		results = append(results, result)
	}

	if err = iter.Error(); err != nil {
		err = fmt.Errorf("faild to iter: %w", err)
		return
	}

	if len(results) == 0 {
		err = ErrNotFound
		return
	}

	return
}

func (l *leveldbStore) ListKey(prefix []byte) (results [][]byte, err error) {
	var ro *opt.ReadOptions
	iter := l.db.NewIterator(util.BytesPrefix(prefix), ro)
	defer iter.Release()

	for iter.Next() {
		var (
			key    = iter.Key()
			result = make([]byte, len(key))
		)
		copy(result, key)

		results = append(results, result)
	}

	if err = iter.Error(); err != nil {
		err = fmt.Errorf("faild to iter: %w", err)
		return
	}

	if len(results) == 0 {
		err = ErrNotFound
		return
	}

	return
}

func (l *leveldbStore) Put(key, value []byte) error {
	return l.db.Put(key, value, nil)
}

func (l *leveldbStore) Delete(key []byte) error {
	return l.db.Delete(key, nil)
}

func (l *leveldbStore) DeleteAll(prefix []byte) error {
	var (
		ro    *opt.ReadOptions
		wo    *opt.WriteOptions
		batch = new(leveldb.Batch)
	)
	iter := l.db.NewIterator(util.BytesPrefix(prefix), ro)
	defer iter.Release()

	for iter.Next() {
		batch.Delete(iter.Key())
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("faild to iter: %w", err)
	}

	if err := l.db.Write(batch, wo); err != nil {
		return fmt.Errorf("faild to write delete batch: %w", err)
	}

	return nil
}

func (l *leveldbStore) Has(key []byte) (bool, error) {
	return l.db.Has(key, nil)
}

func (l *leveldbStore) Dir() string {
	return l.dir
}
