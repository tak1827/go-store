package store

import (
	"io"
)

type Store interface {
	io.Closer

	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error

	Has(key []byte) (bool, error)
}
