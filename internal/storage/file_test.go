package storage

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

type testData struct {
	Value string `json:"value"`
}

func Test_File_Create(t *testing.T) {
	t.Run("must return err if element already exists", func(t *testing.T) {
		s := NewFile[string, testData]()
		s.items["key"] = testData{}
		err := s.Create("key", testData{})
		require.ErrorIs(t, err, ErrAlreadyExists)
	})
	t.Run("must not return err if element does not exist", func(t *testing.T) {
		s := NewFile[string, testData]()
		err := s.Create("key", testData{})
		require.NoError(t, err)
		require.Contains(t, s.items, "key")
	})
}

func Test_File_Get(t *testing.T) {
	t.Run("must return err if element not found", func(t *testing.T) {
		s := NewFile[string, testData]()
		_, err := s.Get("key")
		require.ErrorIs(t, err, ErrNotFound)
	})
	t.Run("must return element if found", func(t *testing.T) {
		s := NewFile[string, testData]()
		element := testData{}
		s.items["key"] = element
		result, err := s.Get("key")
		require.NoError(t, err)
		require.Equal(t, element, result)
	})
}

func Test_File_Delete(t *testing.T) {
	t.Run("must return err if element not found", func(t *testing.T) {
		s := NewFile[string, testData]()
		err := s.Delete("key")
		require.ErrorIs(t, err, ErrNotFound)
	})
	t.Run("must delete element if found", func(t *testing.T) {
		s := NewFile[string, testData]()
		element := testData{}
		s.items["key"] = element
		err := s.Delete("key")
		require.NoError(t, err)
		require.Len(t, s.items, 0)
	})
}

func Test_File_Save(t *testing.T) {
	s1 := NewFile[string, testData]()
	s1.items["key"] = testData{}

	filePath := path.Join(t.TempDir(), "storage.bin")
	file, err := os.Create(filePath)
	require.NoError(t, err)
	err = s1.Save(file)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	file, err = os.Open(filePath)
	require.NoError(t, err)

	s2, err := Load[string, testData](file)
	require.NoError(t, err)

	require.EqualValues(t, s1.items, s2.items)
}
