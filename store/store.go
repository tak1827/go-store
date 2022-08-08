package store

import (
	"io"
)

type Store interface {
	io.Closer

	Get(key []byte) ([]byte, error)
	List(prefix []byte) (results [][]byte, err error)

	Put(key, value []byte) error
	Delete(key []byte) error
	DeleteAll(prefix []byte) error

	Has(key []byte) (bool, error)

	Dir() string
}
