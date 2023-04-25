package index

import (
	"testing"

	"github.com/invopop/validation"
	"github.com/stretchr/testify/require"
)

func Test_Index_Validate(t *testing.T) {
	t.Run("must fail if index name is empty", func(t *testing.T) {
		i := Index{}
		err := validation.Validate(i)
		require.Error(t, err)
	})

	t.Run("must fail if schema is omitted", func(t *testing.T) {
		i := Index{Name: "name"}
		err := validation.Validate(i)
		require.Error(t, err)
	})

	t.Run("must fail if schema is invalid", func(t *testing.T) {
		i := Index{Name: "name", Schema: Schema{Fields: nil}}
		err := validation.Validate(i)
		require.Error(t, err)
	})

	t.Run("must not fail if index is valid", func(t *testing.T) {
		i := Index{
			Name: "name",
			Schema: Schema{
				Fields: map[string]Field{
					"field": {
						Type:     TypeBool,
						Required: false,
					},
				},
			},
		}

		err := validation.Validate(i)
		require.NoError(t, err)
	})
}
