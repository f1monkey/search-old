package usecase

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IndexDelete_Delete(t *testing.T) {
	t.Run("must return error if failed to delete index", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		c := NewIndexDelete(nil, func(name string) error {
			return expectedErr
		})

		err := c.Delete("name")
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("must not return error if index created successfully", func(t *testing.T) {
		c := NewIndexDelete(nil, func(name string) error {
			return nil
		})

		err := c.Delete("name")
		require.NoError(t, err)
	})
}
