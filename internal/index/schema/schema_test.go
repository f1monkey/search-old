package schema

import (
	"testing"

	"github.com/f1monkey/search/internal/index/analyzer"
	"github.com/invopop/validation"

	"github.com/stretchr/testify/require"
)

func Test_Schema_Validate(t *testing.T) {
	t.Run("must fail if field name is empty", func(t *testing.T) {
		s := NewSchema(map[string]Field{
			"": {Type: TypeKeyword},
		}, nil)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if field type is empty", func(t *testing.T) {
		s := NewSchema(map[string]Field{
			"name": {Type: ""},
		}, nil)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if field type is invalid", func(t *testing.T) {
		s := NewSchema(map[string]Field{
			"name": {Type: "invalid"},
		}, nil)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if field type cannot have child types", func(t *testing.T) {
		s := NewSchema(
			map[string]Field{
				"name": {Type: TypeBool, Children: map[string]Field{
					"name": {},
				}},
			},
			nil,
		)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if field type must have children but there aren't any", func(t *testing.T) {
		s := NewSchema(
			map[string]Field{
				"name": {Type: TypeSlice},
			},
			nil,
		)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if field child validation fails", func(t *testing.T) {
		s := NewSchema(
			map[string]Field{
				"name": {Type: TypeSlice, Children: map[string]Field{
					"": {Type: TypeBool},
				}},
			}, nil,
		)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if text field has no analyzers", func(t *testing.T) {
		s := NewSchema(
			map[string]Field{
				"name": {Type: TypeText},
			},
			nil,
		)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if text field has unknown analyzer", func(t *testing.T) {
		s := NewSchema(
			map[string]Field{
				"name": {Type: TypeText, Analyzer: "invalid"},
			},
			nil,
		)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if analyzer has invalid type", func(t *testing.T) {
		s := NewSchema(
			map[string]Field{
				"name": {Type: TypeText, Analyzer: "analyzer"},
			},
			map[string]FieldAnalyzer{
				"analyzer": {
					Analyzers: []analyzer.Analyzer{
						{Type: "invalid", Settings: nil},
					},
				},
			},
		)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must fail if analyzer has invalid settings", func(t *testing.T) {
		s := NewSchema(
			map[string]Field{
				"name": {Type: TypeText, Analyzer: "analyzer"},
			},
			map[string]FieldAnalyzer{
				"analyzer": {
					Analyzers: []analyzer.Analyzer{
						{Type: analyzer.TokenizerRegexp, Settings: nil},
					},
				},
			},
		)
		err := validation.Validate(s)
		require.Error(t, err)
	})

	t.Run("must not fail for vaild fields", func(t *testing.T) {
		s := NewSchema(
			map[string]Field{
				"name":  {Type: TypeBool},
				"name2": {Type: TypeText, Analyzer: "analyzer"},
				"name3": {Type: TypeSlice, Children: map[string]Field{
					"name": {Type: TypeKeyword},
				}},
			},
			map[string]FieldAnalyzer{
				"analyzer": {Analyzers: []analyzer.Analyzer{
					{Type: analyzer.TokenizerRegexp, Settings: map[string]interface{}{"pattern": "\\s"}},
				}},
			},
		)
		err := validation.Validate(s)
		require.NoError(t, err)
	})
}
