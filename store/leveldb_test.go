package store

import (
	"sort"
	"testing"

	// "github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestPutGetHasDelete(t *testing.T) {
	store, _ := NewLevelDB("")
	defer store.Close()

	key := []byte("key")
	value := []byte("value")

	require.NoError(t, store.Put(key, value))

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
	require.EqualError(t, err, ErrNotFound.Error())

	has, err := store.Has(key)
	require.NoError(t, err)
	require.True(t, has)

	err = store.Delete(key)
	require.NoError(t, err)

	// duplicated delete return no error
	err = store.Delete(key)
	require.NoError(t, err)
}

func TestListAllDeleteAll(t *testing.T) {
	store, _ := NewLevelDB("")
	defer store.Close()

	var (
		keys = [][]byte{
			[]byte("key1"),
			[]byte("key2"),
			[]byte("key3"),
		}
		values = [][]byte{
			[]byte("value"),
			[]byte("value2"),
			[]byte("value3"),
		}
	)

	for i := range keys {
		err := store.Put(keys[i], values[i])
		require.NoError(t, err)
	}

	// list all
	results, err := store.List(nil)
	require.NoError(t, err)
	sort.Slice(results, func(i, j int) bool {
		return string(results[i]) < string(results[j])
	})
	require.EqualValues(t, values, results)

	// list all key
	results, err = store.ListKey(nil)
	require.NoError(t, err)
	sort.Slice(results, func(i, j int) bool {
		return string(results[i]) < string(results[j])
	})
	require.EqualValues(t, keys, results)

	// delete all
	err = store.DeleteAll(nil)
	require.NoError(t, err)
	_, err = store.List(nil)
	require.Error(t, ErrNotFound, err.Error())
}
