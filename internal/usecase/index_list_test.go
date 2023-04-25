package usecase

import (
	"testing"

	"github.com/f1monkey/search/internal/index"
	"github.com/stretchr/testify/require"
)

func Test_IndexList_List(t *testing.T) {
	t.Run("must return index list", func(t *testing.T) {
		expected := []index.Index{{Name: "name"}}

		c := NewIndexList(func() []index.Index {
			return expected
		})

		result := c.List()
		require.Equal(t, expected, result)
	})
}
