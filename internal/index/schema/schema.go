package schema

import (
	"context"
	"fmt"
	"regexp"

	"github.com/f1monkey/search/internal/index/analyzer"
	"github.com/invopop/validation"
)

type FieldAnalyzer struct {
	Analyzers []analyzer.Analyzer `json:"analyzers"`
}

func (fa FieldAnalyzer) Build() (analyzer.Func, error) {
	return analyzer.Chain(fa.Analyzers)
}

func (a FieldAnalyzer) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Analyzers, validation.Required, validation.Length(1, 0)),
	)
}

type Schema struct {
	Analyzers map[string]FieldAnalyzer `json:"analyzers"`
	Fields    map[string]Field         `json:"fields"`
}

func NewSchema(fields map[string]Field, analyzers map[string]FieldAnalyzer) Schema {
	return Schema{
		Fields:    fields,
		Analyzers: analyzers,
	}
}

func (s Schema) ValidateDoc(doc map[string]interface{}) error {
	return ValidateDoc(s, doc)
}

func (s Schema) Validate() error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeySchema, s)

	if err := validateKeys("fields", s.Fields); err != nil {
		return err
	}

	if err := validateKeys("fields", s.Analyzers); err != nil {
		return err
	}

	return validation.ValidateStructWithContext(ctx, &s,
		validation.Field(&s.Fields, validation.Required),
		validation.Field(&s.Analyzers),
	)
}

var keyPattern = "[A-z0-9]+"
var keyRegex = regexp.MustCompile(keyPattern)

func validateKeys[T any](name string, m map[string]T) error {
	if _, ok := m[""]; ok {
		return validation.Errors{name: validation.NewError("", "empty keys not allowed")}
	}

	for k := range m {
		if !keyRegex.MatchString(k) {
			return validation.Errors{name: validation.NewError(k, fmt.Sprintf("key %q does not match %q", k, keyPattern))}
		}
	}

	return nil
}
