package usecase

import (
	"fmt"
	"testing"

	"github.com/f1monkey/search/internal/index"
	"github.com/stretchr/testify/require"
)

func Test_IndexGet_Get(t *testing.T) {
	t.Run("must return error if failed to get index", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		c := NewIndexGet(func(name string) (index.Index, error) {
			return index.Index{}, expectedErr
		})

		_, err := c.Get("name")
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("must not return error if index created successfully", func(t *testing.T) {
		val := index.Index{Name: "name"}

		c := NewIndexGet(func(name string) (index.Index, error) {
			return val, nil
		})

		result, err := c.Get("name")
		require.NoError(t, err)
		require.Equal(t, val, result)
	})
}
