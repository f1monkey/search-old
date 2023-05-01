package schema

import (
	"context"

	"github.com/f1monkey/search/pkg/errs"
	"github.com/invopop/validation"
)

type Type string

const (
	TypeAll  Type = "all"
	TypeBool Type = "bool"

	// String types
	TypeKeyword Type = "keyword"
	TypeText    Type = "text"

	TypeSlice Type = "slice"
	TypeMap   Type = "map"

	// Integer types
	TypeUnsignedLong Type = "unsigned_long" // unsigned int64
	TypeLong         Type = "long"          // signed int64
	TypeInteger      Type = "integer"       // signed int32
	TypeShort        Type = "short"         // signed int16
	TypeByte         Type = "byte"          // signed int8

	// Float types
	TypeDouble Type = "double" // float64
	TypeFloat  Type = "float"  // float32
)

func (t Type) Valid() bool {
	return t == TypeBool ||
		t == TypeKeyword ||
		t == TypeText ||
		t == TypeSlice ||
		t == TypeMap ||
		t == TypeUnsignedLong ||
		t == TypeLong ||
		t == TypeInteger ||
		t == TypeShort ||
		t == TypeByte ||
		t == TypeDouble ||
		t == TypeFloat
}

type Field struct {
	Type     Type             `json:"type"`
	Required bool             `json:"required"`
	Children map[string]Field `json:"children"`
	Analyzer string           `json:"analyzer"`
}

func NewField(fieldType Type, required bool, analyzer string) Field {
	return Field{
		Type:     fieldType,
		Required: required,
		Analyzer: analyzer,
	}
}

func NewFieldWithChildren(fieldType Type, required bool, analyzer string, children map[string]Field) Field {
	return Field{
		Type:     fieldType,
		Required: required,
		Analyzer: analyzer,
		Children: children,
	}
}

func (f Field) ValidateWithContext(ctx context.Context) error {
	if f.Children != nil {
		if err := validateKeys("fields", f.Children); err != nil {
			return err
		}
	}

	return validation.ValidateStructWithContext(ctx, &f,
		validation.Field(&f.Type, validation.Required, validation.By(validateFieldType())),
		validation.Field(
			&f.Analyzer,
			validation.When(f.Type == TypeText, validation.Required),
			validation.WithContext(validateFieldAnalyzers(f.Type))),
		validation.Field(&f.Children, validation.By(validateFieldChildren(f.Type))),
	)
}

func validateFieldType() validation.RuleFunc {
	return func(value interface{}) error {
		v := value.(Type)
		if !v.Valid() {
			return errs.Errorf("invalid field type %q", v)
		}
		return nil
	}
}

func validateFieldAnalyzers(t Type) validation.RuleWithContextFunc {
	return func(ctx context.Context, value interface{}) error {
		v := value.(string)
		if v == "" {
			return nil
		}

		s := ctx.Value(ctxKeySchema).(Schema)
		if _, ok := s.Analyzers[v]; !ok {
			return errs.Errorf("unknown analyzer %q", v)
		}

		return nil
	}
}

func validateFieldChildren(t Type) validation.RuleFunc {
	return func(value interface{}) error {
		if value == nil {
			if t == TypeSlice || t == TypeMap {
				return errs.Errorf("type %q must have children defined", t)
			}
			return nil
		}
		v := value.(map[string]Field)
		if len(v) == 0 {
			if t == TypeSlice || t == TypeMap {
				return errs.Errorf("type %q must have children defined", t)
			}
			return nil
		}

		if len(v) != 0 && t != TypeSlice && t != TypeMap {
			return errs.Errorf("type %q cannot have children fields", t)
		}

		return nil
	}
}
