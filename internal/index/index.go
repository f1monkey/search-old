package index

import (
	"github.com/f1monkey/search/internal/index/schema"
	"github.com/invopop/validation"
)

type Index struct {
	Name   string        `json:"name"`
	Schema schema.Schema `json:"schema"`
}

func (i Index) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required),
		validation.Field(&i.Schema, validation.Required),
	)
}
