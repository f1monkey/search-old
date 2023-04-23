package index

import (
	"time"
)

type Index struct {
	Name      string
	CreatedAt time.Time
	Schema    Schema
}

func New(name string, schema Schema) Index {
	return Index{
		Name:      name,
		CreatedAt: time.Now(),
		Schema:    schema,
	}
}
