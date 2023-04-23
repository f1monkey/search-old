package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TokenizerWhitespaceFunc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var data []string
		tokenizer, err := TokenizerWhitespaceFunc(nil)
		require.NoError(t, err)

		result := tokenizer(data)
		require.Equal(t, data, result)
	})

	t.Run("not empty", func(t *testing.T) {
		data := []string{
			"hello world",
			"hello  world ",
		}

		tokenizer, err := TokenizerWhitespaceFunc(nil)
		require.NoError(t, err)

		result := tokenizer(data)
		require.Equal(t, []string{"hello", "world", "hello", "world"}, result)
	})
}
