package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPutGetHasDelete(t *testing.T) {
	store, err := NewLevelDB("")
	require.NoError(t, err)

	key := []byte("key")
	value := []byte("value")

	err = store.Put(key, value)
	require.NoError(t, err)

	v, err := store.Get(key)
	require.NoError(t, err)
	require.EqualValues(t, v, value)

	// update value
	value = []byte("new-value")
	_ = store.Put(key, value)
	v, _ = store.Get(key)
	require.EqualValues(t, v, value)

	// key not found
	_, err = store.Get([]byte("not-found"))
	require.Error(t, err)

	has, err := store.Has(key)
	require.NoError(t, err)
	require.True(t, has)

	err = store.Delete(key)
	require.NoError(t, err)

	// duplicated delete return no error
	err = store.Delete(key)
	require.NoError(t, err)
}
