package index

import "github.com/invopop/validation"

type Index struct {
	Name   string `json:"name"`
	Schema Schema `json:"schema"`
}

func (i Index) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required),
		validation.Field(&i.Schema, validation.Required),
	)
}
