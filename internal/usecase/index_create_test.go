package usecase

import (
	"fmt"
	"testing"

	"github.com/f1monkey/search/internal/index"
	"github.com/f1monkey/search/internal/index/schema"
	"github.com/stretchr/testify/require"
)

func Test_IndexCreate_Create(t *testing.T) {
	validIndex := index.Index{
		Name: "name",
		Schema: schema.Schema{
			Fields: map[string]schema.Field{
				"field": {
					Type:     schema.TypeBool,
					Required: false,
				},
			},
		},
	}

	t.Run("must return error if entity validation fails", func(t *testing.T) {
		c := NewIndexCreate(nil, func(name string, index index.Index) error {
			return nil
		})

		err := c.Create(index.Index{})
		require.Error(t, err)
	})

	t.Run("must return error if failed to create index", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		c := NewIndexCreate(nil, func(name string, index index.Index) error {
			return expectedErr
		})

		err := c.Create(validIndex)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("must not return error if index created successfully", func(t *testing.T) {
		c := NewIndexCreate(nil, func(name string, index index.Index) error {
			return nil
		})

		err := c.Create(validIndex)
		require.NoError(t, err)
	})
}
