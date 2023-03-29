package index

import (
	"time"

	"github.com/f1monkey/search/internal/index/schema"
)

type Index struct {
	Name      string
	CreatedAt time.Time
	Schema    schema.Schema
}

func New(name string, schema schema.Schema) Index {
	return Index{
		Name:      name,
		CreatedAt: time.Now(),
		Schema:    schema,
	}
}
