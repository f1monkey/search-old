package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NopFunc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var data []string
		a, err := NopFunc(nil)
		require.NoError(t, err)
		result := a(data)
		require.Equal(t, data, result)
	})
	t.Run("not empty", func(t *testing.T) {
		data := []string{"qwerty", "asdfgh"}
		a, err := NopFunc(nil)
		require.NoError(t, err)
		result := a(data)
		require.Equal(t, data, result)
	})
}
