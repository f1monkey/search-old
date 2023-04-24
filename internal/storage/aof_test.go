package storage

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AOF_Create(t *testing.T) {
	f := path.Join(t.TempDir(), "tmp.dat")

	storage, err := NewAOFFromPath[string, testData](f)
	require.NoError(t, err)

	t.Run("must write an element to the file on create", func(t *testing.T) {
		key := "key"
		value := testData{"value"}
		require.NoError(t, storage.Create(key, value))

		data, err := os.ReadFile(f)
		require.NoError(t, err)
		require.Equal(t, "{\"key\":\"key\",\"value\":{\"value\":\"value\"}}\n", string(data))
		require.Contains(t, storage.items, key)
		require.Equal(t, storage.items[key], value)
	})

	t.Run("must append next element to file", func(t *testing.T) {
		key := "key2"
		value := testData{"value2"}
		require.NoError(t, storage.Create(key, value))

		data, err := os.ReadFile(f)
		require.NoError(t, err)
		require.Equal(t, `{"key":"key","value":{"value":"value"}}
{"key":"key2","value":{"value":"value2"}}
`,
			string(data),
		)
		require.Contains(t, storage.items, key)
		require.Equal(t, storage.items[key], value)
	})

	t.Run("must return error if element already exists", func(t *testing.T) {
		key := "key"
		value := testData{"value3"}
		require.ErrorIs(t, storage.Create(key, value), ErrAlreadyExists)
	})
}

func Test_AOF_Get(t *testing.T) {
	t.Run("must return err if element not found", func(t *testing.T) {
		f := path.Join(t.TempDir(), "tmp.dat")
		s, err := NewAOFFromPath[string, testData](f)
		require.NoError(t, err)

		_, err = s.Get("key")
		require.ErrorIs(t, err, ErrNotFound)
	})
	t.Run("must return element if found", func(t *testing.T) {
		f := path.Join(t.TempDir(), "tmp.dat")
		s, err := NewAOFFromPath[string, testData](f)
		require.NoError(t, err)
		element := testData{}
		s.items["key"] = element
		result, err := s.Get("key")
		require.NoError(t, err)
		require.Equal(t, element, result)
	})
}

func Test_AOF_Delete(t *testing.T) {
	t.Run("must return err if element not found", func(t *testing.T) {
		f := path.Join(t.TempDir(), "tmp.dat")
		s, err := NewAOFFromPath[string, testData](f)
		require.NoError(t, err)
		err = s.Delete("key")
		require.ErrorIs(t, err, ErrNotFound)
	})
	t.Run("must delete element if found", func(t *testing.T) {
		f := path.Join(t.TempDir(), "tmp.dat")
		s, err := NewAOFFromPath[string, *testData](f)
		require.NoError(t, err)

		key := "key"
		value := &testData{"value"}
		require.NoError(t, s.Create(key, value))

		err = s.Delete("key")
		require.NoError(t, err)
		data, err := os.ReadFile(f)
		require.NoError(t, err)
		require.Equal(t, `{"key":"key","value":{"value":"value"}}
{"key":"key","isDeleted":true}
`,
			string(data),
		)

		require.Len(t, s.items, 0, "must delete element from map")
	})
}
