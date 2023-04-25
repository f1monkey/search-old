package usecase

import (
	"github.com/f1monkey/search/internal/index"
)

type IndexGetter struct {
	getter indexGetter
}

type indexGetter func(name string) (index.Index, error)

func NewIndexGet(getter indexGetter) *IndexGetter {
	return &IndexGetter{
		getter: getter,
	}
}

func (u *IndexGetter) Get(name string) (index.Index, error) {
	return u.getter(name)
}
