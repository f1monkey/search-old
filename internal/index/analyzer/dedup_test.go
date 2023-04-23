package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DedupFunc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var data []string
		a, err := DedupFunc(nil)
		require.NoError(t, err)
		result := a(data)
		require.Equal(t, data, result)
	})

	t.Run("not empty", func(t *testing.T) {
		data := []string{
			"hello world",
			"hello",
			"hello world",
		}
		a, err := DedupFunc(nil)
		require.NoError(t, err)
		result := a(data)
		require.Equal(t, []string{"hello world", "hello"}, result)
	})
}
