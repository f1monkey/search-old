package index

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Storage_Create(t *testing.T) {
	t.Run("must return err if index already exists", func(t *testing.T) {
		s := NewStorage()
		s.items["name"] = Index{}
		err := s.Create(Index{Name: "name"})
		require.ErrorIs(t, err, ErrAlreadyExists)
	})
	t.Run("must not return err if index does not exist", func(t *testing.T) {
		s := NewStorage()
		err := s.Create(Index{Name: "name"})
		require.NoError(t, err)
		require.Contains(t, s.items, "name")
	})
}

func Test_Storage_Get(t *testing.T) {
	t.Run("must return err if index not found", func(t *testing.T) {
		s := NewStorage()
		_, err := s.Get("name")
		require.ErrorIs(t, err, ErrNotFound)
	})
	t.Run("must return index if found", func(t *testing.T) {
		s := NewStorage()
		index := Index{Name: "name"}
		s.items[index.Name] = index
		result, err := s.Get("name")
		require.NoError(t, err)
		require.Equal(t, index, result)
	})
}

func Test_Storage_Delete(t *testing.T) {
	t.Run("must return err if index not found", func(t *testing.T) {
		s := NewStorage()
		err := s.Delete("name")
		require.ErrorIs(t, err, ErrNotFound)
	})
	t.Run("must delete index if found", func(t *testing.T) {
		s := NewStorage()
		index := Index{Name: "name"}
		s.items[index.Name] = index
		err := s.Delete("name")
		require.NoError(t, err)
		require.Len(t, s.items, 0)
	})
}

func Test_Storage_Save(t *testing.T) {
	s1 := NewStorage()
	s1.items["name"] = Index{
		Name: "name",
	}

	filePath := path.Join(t.TempDir(), "storage.bin")
	file, err := os.Create(filePath)
	require.NoError(t, err)
	err = s1.Save(file)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	file, err = os.Open(filePath)
	require.NoError(t, err)

	s2, err := Load(file)
	require.NoError(t, err)

	require.EqualValues(t, s1.items, s2.items)
}
